# Trunks
<img src="https://github.com/hmarf/trunks/blob/master/img/trunks.jpg?raw=true" width="500px">

## Overview
Trunks is a simple command line tool for HTTP load testing. 

## Usage
```
NAME:
   trunks - Trunks is a simple command line tool for HTTP load testing.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   hmarf

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --concurrency value, -c value  [int] Concurrency Level. (default: 10)
   --requests value, -r value     [int] Number of Requests. (default: 100)
   --method value, -m value       [string] http method. (default: "GET")
   --url value, -u value          [string required] URL to hit (default: "None")
   --header value, -H value       [string] HTTP header
   --body value, -b value         [string] HTTP body
   --output value, -o value       [string] File name to output results
   --help, -h                     show help
   --version, -v                  print the version
```

# Reference
### Benchmark
- https://github.com/tsenart/vegeta
### Keep Alive
- https://www.sambaiz.net/article/61/
- https://christina04.hatenablog.com/entry/go-keep-alive
