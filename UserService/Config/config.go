package Config

const (
	MySQLDataRPCServerAddress = "0.0.0.0:23331"
	MongoDataRPCServerAddress = "0.0.0.0:23332"

	DataLayerSrvCAPem       = "/Users/wzy/GitPrograms/PrivateIM/DataLayerService/CATSL/ca.pem"
	DataLayerSrvCAClientPem = "/Users/wzy/GitPrograms/PrivateIM/DataLayerService/CATSL/client/client.pem"
	DataLayerSrvCAClientKey = "/Users/wzy/GitPrograms/PrivateIM/DataLayerService/CATSL/client/client.key"

	StaticFoldPath = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/static/"

	PasswordHashSalt   = "this is a password hash salt"
	AuthTokenSalt      = "this is a auth token salt"
	AuthTokenAliveTime = 3600 * 24 //unit:second
	AuthTokenIssuer    = "userCenter"

	PhotoSaveFoldPath   = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/static/photos/"
	PhotoSuffix         = ".png"
	PhotosUrlPrefix     = "/static/photos/" // if you use oss , should change this value
	DefaultAvatarUrl    = "/static/photos/defaultAvatar.jpg"
	AvatarPicUploadMaxSize = 100 * 2 << 10

	QRCodeBaseUrl = "http://127.0.0.1:8080/qrcontent/?"
)
