#!/usr/bin/env -S bash -e

domain_name=acme.example
if [[ $DOMAIN_NAME != "" ]]; then
  domain_name=$DOMAIN_NAME
fi

rm -rf "private/$domain_name"

mkdir -p "private/$domain_name"

cat <<EOF > "private/$domain_name/req.cnf"
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no
[req_distinguished_name]
C = US
ST = VA
L = SomeCity
O = MyCompany
OU = MyDivision
CN = $domain_name
[v3_req]
keyUsage = critical, digitalSignature, keyAgreement
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = $domain_name
DNS.2 = www.$domain_name
DNS.3 = api.$domain_name
EOF

openssl \
  req -x509 -newkey rsa:4096 -sha256 -nodes \
  -keyout "private/$domain_name/$domain_name.key" \
  -out "private/$domain_name/$domain_name.pem" \
  -subj "/CN=$domain_name" \
  -days 365 \
  -config "private/$domain_name/req.cnf"

mv "private/$domain_name/$domain_name.key" "private/$domain_name/$domain_name.pem.key"
