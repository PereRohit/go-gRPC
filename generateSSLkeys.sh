#!/bin/bash

# OUTPUT FILES
# ca.key:     Certificate Authority private key file         (private)
# ca.crt:     Certificate Authority trust certificate        (public)  -- required by client
# server.key: Server private key, password protected         (private)
# server.csr: Server certificate signing request             (shared)  -- shared with CA owner
# server.crt: Server certificate signed by the CA            (private) -- returned by CA owner
# server.pem: Converted server.key to gRPC recognized format (private)

# All ssl files will be generated in the ssl directory
mkdir -p ssl
cd ssl/

# SERVER_CN contains host URLs in the environment
SERVER_CN=

# PASS is the password used to generate the keys in the environment
PASS=

# Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:${PASS} -des3 -out ca.key 4096
openssl req -passin pass:${PASS} -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# Generate Server Private Key (server.key)
openssl genrsa -passout pass:${PASS} -des3 -out server.key 4096

# Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:${PASS} -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# Sign the certificate with the CA, thereby performing self signing (server.crt)
openssl x509 -req -passin pass:${PASS} -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Convert the server certificate to .pem format to be used by gRPC (server.pem)
openssl pkcs8 -topk8 -nocrypt -passin pass:${PASS} -in server.key -out server.pem