## establish secure connection via ssl
create private_key
``` bash
openssl ecparam -genkey -name secp384r1 -out server.key
``` 

store above as server.cnf
```
[ req ]
default_bits       = 4096
default_md		= sha256
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = BEIJING
localityName                = Locality Name (eg, city)
localityName_default        = BEIJING
organizationName            = Organization Name (eg, company)
organizationName_default    = DEV
commonName                  = Common Name (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = atelier.icu

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1   = localhost
DNS.2   = atelier.icu
IP      = 127.0.0.1
```
generate certificate
```bash
openssl req -nodes -new -x509 -sha256 -days 3650 -config server.cnf -extensions 'req_ext' -key server.key -out server.crt
```