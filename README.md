## 项目启动说明
- ### 简介：
    项目为分布式项目,将服务拆分为`DataService`,`UserService`,`MessageService`等三个部分.
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
        + 群聊新建、加入、退出 \ 获取、修改群聊信息 
        + 订阅号关注、取消关注 

        [HTTP接口文档](./UserService/Doc_HTTP_API.cn.md)

        [gRPC接口文档](./UserService/Doc_gRPC_API.cn.md)

    - #### MessageService
        提供聊天消息服务, 主要负责转发用户之间等一对一聊天和群聊, 系统通知, 订阅号新消息等, 使用`WebSocket`技术.