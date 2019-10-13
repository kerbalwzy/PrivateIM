package ApiRPC

import (
	msgNodesPb "../RpcClientPbs/msgNodesPb"
)

func MSGUserNodeAddFriend(selfId, friendId int64) {
	param := &msgNodesPb.DoubleId{MainId: selfId, OtherId: friendId}
	_, _ = GetMsgNodesDataClient().UserNodeFriendsAdd(getTimeOutCtx(3), param)
}

func MSGUserNodeDelFriend(selfId, friendId int64) {
	param := &msgNodesPb.DoubleId{MainId: selfId, OtherId: friendId}
	_, _ = GetMsgNodesDataClient().UserNodeFriendsDel(getTimeOutCtx(3), param)
}

func MSGUserNodeAddBlacklist(selfId, friendId int64) {
	param := &msgNodesPb.DoubleId{MainId: selfId, OtherId: friendId}
	_, _ = GetMsgNodesDataClient().UserNodeBlacklistAdd(getTimeOutCtx(3), param)
}

func MSGUserNodeMoveFriendIntoBlacklist(selfId, friendId int64) {
	param := &msgNodesPb.DoubleId{MainId: selfId, OtherId: friendId}
	_, _ = GetMsgNodesDataClient().UserNodeMoveFriendIntoBlacklist(getTimeOutCtx(3), param)
}

func MSGUserNodeMoveFriendOutFromBlacklist(selfId, friendId int64) {
	param := &msgNodesPb.DoubleId{MainId: selfId, OtherId: friendId}
	_, _ = GetMsgNodesDataClient().UserNodeMoveFriendOutFromBlacklist(getTimeOutCtx(3), param)
}
