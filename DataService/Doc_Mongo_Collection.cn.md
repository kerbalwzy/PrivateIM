## <span id="0">Project Mongo DB Collection Introduction </span>

- #### Collections

  | Name                                | Description                                                  |
  | ----------------------------------- | ------------------------------------------------------------ |
  | [coll_delay_message](#1)            | 延时消息; 用户因未上线待接收的消息, 用户上线接收后则从集合删除 |
  | [coll_chat_history](#2)             | 一对一用户聊天消息历史记录                                   |
  | [coll_group_chat_history](#3)       | 群聊消息历史记录                                             |
  | [coll_subscription_msg_history](#4) | 订阅号消息历史记录                                           |
  | [coll_user_friends](#5)             | 用户的好友-缓存数据;  MySQL原数据被改变时需要同步更新        |
  | [coll_group_users](#6)              | 群聊的用户成员-缓存数据; MySQL原数据被改变时需要同步更新     |
  | [coll_subscription_users](#7)       | 订阅号的用户成员-缓存数据; MySQL原数据被改变时需要同步更新   |
  | [coll_user_group_chat](#8)          | 用户的加入的群聊-缓存数据; MySQL原数据被改变时需要同步更新   |
  | [coll_user_subscription](#9)        | 用户关注的订阅号-缓存数据; MySQL原数据被改变时需要同步更新   |

----

- #### <span id="1">coll_delay_message 延时消息</span> [TOP](#0)

  | Name            | BsonType   | GolangType | Description      |
  | --------------- | ---------- | ---------- | ---------------- |
  | _id             | NumberLong | int64      | 用户ID           |
  | message         | Array      | slice      | 待发送消息的列表 |
  | message-element | BinData    | []byte     | 消息字符串       |

  ```reStructuredText
  {	// 单条数据结构
      "_id": NumberLong(1234567890123456789),
      "message": [
          BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
          BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
      ]
  }
  ```

---

- #### <span id="2">coll_chat_history 一对一用户聊天记录</span> [TOP](#0)

  | Name            | BsonType | GolangType | Description                                    |
  | --------------- | -------- | ---------- | ---------------------------------------------- |
  | _id             | String   | String     | 聊天记录ID; 两个用户ID排序后拼接组成, 结果唯一 |
  | message_date    | Array    | slice      | 消息记录的列表; 键名中下划线后为具体日期       |
  | message-element | BinData  | []byte     | 消息字符串                                     |

  ```reStructuredText
  {	
  	// 下划线前的用户ID数值必定比下划线后的小
      "_id": "1234567890123456781_1234567890123456782",	
      // 每天的聊天记录存放一个列表
      "message_20190101": [
      	BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	],
  	// 如果当天有聊天记录才会记录, 例如20190102没有则不会有记录
  	"message_20190103": [
          BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	]
  }
  ```

---

- #### <span id="3">coll_group_chat_history 群聊消息历史记录</span> [TOP](#0)

  | Name            | BsonType   | GolangType | Description                              |
  | --------------- | ---------- | ---------- | ---------------------------------------- |
  | _id             | NumberLong | int64      | 群聊ID                                   |
  | message_date    | Array      | slice      | 消息记录的列表; 键名中下划线后为具体日期 |
  | message-element | BinData    | []byte     | 消息字符串                               |

  ```
  {
      "_id": NumberLong(1234567890123456999),
      // 每天的聊天记录存放一个列表
      "message_20190101": [
      	BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	],
  	// 如果当天有聊天记录才会记录, 例如20190102没有则不会有记录
  	"message_20190103": [
          BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	]
  }
  ```

---

- #### <span id="4">coll_subscription_msg_history 订阅号消息历史记录 </span> [TOP](#0)

  | Name            | BsonType   | GolangType | Description                              |
  | --------------- | ---------- | ---------- | ---------------------------------------- |
  | _id             | NumberLong | int64      | 订阅号ID                                 |
  | message_date    | Array      | slice      | 消息记录的列表; 键名中下划线后为具体日期 |
  | message-element | BinData    | []byte     | 消息字符串                               |

  ```
  {
      "_id": NumberLong(1234567890123456999),
      // 每天的消息记录存放一个列表
      "message_20190101": [
      	BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	],
  	// 如果当天有消息记录才会记录, 例如20190102没有则不会有记录
  	"message_20190103": [
          BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="), 						
  		BinData(0,"dGhpcyBpcyBhIHRlc3QgbWVzc2FnZQ=="),
  	]
  }
  ```

---

- #### <span id="5">coll_user_firends 用户的好友-缓存数据 </span> [TOP](#0)

  | Name              | BsonType   | GolangType | Description |
  | ----------------- | ---------- | ---------- | ----------- |
  | _id               | NumberLong | int64      | 用户ID      |
  | friends           | Array      | slice      | 好友ID列表  |
  | friends-element   | NumberLong | int64      | 其他用户ID  |
  | blacklist         | Array      | slice      | 黑名单列表  |
  | blacklist-element | NumberLong | int64      | 其他用户ID  |

  ```
  {
      "_id": NumberLong(1234567890123456789),
      // 好友ID列表
      "friends": [
          NumberLong(1234567890123456790),
          NumberLong(1234567890123456791),
      ],
      // 黑名单ID列表
      "blacklist": [
      	NumberLong(1234567890123456792),
          NumberLong(1234567890123456793),
      ]
  }
  ```

---

- #### <span id="6">coll_group_users 群聊的用户成员-缓存数据</span> [TOP](#0)

  | Name          | BsonType   | GolangType | Description          |
  | ------------- | ---------- | ---------- | -------------------- |
  | _id           | NumberLong | int64      | 群聊ID               |
  | users         | Array      | slice      | 群聊的用户成员ID列表 |
  | users-element | NumberLong | int64      | 用户ID               |

  ```
  {
      "_id": NumberLong(1234567890123458888),
      // 成员用户ID列表
      "users": [
          NumberLong(1234567890123456790),
          NumberLong(1234567890123456791),
      ]
  }
  ```

---

- #### <span id="7">coll_subscription_users 订阅号的用户成员-缓存数据</span> [TOP](#0)

  | Name          | BsonType   | GolangType | Description            |
  | ------------- | ---------- | ---------- | ---------------------- |
  | _id           | NumberLong | int64      | 订阅号ID               |
  | users         | Array      | slice      | 订阅号的用户成员ID列表 |
  | users-element | NumberLong | int64      | 用户ID                 |

  ```
  {
      "_id": NumberLong(1234567890123457777),
      // 成员用户ID列表
      "users": [
          NumberLong(1234567890123456790),
          NumberLong(1234567890123456791),
      ]
  }
  ```

---

- #### <span id="8">coll_user_group_chat 用户的加入的群聊-缓存数据</span> [TOP](#0)

  | Name           | BsonType   | GolangType | Description          |
  | -------------- | ---------- | ---------- | -------------------- |
  | _id            | NumberLong | int64      | 用户ID               |
  | groups         | Array      | slice      | 用户加入的群聊ID列表 |
  | groups-element | NumberLong | int64      | 群聊ID               |

  ```
  {
      "_id": NumberLong(1234567890123456790),
      // 用户加入的群聊ID列表
      "groups": [
          NumberLong(1234567890123458888),
          NumberLong(1234567890123458889),
      ]
  }
  ```

---

- #### <span id="9">coll_user_subscription 用户关注的订阅号-缓存数据</span> [TOP](#0)

  | Name                  | BsonType   | GolangType | Description            |
  | --------------------- | ---------- | ---------- | ---------------------- |
  | _id                   | NumberLong | int64      | 用户ID                 |
  | subscriptions         | Array      | slice      | 用户关注的订阅号ID列表 |
  | subscriptions-element | NumberLong | int64      | 订阅号ID               |

  ```
  {
      "_id":  NumberLong(1234567890123456790),
      // 用户关注的订阅号ID列表
      "subscriptions": [
          NumberLong(1234567890123457777),
          NumberLong(1234567890123457778),
      ]
  }
  ```
