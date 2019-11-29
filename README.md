# Trunks
<img src="https://github.com/hmarf/trunks/blob/master/img/trunks.jpg?raw=true" width="400px">

## Overview
Trunks is a simple command line tool for HTTP load testing. 

## Demo
![demo](https://github.com/hmarf/trunks/blob/master/img/trunks.gif)

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
   --url value, -u value          [required] string
                                   URL to hit (default: "None")
   --concurrency value, -c value  int
                                   Concurrency Level. (default: 10)
   --requests value, -r value     int
                                   Number of Requests. (default: 100)
   --method value, -m value       string
                                   http method. (default: "GET")
   --header value, -H value       string
                                   HTTP header
   --body value, -b value         string
                                   HTTP body
   --output value, -o value       string
                                   File name to output results
   --help, -h                     show help
   --version, -v                  print the version
```

## Install
```
brew tap hmarf/trunks
brew install trunks
```

or 
```
go get -u github.com/hmarf/trunks
```

## Example
### input
#### command option
- 10,000 requests are sent to 127.0.0.1:8080 in 10 parallels
```
trunks -c 10 -r 10000 -u "http://localhost:8080"
```

- 10,000 requests with header and body set are sent to 127.0.0.1:8080 in 10 parallels 
```
trunks -c 10 -r 10000 -u "http://localhost:8080" -H "Content-Type:application/json" -H Accept: "application/json" -b "{"message":"Welcome to underground"}"
```

- Specify a file to save the results
```
trunks -c 10 -r 10000 -u "http://localhost:8080" -o "output-file.json"
```

### output
```
$ trunks -c 10 -r 10000 -u "http://localhost:8080"

[00:00:01] [##################################################] 100%

Concurrency Level:   10
Total Requests:      10000
Succeeded requests:  10000
Failed requests:     0
Requests/sec:        0
Total data received: 120000 bytes

Status code:
   [200] 10000 responses

Latency:
   total: 1.238740417s
   max:   49.110501ms
   min:   161.2µs
   ave:   1.215923ms
```

# Reference
### Benchmark
- https://github.com/tsenart/vegeta
### Keep Alive
- https://www.sambaiz.net/article/61/
- https://christina04.hatenablog.com/entry/go-keep-alive
