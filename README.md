## 项目说明

`gRPC`  `Proto3`  `CA`  `TLS`  `SnowFlake`  `MySQL`  `MongoDB`

`Gin`  `JWT`  `validator.v8`  `QrCode`  `Elasticsearch`  `gomail.v2`  `Redis`

`WebSocket` 

- ### 简介：
    项目将服务拆分为 **DataService**, **UserService**, **MessageService**, **SubscriptionService** 等四个部分.
    - #### DataService
        提供数据服务, 当其他服务需要操作数据库(MySQL,MongoDB)时都要通过调用`DataService`提供的RPC服务(通过gRPC实现).

        [MySQL相关表说明](./DataService/Doc_MySQL_TABEL.cn.md)

        [Mongo存储数据结构说明](./DataService/Doc_Mongo_Collection.cn.md)

        [Proto文件-操作MySQL数据的接口定义](./DataService/Protos/mysqlProto/mysqlBind.proto)

        [Proto文件-操作Mongo数据的接口定义](./DataService/Protos/mongoProto/mongoBind.proto)

        数据层服务的**RPC函数接口**应该简单而底层(对数据做基本CURD即可).  接口函数中不应该有业务逻辑判断的代码,  所有的业务相关代码都应该写在业务服务层, 这样才能让一个数据层服务能同时为多个(多种)业务层更好的提供数据服务. 

        虽然这样可能会因为在进行某个业务处理时需要多次请求数据层的服务造成性能的额外消耗,  但是为了让我们的数据层和业务层充分解耦, 能够更灵活的组合使用, 我认为这些消耗是值得的.

        如果是哪个具体的业务操作真的很需要优化性能, 并且性能瓶颈就是在于在业务处理中需要多次和数据层进行交互的话, 可以在数据层专门为该项业务定制接口, 在接口中可以包含些业务处理, 减少业务层和数据的交互次数以提高效率, 这样的函数接口在命名时应该显著的表现出这是一个为特定业务定制的RPC函数接口(本项目中为以**Plus**结尾函数). 

- #### UserService
    提供用户相关服务, 主要包括: 
    + 注册 \ 登录 \ 忘记密码 \ 发送重置密码邮件
    + 获取、修改基本信息 
    + 好友搜索、添加、删除、黑名单、修改备注
    + 群聊搜索、新建、加入、退出 \ 获取、修改群聊信息 
    + 订阅号搜索、关注、取消关注 

      [HTTP接口文档](./UserService/Doc_HTTP_API.cn.md)

      [Proto文件-用户Token检查](./UserService/Protos/UserAuth.proto)

- #### MessageService
    提供聊天消息服务, 主要负责转发用户之间等一对一聊天和群聊, 系统通知, 订阅号新消息等, 使用`WebSocket`技术.

    目前已经实现了:

    ​	**1. 同一用户ID最多可以在三个客户端登录后与服务建立连接, 最新的连接会插入到连接对象数组的0号位, 原其他位往后移动, 如果已经有数据元素已满, 则会弹出最后一个连接对象, 并由服务器主动关闭这个连接**

    ​	**2. 用户某一个客户端连接对象失效后自动从连接对象数据组中移出, 空位由后面的非空元素往前移动补位**

    ​	**3. 用户节点在所有其客户端断开链接后,自动从节点池中删除**

    ​	**4. 一对一用户聊天**

    ​	**5. 群聊节点每天定时按照存活时间和活跃度百分比自动清理**

    ​	**6. 群聊天**

    ​	**7. 订阅号节点每天定时按照存活时间自动清理**

    ​	**8. 订阅号向所有订阅者发送消息**

    ​	**9. 返回异常信息给发送者**

    ​	[消息类型与结构说明](./MessageService/Doc_Message_Struct.md)

    ​	[Proto文件-节点数据更新](./MessageService/Protos/MessageNodes.proto)

    ##### Future TODO:

    ​	实现集群功能, 能够同时启动多个MessageService, 并且保证用户的多点登录后即时Connecter分布在不同的服务中, 也能同步消息的收发.

- ### SubscriptionService (未开始)

    提供订阅号创建\管理, 文章发布, 审核, 订阅消息主动推送等服务; 

---

### Future:

​	使用微服务架构重构此项目, 使用功能划分更加合理, 提高并发能力, 和可用性.



