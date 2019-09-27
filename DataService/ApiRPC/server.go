package ApiRPC

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"time"

	conf "../Config"
	mongoPb "../Protos/mongoProto"
	mysqlPb "../Protos/mysqlProto"
)

var CtxCanceledErr = errors.New("the client canceled or connection time out")

// out the log information for every handler
func logWithHandler(handlerName string, err error, timeConsumed float64) {
	outString := "[%s] %s: "
	if nil != err {
		outString += "%s TimeConsumed(%f)s"
		log.Printf(outString, "error", handlerName, err.Error(), timeConsumed)
	} else {
		outString += "success TimeConsumed(%f)s"
		log.Printf(outString, "info", handlerName, timeConsumed)
	}
}

// Check the client if canceled the calling or connection is time out
func checkCtxCanceled(ctx context.Context) error {
	if ctx.Err() == context.Canceled {
		return CtxCanceledErr
	}
	return nil
}

// Create a server option to use CA TSL authentication for keeping safe for data transmission.
func getCAOption() grpc.ServerOption {
	cert, err := tls.LoadX509KeyPair(conf.DataLayerSrvCAServerPem, conf.DataLayerSrvCAServerKey)
	if err != nil {
		log.Fatalf("[error] getCAOption: %s", err.Error())
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(conf.DataLayerSrvCAPem)
	if err != nil {
		log.Fatalf("[error] getCAOption: %s", err.Error())
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("[error] getCAOption: certPool.AppendCertsFromPEM Error")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return grpc.Creds(c)
}

// new an interceptor, function similar to middleware
func getUnaryInterceptorOption() grpc.ServerOption {
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// check the call if canceled by client or time out before every handler
		err = checkCtxCanceled(ctx)
		if err != nil {
			return
		}

		// continue the handel, record the time consumed and log information
		startTime := time.Now()
		resp, err = handler(ctx, req)
		timeConsumed := time.Now().Sub(startTime).Seconds()

		logWithHandler(info.FullMethod, err, timeConsumed)

		return resp, err
	}
	// add the interceptor for every Unary-Unary handler
	return grpc.UnaryInterceptor(interceptor)
}

// Start the gRPC server for MySQL data operation.
func StartMySQLDataRPCServer() {
	// using CA TSL authentication
	caOption := getCAOption()
	log.Printf("[info] StartMySQLDataRPCServer: load CA TSL authentcation files success")

	// get an interceptor server option for Unary-Unary handler
	unaryOption := getUnaryInterceptorOption()
	log.Printf("[info] StartMySQLDataRPCServer: load unary method interceptor function success")

	server := grpc.NewServer(unaryOption, caOption)

	mysqlPb.RegisterMySQLBindServiceServer(server, &MySQLData{})

	listener, err := net.Listen("tcp", conf.MySQLDataRPCServerAddress)
	if nil != err {
		log.Fatalf("[error] StartMySQLDataRPCServer: %s", err.Error())
	}

	log.Printf("[info] StartMySQLDataRPCServer: start the server with tcp address %s", conf.MySQLDataRPCServerAddress)
	err = server.Serve(listener)
	if nil != err {
		log.Fatalf("[error] StartMySQLDataRPCServer: %s", err.Error())
	}
}

func StartMongoDataRPCServer() {
	// Start the gRPC server for MySQL data operation.

	// using CA TSL authentication
	caOption := getCAOption()

	// get an interceptor server option for Unary-Unary handler
	unaryOption := getUnaryInterceptorOption()

	server := grpc.NewServer(unaryOption, caOption)

	mongoPb.RegisterMongoBindServiceServer(server, &MongoData{})

	log.Println(":::Start Mongo Data Layer gRPC Server")
	listener, err := net.Listen("tcp", conf.MongoDataRPCServerAddress)
	if nil != err {
		log.Fatalf("Start gRPC server error: %s", err.Error())
	}
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}

}
