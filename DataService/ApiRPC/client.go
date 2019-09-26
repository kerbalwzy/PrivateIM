package ApiRPC

import (
	conf "../Config"
	pb "../Protos"

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
	mysqlDataClient pb.MySQLBindServiceClient
	mongoDateClient pb.MongoBindServiceClient
)

func init() {
	// Add CA TSL authentication data
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

	mysqlDataClient = pb.NewMySQLBindServiceClient(conn1)

	conn2, err := grpc.Dial(conf.MongoDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}

	mongoDateClient = pb.NewMongoBindServiceClient(conn2)
}

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

// Return the client for mysql data rpc sever. The type of client is pointer.
func GetMySQLDateClient() pb.MySQLBindServiceClient {
	return mysqlDataClient
}

// Return the client for mongo data rpc sever. The type of client is pointer.
func GetMongoDateClient() pb.MongoBindServiceClient {
	return mongoDateClient
}
