#!/usr/bin/env bash

# Create certs dir if it does not exists
mkdir -p ../certs

######################
# Become a Certificate Authority
######################

# Generate private key
openssl genrsa -des3 -out ../certs/visualonCA.key 2048
# Generate root certificate
openssl req -x509 -new -nodes -key ../certs/visualonCA.key -sha256 -days 825 -out ../certs/visualonCA.pem

######################
# Create CA-signed certs
######################

NAME=live-encoding-status.visualon.info # Use your own domain name
# Generate a private key
openssl genrsa -out ../certs/$NAME.key 2048
# Create a certificate-signing request
openssl req -new -key ../certs/$NAME.key -out ../certs/$NAME.csr
# Create a config file for the extensions
>../certs/$NAME.ext cat <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = DNS:live-encoding-status.visualon.info
[alt_names]
DNS.1 = $NAME # Be sure to include the domain name here because Common Name is not so commonly honoured by itself
IP.1 = 23.20.1.104 # Optionally, add an IP address (if the connection which you have planned requires it)
EOF
# Create the signed certificate
openssl x509 -req -in ../certs/$NAME.csr -CA ../certs/visualonCA.pem -CAkey ../certs/visualonCA.key -CAcreateserial \
-out ../certs/$NAME.crt -days 825 -sha256 -extfile ../certs/$NAME.ext