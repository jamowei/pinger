# Pinger
Pinger is a small lightweigt command line programm written with ❤️ in go.
It's made for testing TCP connections between to host. Therefore you must run it on both hosts, on one site in server mode and on the other side in client mode. Then the programm on the client side tries to reach the endpoint on the server side and logs it it was successfull.

[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](https://github.com/jamowei/pinger/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/jamowei/pinger.svg?branch=master)](https://travis-ci.org/jamowei/pinger)

# Installation

You can get the latest binary using Go:

`> go get -u github.com/jamowei/pinger`

or download released binary from [here](https://github.com/jamowei/pinger/releases/latest).

# Commandline

If you define a server name (or hostname) with `-s` r `--server`, its starts in client mode an tries to connect to the server endpoint on the specified port(s).

Otherwise if you define only the ports you want to test, it starts a small http server for each port.

```
usage: pinger [-h|--help] [-s|--server "<value>"] [-r|--range "<value>"]
              [-p|--port "<value>" [-p|--port "<value>" ...]]

              Pings a running pinger server or starts a server to look for open
              tcp ports

Arguments:

  -h  --help    Print help information
  -s  --server  host to connect with (server mode)
  -r  --range   range of port numbers, e.g. 8080-8090
  -p  --port    ports to listen on (server mode) or to connect with (client
                mode)
```

# License

Pinger is released under the MIT license. See [LICENSE](https://github.com/jamowei/pinger/blob/master/LICENSE)