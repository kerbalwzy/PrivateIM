package ApiRPC

import (
	conf "../Config"
	mongoPb "../Protos/mongoProto"
	mysqlPb "../Protos/mysqlProto"

	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"time"
)

var (
	mysqlDataClient mysqlPb.MySQLBindServiceClient
	mongoDateClient mongoPb.MongoBindServiceClient
)

func init() {
	// Add CA TLS authentication data
	cert, err := tls.LoadX509KeyPair(conf.PrivateIMClientPem, conf.PrivateIMClientKey)
	if err != nil {
		log.Fatalf("[error] load CA X509 key files fail: %s", err.Error())
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(conf.PrivateIMRootCAPem)
	if err != nil {
		log.Fatalf("[error] load CA Root Pem fail: %s", err.Error())
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("[error] certPool.AppendCertsFromPEM error")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "PrivateIM",
		RootCAs:      certPool,
	})

	conn1, err := grpc.Dial(conf.MySQLDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}

	mysqlDataClient = mysqlPb.NewMySQLBindServiceClient(conn1)

	conn2, err := grpc.Dial(conf.MongoDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}

	mongoDateClient = mongoPb.NewMongoBindServiceClient(conn2)
}

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

// Return the client for mysql data rpc sever. The type of client is pointer.
func GetMySQLDateClient() mysqlPb.MySQLBindServiceClient {
	return mysqlDataClient
}

// Return the client for mongo data rpc sever. The type of client is pointer.
func GetMongoDateClient() mongoPb.MongoBindServiceClient {
	return mongoDateClient
}
