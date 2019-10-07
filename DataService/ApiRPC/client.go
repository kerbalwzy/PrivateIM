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
	cert, err := tls.LoadX509KeyPair(conf.DataLayerSrvCAClientPem, conf.DataLayerSrvCAClientKey)
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(conf.DataLayerSrvCAPem)
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
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
