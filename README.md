## 项目启动说明
- ### 简介：
    项目为分布式项目,将服务拆分为`DataService`,`UserService`,`MessageService`等三个部分.
    - #### DataService
        提供数据服务, 当其他服务需要操作数据库(MySQL,MongoDB)时都要通过调用`DataService`提供的RPC服务(通过gRPC实现)
    - #### UserService
        提供用户相关服务, 主要包括登录注册, 个人信息管理, 好友管理等. 
        
        [HTTP接口文档](./UserService/Doc_HTTP_API.md)
        
        [MySQL相关表说明](./UserService/Doc_MySQL_Tabel.md)
        
    - #### MessageService
        提供聊天消息服务, 主要负责转发用户之间等一对一聊天和群聊, 系统通知, 订阅号新消息等, 使用`WebSocket`技术.