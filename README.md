[![Go Report Card](https://goreportcard.com/badge/github.com/Tanmay9511/Go-Portscanner)](https://goreportcard.com/report/github.com/Tanmay9511/Go-Portscanner)
[![GoDoc](https://godoc.org/github.com/Tanmay9511/Go-Portscanner?status.svg)](https://godoc.org/github.com/Tanmay9511/Go-Portscanner)

# Go-Portscanner
Port Scanner in Go language

A port scanner written in Go language which scans TCP Ports with very simple thread-safe which should work on every OS without any problems. Port scanning will perform 65535 ports and show the results.

Docker implemention are also availabe for port scanning

![alt text](https://github.com/Tanmay-N/Go-Portscanner/blob/master/GO_Portscanner%20Screenshot.png?raw=true)

## Installation 

```javascript 
git clone https://github.com/Tanmay-N/Go-Portscanner.git portscanner
cd portscanner
go build scanner.go -o scanner
```
Or 

```javascript 
go get -u github.com/Tanmay-N/Go-Portscanner
```

## Usage

### Get help

```javascript 
scanner.go -h 
-IP string
        IP Address/Domain name. (Required) (default "127.0.0.1")
```
Shows the available options.

### Localhost Port Scan

```javascript 
go run .\scanner.go 
```
Scans all local machine ports, from 1 to 65535.

### Port Scan on Domain 

```javascript 
go run .\scanner.go -IP="scanme.nmap.org"
```
Scan for particular domain.

### Network Discovery 

```javascript 
go run .\scanner.go -IP="192.168.0.1-225"
```
Looks for HTTP or FTP servers on 192.168.0.0/24.

## Docker installation

To start a port-scan container, simply run gievn command in docker:

```javascript
docker run -it tanmay95/go-portscanner -IP="scanme.nmap.org"
```
