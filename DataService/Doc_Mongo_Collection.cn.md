## Project Mongo DB Collection Introduction [TOP](#0)

- #### Collections

  | Name                                | Description                                                  |
  | ----------------------------------- | ------------------------------------------------------------ |
  | [coll_delay_message](#1)            | 延时消息, 用户因未上线带接收的消息, 用户上线接收后则从集合删除 |
  | [coll_chat_history](2)              | 一对一用户聊天消息历史记录                                   |
  | [coll_group_chat_history](#3)       | 群聊消息历史记录                                             |
  | [coll_subscription_msg_history](#4) | 订阅号消息历史记录                                           |
  | [coll_user_friends](#5)             | 用户好友缓存数据;  UserService改变原数据的同时需要主动更新   |
  | [coll_group_users](#6)              | 群聊用户成员缓存数据; UserService改变原数据的同时需要主动更新 |
  | [coll_subscription_users](#7)       | 订阅用户成员缓存数据; UserService改变原数据的同时需要主动更新 |
  |                                     |                                                              |

----

- #### <span id="1">coll_delay_message 延时消息</span> [TOP](#0)

  | Name            | BsonType   | GolangType | Description                |
  | --------------- | ---------- | ---------- | -------------------------- |
  | _id             | NumberLong | int64      | 用户ID, 雪花算法生成的数值 |
  | message         | Array      | slice      | 待发送消息的列表           |
  | message-element | BinData    | []byte     | 消息字符串                 |

  ```reStructuredText
  {	// 单条数据结构
      "_id": NumberLong(1234567890123456789),
      "message" [
      	BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	]
  }
  ```

---

- #### <span id="2">coll_chat_history 一对一用户聊天记录</span> [TOP](#0)

  | Name | BsonType | GolangType | Description |
  | ---- | -------- | ---------- | ----------- |
  |      |          |            |             |
  |      |          |            |             |
  |      |          |            |             |
