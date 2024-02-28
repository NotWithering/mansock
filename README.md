# ManSock [![MIT License](https://img.shields.io/badge/License-MIT-a10b31)](https://github.com/NotWithering/mansock/blob/master/LICENSE)
<img src="sock.png" width=150px alt="sock with roblox man face">

**ManSock** is a simple CLI program used to debug tcp/udp servers with an interactive command-line

## Installing
```bash
go install github.com/NotWithering/mansock@latest
```
## Example
```go
$ ./udpserver &

$ mansock
MANSOCK> set PORT string "42480"
MANSOCK> set PROTO string "udp"
MANSOCK> wb string "Hello, world!"
MANSOCK> c
MANSOCK> ws

UDP Server: "Hello, world!"
```

## What does ManSock mean?
Despite the funny name, it actually stands for something

**Man**ual **Sock**et