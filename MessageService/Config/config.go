package Config

const (
	MessageServerAddress        = "0.0.0.0:8080"
	MongoDataRPCServerAddress   = "0.0.0.0:23332"
	UserAuthRPCServerAddress    = "0.0.0.0:11111"
	MsgTransferRPCServerAddress = "0.0.0.0:12222"

	PrivateIMRootCAPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/ca.pem"
	PrivateIMServerPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/server/server.pem"
	PrivateIMServerKey = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/server/server.key"
	PrivateIMClientPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/client/client.pem"
	PrivateIMClientKey = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/client/client.key"

	GroupChatNodeLifeTime                   = 60 * 60 * 24 * 3 // unit: sec
	GroupChatNodeLowActivityCleanPercentage = 30               // 0-100 / 100 %
	GroupChatNodeCleanTime                  = 3                // 0-23 h/d (every day)

	// when the count of group chat nodes is less then this value, would not do the clear up by activity count.
	GroupChatNodeLowActivityCleanStartLimit = 1000

	SubscriptionNodeLifeTime  = 60 * 60 * 23 // unit: sec
	SubscriptionNodeCleanTime = 3            // 0-23 h/d (every day)

	CleanWorkRPCCallSecretKey = "1234!@#$%^&*nsb#$%^78qwe"
)
