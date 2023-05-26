<h1 align="center">nmapParser</h1>
<h4 align="center">Tool for parsing NMAP output and creating different formats</h4>
<p align="center">
  
  <img src="https://img.shields.io/github/watchers/secinto/nmapParser?label=Watchers&style=for-the-badge" alt="GitHub Watchers">
  <img src="https://img.shields.io/github/stars/secinto/nmapParser?style=for-the-badge" alt="GitHub Stars">
  <img src="https://img.shields.io/github/license/secinto/nmapParser?style=for-the-badge" alt="GitHub License">
</p>

Developed by Stefan Kraxberger (https://twitter.com/skraxberger/)  

Released as open source by secinto GmbH - https://secinto.com/  
Released under Apache License version 2.0 see LICENSE for more information

Description
----
nmapParser is a GO tool which parses NMAP output and creates JSON outputs for further processing.

# Installation Instructions

`nmapParser` requires **go1.20** to install successfully. Run the following command to get the repo:

```sh
git clone https://github.com/secinto/nmapParser.git
cd parser
go build
go install
```

or the following to directly install it from the command line:

```sh
go install -v github.com/secinto/parser/cmd/parser@latest
```

# Usage

```sh
parser -help
```

This will display help for the tool. Here are all the switches it supports.


```console
Usage:
  ./nmapParser [flags]

Flags:
   -p                    project name which will be added as additional information to the data
