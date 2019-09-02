## 后端逻辑说明(Logic note of backend)

- #### 接收到客户端的WebSocket升级请求(recv a WebSocket upgrade request of  client)

  ```reStructuredText
  1. 检查请求参数中的Token值(check the token which in request params)
  通过gRPC调用将Token值发送到用户中心检查, 并获取返回的结果, 如果检查结果为真, 则将连接升级, 并为该用户的建立一个客户端节点, 否则返回失败信息给前端.
  Check the token by send it to UserCenter through gRPC call, and get the response. if the result is true, then upgrade the connection and new a client node for this user, else return an error information to client.
  ```

  ```reStructuredText
  2. 检查用户是否有待接收的数据(check if there are waiting recv messages for the user)
  
  ```
