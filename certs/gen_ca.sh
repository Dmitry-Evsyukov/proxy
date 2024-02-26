#!/bin/sh
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=dmitry proxy CA"
openssl genrsa -out cert.key 2048enssl.cnf -key ca.key.pem -new -x509 -days 7300 -sha256 -extensions v3_ca -out ca.pem
