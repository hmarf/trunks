# Trunks
<img src="https://github.com/hmarf/trunks/blob/master/img/trunks.jpg?raw=true" width="500px">

## Overview
Trunks is a simple command line tool for HTTP load testing. 

## Usage
```
NAME:
   trunks - Trunks is a simple command line tool for HTTP load testing.

USAGE:
   trunks [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   hmarf

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --concurrency value, -c value  Concurrency Level (default: 10)
   --requests value, -r value     Number of Requests (default: 100)
   --url value, -u value          URL to hit (default: "None")
   --help, -h                     show help
   --version, -v                  print the version
```

# Reference
### Benchmark
- https://github.com/tsenart/vegeta
### Keep Alive
- https://www.sambaiz.net/article/61/
- https://christina04.hatenablog.com/entry/go-keep-alive
