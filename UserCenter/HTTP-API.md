# UserCenter
- #### HTTP-API Functions 

  | Name       | Method | URL          | Auth Require | Description        |
  | :--------- | :----- | :----------- | :----------- | :----------------- |
  | SignUp     | POST   | /user        | 0            | 注册               |
  | SignIn     | POST   | /profile     | 0            | 登录               |
  | SignOut    | DELETE | /profile     | 1            | 退出登录           |
  | GetProfile | GET    | /profile     | 1            | 获取个人信息       |
  | ——         | ——     | /profile/:id | 1            | 获取好友信息       |
  | PutProfile | PUT    | /profile     | 1            | 修改个人信息       |
  | ——         | ——     | /profile/:id | 1            | 修改好友备注和标签 |
  | GetFriends | GET    | /firend      | 1            | 获取好友列表       |
  | AddFriend  | POST   | /firend      | 1            | 添加好友           |
  | DelFriend  | DELETE | /firend      | 1            | 删除好友           |

----

- #### SignUp 注册

  Path:  `/user`		Method: `POST`	 `  

  Params: `Content-Type: Application/json;`

  | Column           | DataType | Constraints                     | Description    |
  | ---------------- | -------- | ------------------------------- | -------------- |
  | name             | string   | 1到10个字符                     | 用户昵称       |
  | email            | string   | 符合邮箱格式,最多100个字符      | 注册的邮箱地址 |
  | password         | string   | 8到12位个符                     | 密码           |
  | confirm_password | string   | 8到12位个符, 与password完全一样 | 确认密码       |

  Return: `Content-Type: Application/json;`

  | Column      | DataType | Description                                                  |
  | ----------- | -------- | ------------------------------------------------------------ |
  | AuthToken   | string   | 认证Token字符串, 原理与JWT类似, 但更加简易.                  |
  | id          | int64    | 用户ID                                                       |
  | name        | string   | 用户昵称                                                     |
  | mobile      | string   | 手机号, 默认为空                                             |
  | email       | string   | 注册邮箱号                                                   |
  | gender      | int      | 性别, 默认(-1: 未知) ;    (0: 女)    (1: 男)                 |
  | create_time | string   | 账号创建时间 (格式: "2006-01-02T15:04:05Z") [**rfc3339格式**] |
  | update_time | string   | 账号更新时间 (格式: "2006-01-02T15:04:05Z") [**rfc3339格式**] |

- #### SignIn 登录

  Path: `/profile`		Method: `POST`

  Params: `Content-Type: Application/json;`

  | Column   | DataType | Constraints                | Description    |
  | -------- | -------- | -------------------------- | -------------- |
  | email    | string   | 符合邮箱格式,最多100个字符 | 注册的邮箱地址 |
  | password | string   | 8到12位个符                | 密码           |

  Return: `Content-Type: Application/json;` 「⚠️**与SignUp的返回值完全一致**」





