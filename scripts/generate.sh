#!/bin/bash

function main {
   case $1 in
     cleanup)
         cleanup
     ;;
     init)
         init
         openssl_conf
     ;;
     generate)
         gen_ca
         gen_server
         gen_client
     ;;
     install)
         install_fluentbit_secret
         install_fluentd_secret
     ;;
     all)
         init
         openssl_conf
         gen_ca
         gen_server
         gen_client
         install_fluentbit_secret
         install_fluentd_secret
     ;;
     *)
         commands
     ;;
   esac
}

function commands {
    echo "Usage: "
    echo "$0 cleanup  - Remove all files and folders created during process"
    echo "$0 init     - Create directories and generate openssl.conf"
    echo "$0 generate - Generate CA, Server, Client keys and certificates"
    echo "$0 install  - Installs fluentbit-tls and fluentd-tls secrets on the kubernetes cluster"
    echo "$0 all      - Do all the steps automatically"
}

function cleanup {
    rm -rf certs csr private
    rm index.txt openssl.cnf serial
}

function init {
  mkdir -p certs csr private
  touch index.txt
  echo "1000" > serial
}

function openssl_conf {
cat << EOF > openssl.cnf
[ ca ]
#
default_ca = CA_default

[ CA_default ]
# Directory and file locations.
dir               = .
certs             = ./certs
crl_dir           = ./crl
new_certs_dir     = ./newcerts
database          = ./index.txt
serial            = ./serial
RANDFILE          = ./private/.rand

# The root key and root certificate.
private_key       = /private/ca.key.pem
certificate       = /certs/ca.crt.pem

# SHA-1 is deprecated, so use SHA-2 instead.
default_md        = sha256

name_opt          = ca_default
cert_opt          = ca_default
default_days      = 365
preserve          = no
policy            = policy_strict

[ req ]
# Options for the  tool ().
default_bits        = 4096
distinguished_name  = req_distinguished_name
string_mask         = utf8only

# SHA-1 is deprecated, so use SHA-2 instead.
default_md          = sha256

# Extension to add when the -x509 option is used.
x509_extensions     = v3_ca

[ req_distinguished_name ]
# See <https://en.wikipedia.org/wiki/Certificate_signing_request\>.
countryName                     = Country Name (2 letter code)
stateOrProvinceName             = State or Province Name
localityName                    = Locality Name
0.organizationName              = Organization Name
organizationalUnitName          = Organizational Unit Name
commonName                      = Common Name (required)
emailAddress                    = Email Address

# Optionally, specify some defaults.
countryName_default             = US
stateOrProvinceName_default     = CA
#localityName_default           = Mountain View
0.organizationName_default      = Your company name
#organizationalUnitName_default =
emailAddress_default            = foo@example.com

[v3_ca]
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always, issuer
basicConstraints = critical,CA:true
keyUsage = critical, cRLSign, digitalSignature, keyCertSign

[ client_cert ]
# Extensions for client certificates ().
basicConstraints = CA:FALSE
nsCertType = client, email
nsComment = "OpenSSL Generated Client Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer
keyUsage = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, emailProtection

[ server_cert ]
# Extensions for server certificates ().
basicConstraints = CA:FALSE
nsCertType = server
nsComment = "OpenSSL Generated Server Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth

[ policy_strict ]
# The root CA should only sign intermediate certificates that match.
# See the POLICY FORMAT section of .
countryName             = match
stateOrProvinceName     = match
organizationName        = match
organizationalUnitName  = optional
commonName              = supplied
emailAddress            = optional
EOF
}

function gen_ca {
  echo "Generate CA private key..."
  openssl genrsa -out private/ca.key 4096

  echo "Generate CA certificate..."
  openssl req -config openssl.cnf \
    -key private/ca.key \
    -new -x509 -days 365 -sha256 -extensions v3_ca \
    -out certs/ca.crt \
    -subj '/CN=logging-ca/O=Banzai Cloud/C=HU/ST=Budapest'
}

function gen_server {
  echo "Generate Server (fluentd) private key..."
  openssl genrsa -out private/server.key 4096

  echo "Generate Server CSR (fluentd)..."
  openssl req -config openssl.cnf \
    -key private/server.key \
    -new -out csr/server.csr \
    -subj '/CN=fluentd/O=Banzai Cloud/C=HU/ST=Budapest'

  echo "Generate Server certificate (fluentd)..."
  openssl ca -config openssl.cnf -outdir certs  \
    -cert certs/ca.crt \
    -keyfile private/ca.key \
    -extensions server_cert -days 365 -notext -md sha256 \
    -in csr/server.csr \
    -out certs/server.crt \
    -subj '/CN=fluentd/O=Banzai Cloud/C=HU/ST=Budapest'
}

function gen_client {
  echo "Generate Client (fluent-bit) private key..."
  openssl genrsa -out private/client.key 4096

  echo "Generate Client CSR (fluent-bit)..."
  openssl req -config openssl.cnf \
    -key private/client.key \
    -new \
    -out csr/client.csr \
    -subj '/CN=fluent-bit/O=Banzai Cloud/C=HU/ST=Budapest'

  echo "Generate Client certificate (fluent-bit)..."
  openssl ca -config openssl.cnf -outdir certs  \
    -cert certs/ca.crt \
    -keyfile private/ca.key \
    -extensions client_cert -days 365 -notext -md sha256 \
    -in csr/client.csr \
    -out certs/client.crt \
    -subj '/CN=fluent-bit/O=Banzai Cloud/C=HU/ST=Budapest'
}

function install_fluentbit_secret {
  kubectl create secret generic fluentbit-tls \
    --from-file=ca.crt=./certs/ca.crt \
    --from-file=tls.crt=./certs/client.crt \
    --from-file=tls.key=./private/client.key
}

function install_fluentd_secret {
kubectl create secret generic fluentd-tls \
    --from-file=ca.crt=./certs/ca.crt \
    --from-file=tls.crt=./certs/server.crt \
    --from-file=tls.key=./private/server.key
}

main $@
