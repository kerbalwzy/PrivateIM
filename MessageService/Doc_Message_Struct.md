# Type and struct of message

- #### Message Type

  | ConstName                 | Value(int) | Description                    |
  | ------------------------- | ---------- | ------------------------------ |
  | UserChatMessageTypeId     | 0          | 普通的两个用户之间的聊天消息   |
  | GroupChatMessageTypeId    | 1          | 用户向群聊发送的消息           |
  | SubscriptionMessageTypeId | 2          | 订阅号向其订阅用户发送的消息   |
  | ErrorMessageTypeId        | 3          | 后台系统返回的消息传输结果信息 |

- #### Conten Type

  | ConstName    | Value(int) | Description      |
  | ------------ | ---------- | ---------------- |
  | TextContent  | 0          | 普通文字文本消息 |
  | ImageContent | 1          | 图片消息         |
  | VideoContent | 2          | 视频消息         |
  | VioceContent | 3          | 语音消息         |

----

- #### Message  Structure Note

  - ##### BasicMessage (the base message struct)

    ```go
    type BasicMessage struct {
    	TypeId       int   `json:"type_id"`                 // the type number of message
    	SenderId     int64 `json:"sender_id"`               // who send this message, the sender id
    	ReceiverId   int64 `json:"receiver_id"`             // who will recv this message, the receiver id
    	CreateTime   int64 `json:"create_time,omitempty"`   // set by sender or add by the message center, timestamp, unit:sec.
    	DeliveryTime int64 `json:"delivery_time,omitempty"` // the time for message want be sent, use for timing message
    }
    ```

  - ##### ChatMessage  (include the user and group chat message)

    ```go
    type ChatMessage struct {
    	BasicMessage
    	ContentType int    `json:"content_type"`           // how to show the message in client
    	Content     string `json:"content,omitempty"`      // text content
    	PreviewPic  string `json:"preview_pic,omitempty"`  // preview picture url
    	ResourceUrl string `json:"resource_url,omitempty"` // resource URL
    	Description string `json:"description,omitempty"`  // simple description
    	Additional  string `json:"additional,omitempty"`   // other additional information
        }
    ```

  - ##### SubscriptionMessage (the message which the subscription send to their every fans)

    ```go
    type SubscriptionMessage struct {
    	BasicMessage
    	Title       string `json:"title"`                  // the title
    	Abstract    string `json:"abstract,omitempty"`     // the brief introduction of this message
    	PreviewPic  string `json:"preview_pic,omitempty"`  // the preview picture url
    	ResourceUrl string `json:"resource_url,omitempty"` // resource URL
    }
    ```

  - ##### ErrorMessage (the message which the service back to the sender)

    ```go
    type ErrorMessage struct {
    	BasicMessage
    	Code       int    `json:"code"`                  // the code of error type
    	Error      string `json:"error"`                 // the detail error information
    	RawMessage []byte `json:"raw_message,omitempty"` // the row message which the user want to send.
    }
    ```

----

- #### Json String Message Demo

  #####   ⚠️: 这里出现的键值对都是必传的

  - Text Chat Message Demo 

    ```json
    {
        "type_id": 0, // or 1
        "sender_id": 0, 
        "receiver_id": 1, // can be user or group chat id
        "content_type": 0,
        "content": "<test group chat text message>"
    }
    ```

  - Media Chat Message Demo - Image Or Video

    ```json
    {
        "type_id": 0, // or 1
        "sender_id": 0, 
        "receiver_id": 1, // can be user or group chat id
        "content_type": 1, // or 2
        "preview_pic": "<preview picture link>",
        "resource_url": "<the raw raw resource link>"
    }
    ```

  - Media Chat Message Demo - Vioce

    ```json
    {
        "type_id": 0, // or 1
        "sender_id": 0, 
        "receiver_id": 1, // can be user or group chat id
        "content_type": 3,
        "resource_url": "<the raw raw resource link>"
    }
    ```

  - Subscription Message Demo 

    ```json
    {
        "type_id": 2,
        "sender_id": 1, // also the manager id of the subscription
        "receiver_id": 110,	// the subscription id
        "create_time": 12312312,
        "title": "<test title>",
        "abstract": "<test abstract>"
    }
    ```






