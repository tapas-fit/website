#!/bin/bash
GOOS=freebsd go build -o waitlist-freebsd
go build -o waitlist-linux
