package ApiRPC

import (
	"context"
	"errors"
	"log"

	"../ApiHTTP"
	"../utils"
)

var ClientWhiteList = map[string]string{
	"this is auth key for MsgCenter gRPC client": "127.0.0.1",
}

func checkClientInfo(info *ClientInfo) (bool, error) {
	host, ok := ClientWhiteList[info.AuthKey]
	if !ok {
		return ok, errors.New("client auth key validate fail")
	}
	if host != info.Host {
		return false, errors.New("client host not in white list")
	}
	return ok, nil
}

type MygRPCServer struct {
}

func (obj *MygRPCServer) CheckAuthToken(ctx context.Context, token *AuthToken) (*TokenCheckResult, error) {
	tokenString := token.Token
	clientInfo := token.ClientInfo
	log.Printf("gRPC server <CheckAuthToken> called by client(NAME %s HOST: %s )", clientInfo.Name, clientInfo.Host)
	isOk, err := checkClientInfo(clientInfo)
	if isOk {
		// parseToken
		_, err = utils.ParseJWTToken(tokenString, []byte(ApiHTTP.AuthTokenSalt))
		if nil != err {
			log.Printf("gRPC server <CheckAuthToken> Error:%s", err.Error())
			isOk = false
		}

	}
	return &TokenCheckResult{Ok: isOk}, err
}
