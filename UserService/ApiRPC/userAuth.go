package ApiRPC

import (
	"context"

	"../Protos"
	"../utils"

	conf "../Config"
)

type UserAuth struct {
}

func (obj *UserAuth) CheckAuthToken(ctx context.Context, param *userAuthPb.Token) (*userAuthPb.Id, error) {
	// parseToken
	claims, err := utils.ParseJWTToken(param.Value, []byte(conf.AuthTokenSalt))
	if err != nil {
		return nil, err
	}
	return &userAuthPb.Id{Value: claims.Id}, nil
}
