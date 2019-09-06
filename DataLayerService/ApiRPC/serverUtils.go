package ApiRPC

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"

	conf "../Config"
)

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
		log.Fatalf("tls.LoadX509KeyPair Error: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(conf.DataLayerSrvCAPem)
	if err != nil {
		log.Fatalf("Read DataLayerSrvCAPem Error: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("certPool.AppendCertsFromPEM Error")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	return grpc.Creds(c)
}

func recordLog() {

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

		// continue the handel
		return handler(ctx, req)
	}
	// add the interceptor for every Unary-Unary handler
	return grpc.UnaryInterceptor(interceptor)
}
