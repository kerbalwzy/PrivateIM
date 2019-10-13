package ApiRPC

import (
	"context"
	"errors"

	"../MSGNode"

	conf "../Config"
	msgNodesPb "../Protos"
)

type NodesData struct {
}

func (obj *NodesData) UserNodeFriendsAdd(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalUsers.Get(param.MainId); ok {
		node.Friends.Add(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) UserNodeFriendsDel(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalUsers.Get(param.MainId); ok {
		node.Friends.Del(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) UserNodeBlacklistAdd(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalUsers.Get(param.MainId); ok {
		node.BlackList.Add(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) UserNodeBlacklistDel(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalUsers.Get(param.MainId); ok {
		node.BlackList.Del(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) UserNodeMoveFriendIntoBlacklist(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalUsers.Get(param.MainId); ok {
		node.Friends.Del(param.OtherId)
		node.BlackList.Add(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) UserNodeMoveFriendOutFromBlacklist(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalUsers.Get(param.MainId); ok {
		node.Friends.Add(param.OtherId)
		node.BlackList.Del(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) GroupChatNodeUserAdd(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalGroupChats.Get(param.MainId); ok {
		node.Users.Add(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) GroupChatNodeUserDel(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalGroupChats.Get(param.MainId); ok {
		node.Users.Del(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) SubsNodeFansAdd(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalSubscriptions.Get(param.MainId); ok {
		node.Fans.Add(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) SubsNodeFansDel(ctx context.Context, param *msgNodesPb.DoubleId) (
	*msgNodesPb.Result, error) {
	if node, ok := MSGNode.GlobalSubscriptions.Get(param.MainId); ok {
		node.Fans.Del(param.OtherId)
	}
	return &msgNodesPb.Result{Code: 200}, nil
}

var ErrWrongSecretKey = errors.New("the secret key is wrong")

func (obj *NodesData) GroupChatNodesPoolCleanByLifeTime(ctx context.Context, param *msgNodesPb.CleanNodePoolAuth) (
	*msgNodesPb.Result, error) {
	if param.SecretKey != conf.CleanWorkRPCCallSecretKey {
		return &msgNodesPb.Result{Code: 400}, ErrWrongSecretKey
	}
	MSGNode.GlobalGroupChats.CleanByLifeTime()
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) GroupChatNodesPoolCleanByActivity(ctx context.Context, param *msgNodesPb.CleanNodePoolAuth) (
	*msgNodesPb.Result, error) {
	if param.SecretKey != conf.CleanWorkRPCCallSecretKey {
		return &msgNodesPb.Result{Code: 400}, ErrWrongSecretKey
	}
	MSGNode.GlobalGroupChats.CleanByActiveCount()
	return &msgNodesPb.Result{Code: 200}, nil
}

func (obj *NodesData) SubsNodesPoolCleanByLifeTime(ctx context.Context, param *msgNodesPb.CleanNodePoolAuth) (
	*msgNodesPb.Result, error) {
	if param.SecretKey != conf.CleanWorkRPCCallSecretKey {
		return &msgNodesPb.Result{Code: 400}, ErrWrongSecretKey
	}
	MSGNode.GlobalSubscriptions.CleanByLifeTime()
	return &msgNodesPb.Result{Code: 200}, nil
}
