package ApiRPC

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"time"

	conf "../Config"

	"../RpcClientPbs/mongoPb"
	"../RpcClientPbs/mysqlPb"
	"../RpcClientPbs/userAuthPb"
)

var (
	mongoDataClient mongoPb.MongoBindServiceClient
	mysqlDataClient mysqlPb.MySQLBindServiceClient
	userAuthClient  userAuthPb.UserAuthClient
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

	// get the client for calling the user-auth-rpc-server
	conn1, err := grpc.Dial(conf.UserAuthRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}
	userAuthClient = userAuthPb.NewUserAuthClient(conn1)

	// get the client for calling the mongo-data-rpc-server
	conn2, err := grpc.Dial(conf.MongoDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}

	mongoDataClient = mongoPb.NewMongoBindServiceClient(conn2)

	conn3, err := grpc.Dial(conf.MySQLDataRPCServerAddress, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatal(err.Error())
	}
	mysqlDataClient = mysqlPb.NewMySQLBindServiceClient(conn3)
}

// Return a context instance with deadline
func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

// Return the client for mongo data rpc sever. The type of client is pointer.
func GetMongoDataClient() mongoPb.MongoBindServiceClient {
	return mongoDataClient
}

// Return the client for mysql data rpc server. The type of client is pointer.
func GetMySQLDataClient() mysqlPb.MySQLBindServiceClient {
	return mysqlDataClient
}

// Return the client for user aut rpc server.
func GetUserAuthClient() userAuthPb.UserAuthClient {
	return userAuthClient
}
