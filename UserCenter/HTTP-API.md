# UserCenter
- #### HTTP-API Functions 

  | Name       | Method | URL               | Auth Required | Description                    |
  | :--------- | :----- | :---------------- | :------------ | :----------------------------- |
  | SignUp     | POST   | /auth/user        | 0             | 注册                           |
  | SignIn     | POST   | /auth/profile     | 0             | 登录                           |
  | GetProfile | GET    | /info/profile     | 1             | 获取个人信息                   |
  | PutProfile | PUT    | /info/profile     | 1             | 修改个人信息                   |
  | GetQrCode  | GET    | /info/qrcode      | 1             | 获取个人二维码                 |
  | GetFriends | GET    | /relation/friends | 1             | 获取好友列表                   |
  | GetFriend  | GET    | /relation/friend  | 1             | 获取单个好友信息               |
  | AddFriend  | POST   | /relation/friend  | 1             | 添加好友                       |
  | PutFriend  | PUT    | /relation/friend  | 1             | 修改好友备注; 加入、移出黑名单 |
  | DelFriend  | DELETE | /relation/friend  | 1             | 删除好友                       |

----

- #### SignUp 注册

##### Request: 

Path:  `/auth/user`		Method: `POST`

Headers: `Content-Type: application/json;`

JsonBodyParams:  `所有参数均为必传`

| Column           | DataType | Constraints                       | Description    |
| ---------------- | -------- | --------------------------------- | -------------- |
| name             | string   | 1到10个字符                       | 用户昵称       |
| email            | string   | 符合邮箱格式,最多100个字符        | 注册的邮箱地址 |
| password         | string   | 8到12位个字符                     | 密码           |
| confirm_password | string   | 8到12位个字符, 与password完全一样 | 确认密码       |

```json
{
	"name":"test",
	"email":"test@test.com",
	"password":"12345678",
	"confirm_password":"12345678"
}
```

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column      | DataType | Description                                                  |
| ----------- | -------- | ------------------------------------------------------------ |
| AuthToken   | string   | 认证Token字符串, 原理与JWT类似, 但更加简易.                  |
| id          | int64    | 用户ID                                                       |
| name        | string   | 用户昵称                                                     |
| mobile      | string   | 手机号, 默认为空                                             |
| email       | string   | 注册邮箱号                                                   |
| gender      | int      | 性别, 默认(-1: 未知) ;    (0: 女)    (1: 男)                 |
| create_time | string   | 账号创建时间 (格式: "2006-01-02T15:04:05Z") [**rfc3339格式**] |

  ```json
  {
      "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMTYwMDAxODE5ODg0MTI2MjA4LCJleHAiOjE1NjU0MDMxODQsImlzcyI6InVzZXJDZW50ZXIifQ.TYrJufslBSiDFZ23lH6TuwHuuORGEVXv73HD_rh0sZQ",
      "user": {
          "id": 1160001819884126208,
          "name": "test",
          "mobile": "",
          "email": "test@test.com",
          "gender": -1,
          "create_time": "2019-08-09T10:51:29Z"
      }
  }
  ```

----

- #### SignIn 登录

##### Request:

Path: `/auth/profile`		Method: `POST`

Headers: `Content-Type: application/json;`

JsonBodyParams: `所有参数均为必传`

| Column   | DataType | Constraints                 | Description    |
| -------- | -------- | --------------------------- | -------------- |
| email    | string   | 符合邮箱格式,最多100个字符; | 注册的邮箱地址 |
| password | string   | 8到12位个符                 | 密码           |

```json
{
	"email":"nihao@qq.com",
	"password":"xixiixi123"
}
```

##### Response: ⚠️**与SignUp的完全一致**

----

- #### GetProfile 获取用户个人详细信息

##### Request:

Path: `/info/profile`		Method: `GET`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column | DataType   | Description |
| ------ | ---------- | ----------- |
| id          | int64    | 用户ID                                                       |
| name        | string   | 用户昵称                                                     |
| mobile      | string   | 手机号, 默认为空                                             |
| email       | string   | 注册邮箱号                                                   |
| gender      | int      | 性别, 默认(-1: 未知) ;    (0: 女)    (1: 男)                 |
| create_time | string   | 账号创建时间 (格式: "2006-01-02T15:04:05Z") [**rfc3339格式**] |

```json
{
    "id": 1160001819884126208,
    "name": "test",
    "mobile": "",
    "email": "test@test.com",
    "gender": -1,
    "create_time": "2019-08-09T10:51:29Z"
}
```

----

- #### PutProfile 修改用户的个人信息

##### Request:

Path: `/info/profile`		Method: `PUT`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

JsonBodyParams: `所有参数为必填, 如果未发生改变则填写原值`

| Columns | DataType | Constraints                                  | Descripton |
| ------- | -------- | -------------------------------------------- | ---------- |
| name    | string   | 1到10个字符                                  | 用户昵称   |
| mobile  | string   | 0个或者11个数字字符                          | 用户手机号 |
| gender  | int      | -1/0/1;  (-1: 未知) ;    (0: 女)    (1: 男); | 性别       |

##### Response: 

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column | DataType | Description                         |
| ------ | -------- | ----------------------------------- |
| name   | string   | 用户昵称                            |
| mobile | string   | 手机号, 默认为空                    |
| gender | int      | 性别 (-1: 未知)   (0: 女)   (1: 男) |

```json
{
    "name": "newName",
    "mobile": "13122222221",
    "gender": 1
}
```

---
