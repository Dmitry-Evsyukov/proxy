#!/bin/sh
openssl req -new -key ca.key -subj "/CN=dmitry" -sha256 | openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -set_serial "$1" > nck.crt
