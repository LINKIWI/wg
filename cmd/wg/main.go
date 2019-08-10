package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/proxy"

	"wg/internal/cli"
	"wg/internal/meta"
	"wg/pkg/webgrep"
)

var (
	flagWebgrepURL    = flag.String("webgrep-url", "", "base URL of the webgrep instance")
	flagRegex         = flag.Bool("regex", false, "interpret search query as a regular expression")
	flagCaseSensitive = flag.Bool("case-sensitive", false, "respect search query case sensitivity")
	flagFile          = flag.String("file", "", "filter matches by file path pattern")
	flagMaxMatches    = flag.Int("max-matches", 50, "maximum number of matches in search results")
	flagProxy         = flag.String("proxy", "", "optional address of a SOCKS5 proxy server")
	flagVersion       = flag.Bool("version", false, "print the application version and exit")
	flagRepos         = cli.NewArrayFlag()
	flagSearchType    = cli.NewChoicesFlag([]string{"files", "code"}, "code")
)

func init() {
	flag.Var(flagRepos, "repo", "filter matches by repository name")
	flag.Var(flagSearchType, "search-type", "search results type to print; one of {files, code}")
	flag.Parse()
}

func main() {
	if *flagVersion {
		fmt.Printf("wg/%s\n", meta.Version)
		return
	}

	// Rudimentary input validation
	if *flagWebgrepURL == "" {
		panic(fmt.Errorf("main: no value specified for webgrep instance URL"))
	}

	// Optional proxy server configuration
	var backend *http.Client
	if *flagProxy != "" {
		dialer, err := proxy.SOCKS5("tcp", *flagProxy, nil, proxy.Direct)
		if err != nil {
			panic(err)
		}

		backend = &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	}

	// Instantiate a webgrep client
	client, err := webgrep.NewClient(*flagWebgrepURL, backend)
	if err != nil {
		panic(err)
	}

	// Read search query as input from stdin
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// Execute the search, respecting CLI flags as parameters
	resp, searchErr := client.Search(&webgrep.SearchQueryRequest{
		Query:         strings.TrimSpace(input),
		File:          *flagFile,
		Repos:         flagRepos.Values(),
		Regex:         *flagRegex,
		CaseSensitive: *flagCaseSensitive,
		MaxMatches:    *flagMaxMatches,
	})
	if searchErr != nil {
		panic(searchErr)
	}

	rendered := cli.NewTable()

	// Format results as requested by parameters
	switch flagSearchType.Choice() {
	case "code":
		for _, result := range resp.Code {
			for _, line := range result.Lines {
				source := line.Line
				if len(line.Bounds) == 2 {
					// Apply highlighting to the matching portion of the source
					// line, if applicable
					source = strings.Join([]string{
						line.Line[:line.Bounds[0]],
						cli.Highlight(line.Line[line.Bounds[0]:line.Bounds[1]]),
						line.Line[line.Bounds[1]:],
					}, "")
				}

				row := []string{
					result.Version,
					result.Repo,
					result.Path,
					strconv.Itoa(line.Number),
					fmt.Sprintf("|%s", source),
				}

				if err := rendered.Add(row); err != nil {
					panic(err)
				}
			}
		}

	case "files":
		for _, result := range resp.Files {
			path := strings.Join([]string{
				result.Path[:result.Bounds[0]],
				cli.Highlight(result.Path[result.Bounds[0]:result.Bounds[1]]),
				result.Path[result.Bounds[1]:],
			}, "")

			row := []string{
				result.Version,
				result.Repo,
				path,
			}

			if err := rendered.Add(row); err != nil {
				panic(err)
			}
		}

	default:
	}

	if !rendered.IsEmpty() {
		fmt.Println(rendered)
	}
}
