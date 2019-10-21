[![Go Report Card](https://goreportcard.com/badge/github.com/Tanmay9511/Go-Portscanner)](https://goreportcard.com/report/github.com/Tanmay9511/Go-Portscanner)
[![GoDoc](https://godoc.org/github.com/Tanmay9511/Go-Portscanner?status.svg)](https://godoc.org/github.com/Tanmay9511/Go-Portscanner)

# Go-Portscanner
Port Scanner in Go language

A port scanner written in Go language which scans TCP Ports with very simple thread-safe which should work on every OS without any problems. Port scanning will perform 65535 ports and show the results.

Docker implemention are also availabe for port scanning

## Installation 

```javascript 
git clone https://github.com/Tanmay-N/Go-Portscanner.git
```

Or 
```javascript 
go get -u github.com/Tanmay-N/Go-Portscanner
```

## Usage

```javascript 
go run .\scanner.go -IP="scanme.nmap.org"
```

## Docker installation

To start a port-scan container, simply run gievn command in docker:

```javascript
docker run -it tanmay95/go-portscanner -IP="scanme.nmap.org"
```
