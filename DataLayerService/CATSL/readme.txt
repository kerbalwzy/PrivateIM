// CA TLS auth files generate record, you can update the files by yourself.

openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 3650 -key ca.key -out ca.pem
/* input value:
----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:PrivateIM
Email Address []:
*/




openssl ecparam -genkey -name secp384r1 -out ./server/server.key

openssl req -new -key ./server/server.key -out ./server/server.csr
/* input value:
-----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:PrivateIM
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
*/

openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in ./server/server.csr -out ./server/server.pem





openssl ecparam -genkey -name secp384r1 -out ./client/client.key

openssl req -new -key ./client/client.key -out ./client/client.csr
/*input value:
-----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:PrivateIM
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
*/

openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in ./client/client.csr -out ./client/client.pem