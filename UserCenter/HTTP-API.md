# UserCenter
- #### HTTP-API Functions  <span id="0"> </span>

  | Name       | Method | URL               | Auth Required | Description                    |
  | :--------- | :----- | :---------------- | :------------ | :----------------------------- |
  | [SignUp](#1)     | POST   | /auth/user        | 0             | 注册                           |
  | [SignIn](#2)     | POST   | /auth/profile     | 0             | 登录                           |
  | [GetProfile](#3) | GET    | /info/profile     | 1             | 获取个人信息                   |
  | [PutProfile](#4) | PUT    | /info/profile     | 1             | 修改个人信息                   |
  | [GetAvatar](#5) | GET    | /info/avatar      | 1             | 获取个人头像                   |
  | [PutAvatar](#6)  | PUT    | /info/avatar      | 1             | 更新个人头像                   |
  | [GetQrCode](#7) | GET    | /info/qrcode      | 1             | 获取个人二维码                 |
  | [ParseQrCode](#8) | POST | /info/qrcode | 1 | 解析上传的二维码 |
  | [GetFriend](#9) | GET | /relation/friend | 1 | 搜索用户/获取好友信息 |
  | [AddFriend](#10) | POST   | /relation/friend  | 1             | 添加好友                       |
  | [PutFriend](#11) | PUT    | /relation/friend  | 1             | 修改好友备注; 接受\拒绝好友申请; 加入\移出黑名单 |
  | [DelFriend](#12) | DELETE | /relation/friend  | 1             | 删除好友                       |
  | [GetFriends] | GET | /relation/friends | 1 | 获取好友列表 |

----

- #### <span id="1">SignUp 注册</span> [Top](#0)

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

- #### <span id="2">SignIn 登录</span> [Top](#0)

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

- #### <span id="3">GetProfile 获取用户个人详细信息</span> [Top](#0)

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

- #### <span id="4">PutProfile 修改用户的个人信息</span> [Top](#0)

##### Request:

Path: `/info/profile`		Method: `PUT`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​		    `Content-Type: application/json;`

JsonBodyParams: `所有参数为必填, 如果未发生改变则填写原值`

| Columns | DataType | Constraints                                | Descripton |
| ------- | -------- | ------------------------------------------ | ---------- |
| name    | string   | 1到10个字符                                | 用户昵称   |
| mobile  | string   | 0个或者11个数字字符                        | 用户手机号 |
| gender  | int      | 0/1/2;  (0: 未知) ;    (1: 女)    (2: 男); | 性别       |

##### Response: 

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column | DataType | Description                        |
| ------ | -------- | ---------------------------------- |
| name   | string   | 用户昵称                           |
| mobile | string   | 手机号, 默认为空                   |
| gender | int      | 性别 (0: 未知)   (1: 女)   (2: 男) |

```json
{
    "name": "newName",
    "mobile": "13122222221",
    "gender": 1
}
```
---

- #### <span id="5">GetAvatar 获取用户头像</span> [Top](#0)

##### Request:

Path: `/info/avatar`		Method: `GET`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

##### Response: 

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column     | DataType | Description  |
| ---------- | -------- | ------------ |
| avatar_url | string   | 用户头像地址 |

```json
{
    "avatar_url": "this is the avatar url"
}
```

----

- #### <span id="6">PutAvatar 更新用户头像 </span> [Top](#0)

##### Request: 

Path: `/info/avatar`		Method: `PUT`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​                 `Content-Type: multipart/form-data"`

FormBodyParams: `所有参数为必填`

| Column     | DataType | Constraints                                                  | Description  |
| ---------- | -------- | ------------------------------------------------------------ | ------------ |
| new_avatar | file     | 图片文件,上传前进行裁剪压缩处理, 保证图片尺寸为618*618像素, 分辨率80, 最大体积100kb. | 用户头像文件 |

##### Response: ⚠️ 与GetAvatar完全一致

---

- #### <span id="7"> GetQrCode 获取个人二维码</span> [Top](#0)

##### Request:

Path: `/info/qrcode`		Method: `GET`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column      | DataType | Description    |
| ----------- | -------- | -------------- |
| qr_code_url | string   | 二维码链接地址 |

```json
{
    "qr_code_url": "this is the QRCode url"
}
```

---

- #### <span id="8"> ParseQrCode 解析上传的二维码</span> [Top](#0)

##### Requset:

Path: `/info/qrcode`		Method: `POST`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​                 `Content-Type: multipart/form-data"`

FormBodyParams: `所有参数为必填`

| Column  | DataType | Constraints                                                  | Description        |
| ------- | -------- | ------------------------------------------------------------ | ------------------ |
| qr_code | file     | 图片文件,上传前进行裁剪压缩处理, 图片内容只能包含纯二维码区域, 分辨率80, 最大体积100kb. | 截取到的二维码图片 |

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column     | DataType | Description          |
| ---------- | -------- | -------------------- |
| qr_content | string   | 二维码包含的真实信息 |

```json
{
    "qr_content": "the real content of qr code"
}
```

---

- #### <span id="9">GetFriend 搜索用户/获取单个好友信息</span> [Top](#0) 

##### Request:

Path: `/relation/friend`		Method: `GET`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

JsonBodyParams: `至少包含三个参数中的一个, 参数优先级:id > email > name`

| Column | DataType | Constraints                 | Descritption       |
| ------ | -------- | --------------------------- | ------------------ |
| id     | int64    | 无符号整型                  | 目标用户的ID       |
| email  | string   | 符合邮箱格式,最多100个字符; | 目标用户的邮箱地址 |
| name   | string   | 1到10个字符                 | 目标用户的昵称     |

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult: `id, mobile, gender, note只当两个用户存在有效好友关系,且好友保存了相关信息才会返回有效值, 否则都返回数据类型的默认值; 使用name作为搜索参数时可能返回多条数据`

| Column | DataType | Description      |
| ------ | -------- | ---------------- |
| id     | int64    | 目标用户ID       |
| email  | string   | 目标用户邮箱地址 |
| name   | string   | 目标用户昵称     |
| mobile | string   | 目标用户手机号   |
| gender | int      | 目标用户的性别   |
| note   | string   | 给好友备注的名称 |

```json
{	// the demo result is search by name = "test"
    "result": [
        {	// friend
            "id": 1162262948794597376,
            "email": "test@test.com",
            "name": "test",
            "mobile": "13122222223",
            "gender": 1,
            "note": "Li"
        },
        {	// not friend
            "id": 1162663753959866368,
            "email": "demo@demo.com",
            "name": "test",
            "mobile": "",
            "gender": 0,
            "note": ""
        }
        ....
    ]
}
```

---

- #### <span id="10">AddFriend 添加好友</span> [Top](#0)

##### Requset:

Path: `/relation/friend`		Method: `POST`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​		    `Content-Type: application/json;`

JsonBodyParams: `所有参数均为必传`

| Column | DataType | Constraints                                 | Description          |
| ------ | -------- | ------------------------------------------- | -------------------- |
| dst_id | int64    | 无符号整型                                  | 目标用户的ID         |
| note   | string   | 最多10个字符,可以为空字符串, 表示不设置昵称 | 给目标用户设置的昵称 |

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column  | DataType | Description        |
| ------- | -------- | ------------------ |
| message | string   | 操作结果提示字符串 |

```json
{
    "message": "initiate and add friends successfully, wait for the target user to agree"
}
```

---

- #### <span id="11">PutFriend 修改好友备注; 接受\拒绝好友申请; 加入\移出黑名单</span> [Top](#0)

##### Request:

Path: `/relation/friend`		Method: `PUT`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​		    `Content-Type: application/json;`

JsonBodyParams: `在拒绝好友申请的同时,会将用户加入黑名单;`

| Column    | DataType | Constraints                                                  | Description      |
| --------- | -------- | ------------------------------------------------------------ | ---------------- |
| action    | int      | 1/2/3 (1:修改备注 2:接受/拒绝好友申请 3:加入/移除黑名单 )[必填] | 要执行的操作     |
| dst_id    | int64    | 无符号整型[必填]                                             | 目标用户的ID     |
| note      | string   | 小于10个字符[action为1时必填, action为2时选填]               | 给目标用户的备注 |
| is_accept | bool     | 默认为false[action为2时必填]                                 | 是否接受好友申请 |
| is_black  | bool     | 默认为false[action为3时必填]                                 | 是否加入黑名单   |

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column  | DataType | Description                                               |
| ------- | -------- | --------------------------------------------------------- |
| action  | int      | 1/2/3 (1:修改备注 2:接受/拒绝好友申请 3:加入/移除黑名单 ) |
| message | string   | 执行结果描述                                              |

```json
{
    "action": 3,
    "message": "move friend into blacklist successful"
}
```

---

- #### <span id="12">DelFriend 删除好友</span> [Top](#0)

##### Request:

Path: `/relation/friend` 		Method: `DELETE`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​		    `Content-Type: application/json;`

JsonBodyParams: `在删除好友的同时,会将好友关系记录中的is_accept, is_black都重置为false; 除非自己重新请求添加对方为好友,或者对象重新请求添加好友,删除状态不可更改`

| Column | DataType | Constraints      | Description  |
| ------ | -------- | ---------------- | ------------ |
| dst_id | int64    | 无符号整型[必填] | 目标用户的ID |

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column  | DataType | Description  |
| ------- | -------- | ------------ |
| message | string   | 操作结果提示 |

```json

```

