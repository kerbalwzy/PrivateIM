// CA TLS auth files generate record, you can update the files by yourself.
// When the terminal need input values for generate the file, remember the information
// should always be same in every input time.

openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 3650 -key ca.key -out ca.pem

openssl ecparam -genkey -name secp384r1 -out ./server/server.key

openssl req -new -key ./server/server.key -out ./server/server.csr

openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in ./server/server.csr -out ./server/server.pem

openssl ecparam -genkey -name secp384r1 -out ./client/client.key

openssl req -new -key ./client/client.key -out ./client/client.csr

openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in ./client/client.csr -out ./client/client.pem