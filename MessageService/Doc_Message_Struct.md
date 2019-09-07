# Type and struct of message

- #### MessageType

  ```go
  // The type code of message
  const (
  	NormalMessage  = 0 // Users chat with each other one to one
  	GroupMessage   = 1 // User group chat
  	ChannelMessage = 2 // From system notification or user subscription
  	DebugMessage   = 3 // Tell the client what error happened
  )
  ```

- #### 

  ----

- #### Struct

  ```go
  // Basic message struct
  type BasicMessage struct {
  	TypeId      int   `json:"type_id"` // the type number of message
  	SrcId       int64 `json:"src_id"`  // who send this message, the sender id
  	DstId       int64 `json:"dst_id"`  // who will recv this message, the receiver id
  	ProduceTime int64 `json:"produce_time,omitempty"` // the message produce time
  }
  ```
