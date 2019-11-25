package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"wg/internal/cli"
	"wg/internal/meta"
	"wg/pkg/webgrep"
)

const (
	envWebgrepURL = "WG_WEBGREP_URL"
	envProxyAddr  = "WG_PROXY_ADDR"
)

var (
	flagWebgrepURL    = flag.String("webgrep-url", os.Getenv(envWebgrepURL), "base URL of the webgrep instance")
	flagRegex         = flag.Bool("regex", false, "interpret search query as a regular expression")
	flagCaseSensitive = flag.Bool("case-sensitive", false, "respect search query case sensitivity")
	flagFile          = flag.String("file", "", "filter matches by file path pattern")
	flagMaxMatches    = flag.Int("max-matches", 50, "maximum number of matches in search results")
	flagContext       = flag.Int("context", 4, "number of surrounding context lines to include around matches")
	flagProxy         = flag.String("proxy", os.Getenv(envProxyAddr), "optional address of a SOCKS5 proxy server")
	flagAbout         = flag.Bool("about", false, "print current server-side index metadata")
	flagRepos         = cli.NewArrayFlag()
	flagSearchType    = cli.NewChoicesFlag([]string{"files", "code"}, "code")
)

func init() {
	flag.Var(flagRepos, "repo", "filter matches by repository name")
	flag.Var(flagSearchType, "search-type", "search results type to print; one of {files, code}")
	flag.Parse()
}

func main() {
	// Rudimentary input validation
	if *flagWebgrepURL == "" {
		panic(fmt.Errorf("main: no value specified for webgrep instance URL (see --help for docs)"))
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

	// Application and index metadata
	if *flagAbout {
		if err := about(client); err != nil {
			panic(err)
		}
		return
	}

	// Code search query
	if err := search(client); err != nil {
		panic(err)
	}
}

func about(client *webgrep.Client) error {
	metadata, err := client.Metadata()
	if err != nil {
		return err
	}

	table := cli.NewTable()
	table.Add([]string{"wg client version:", meta.Version})
	table.Add([]string{"webgrep server version:", metadata.Version})
	table.Add([]string{"index name:", metadata.Name})
	table.Add([]string{"index timestamp:", time.Unix(int64(metadata.Timestamp), 0).String()})
	table.Add([]string{"index repositories:", strconv.Itoa(len(metadata.Repositories))})

	fmt.Println(table)

	return nil
}

func search(client *webgrep.Client) error {
	// Read search query as input from stdin
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Execute the search, respecting CLI flags as parameters
	resp, searchErr := client.Search(&webgrep.SearchQueryRequest{
		Query:         strings.TrimSpace(input),
		File:          *flagFile,
		Repos:         flagRepos.Values(),
		Regex:         *flagRegex,
		CaseSensitive: *flagCaseSensitive,
		MaxMatches:    *flagMaxMatches,
		Context:       *flagContext,
	})
	if searchErr != nil {
		return searchErr
	}

	table := cli.NewTable()

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

				if err := table.Add(row); err != nil {
					return err
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

			if err := table.Add(row); err != nil {
				return err
			}
		}

	default:
	}

	if !table.IsEmpty() {
		fmt.Println(table)
	}

	return nil
}
