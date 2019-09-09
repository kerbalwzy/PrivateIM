package ApiRPC

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"

	"../utils"

	conf "../Config"
	pb "../Protos"
)

type AuthServer struct{}

func (obj *AuthServer) CheckAuthToken(ctx context.Context, params *pb.AuthToken) (*pb.TokenCheckResult, error) {
	claims, err := utils.ParseJWTToken(params.Data, []byte(conf.AuthTokenSalt))
	if nil != err {
		return nil, err
	}
	return &pb.TokenCheckResult{UserId: claims.Id}, nil
}

var CtxCanceledErr = errors.New("the client canceled or connection time out")

// Check the client if canceled the calling or connection is time out
func checkCtxCanceled(ctx context.Context) error {
	if ctx.Err() == context.Canceled {
		return CtxCanceledErr
	}
	return nil
}

func getUnaryInterceptorOption() grpc.ServerOption {
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp interface{}, err error) {

		err = checkCtxCanceled(ctx)
		if nil != err {
			return
		}

		return handler(ctx, req)
	}
	return grpc.UnaryInterceptor(interceptor)
}

func StartUserAuthRPCServer() {
	listener, err := net.Listen("tcp", conf.UserAuthRPCServerAddress)
	if nil != err {
		log.Fatal(err)
	}

	unaryOption := getUnaryInterceptorOption()
	server := grpc.NewServer(unaryOption)

	pb.RegisterUserAuthServer(server, &AuthServer{})
	log.Println(":::Start UserAuth gRPC Server")
	err = server.Serve(listener)
	if nil != err {
		log.Fatal(err)
	}

}
