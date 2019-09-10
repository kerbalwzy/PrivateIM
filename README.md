## 项目启动说明
- ### 简介：
    项目为分布式项目,将服务拆分为`DataService`,`UserService`,`MessageService`等三个部分.
    - #### DataService
        提供数据服务, 当其他服务需要操作数据库(MySQL,MongoDB)时都要通过调用`DataService`提供的RPC服务(通过gRPC实现)
        
        [MySQL相关表说明](./DataService/Doc_MySQL_TABEL.cn.md)
        
        [Mongo存储数据结构说明](./DataService/Doc_Mongo_Collection.cn.md)
        
        [gRPC接口文档](./DataService/Doc_gRPC_API.cn.md)
        
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