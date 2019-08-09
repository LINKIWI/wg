# wg

**wg** is a command line client for webgrep.

## Usage

### Building

```bash
$ make
```

### Runtime

`wg` accepts search queries via stdin.

```bash
$ echo "query" | ./bin/wg-$OS-$ARCH --webgrep-url https://grep.example.com
```
