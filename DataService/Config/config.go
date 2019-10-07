package Config

const (
	MySQLURI = "root:123456@tcp(127.0.0.1:3306)/PrivateIM?charset=utf8&parseTime=true"

	MongoDBURI                     = "mongodb://localhost:27017"
	MongoDBName                    = "privateIM"
	CollDelayMessageName           = "delay_message"
	CollChatHistoryName            = "user_chat_history"
	CollGroupChatHistoryName       = "group_chat_history"
	CollSubscriptionMsgHistoryName = "subscription_msg_history"
	CollUserFriendsName            = "user_friends"
	CollUserGroupChatsName         = "user_group_chats"
	CollUserSubscriptions          = "user_subscriptions"
	CollGroupChatUsersName         = "group_chat_users"
	CollSubscriptionUsersName      = "subscription_users"

	MySQLDataRPCServerAddress = "0.0.0.0:23331"
	MongoDataRPCServerAddress = "0.0.0.0:23332"

	// todo: when the file path changed, the config should be change
	DataLayerSrvCAPem       = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/ca.pem"
	DataLayerSrvCAServerPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/server/server.pem"
	DataLayerSrvCAServerKey = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/server/server.key"
	DataLayerSrvCAClientPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/client/client.pem"
	DataLayerSrvCAClientKey = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/client/client.key"
)
