# UserService
- #### HTTP-API Functions  <span id="0"> </span>

  | Name       | Method | URL               | Auth Required | Description                    |
  | :--------- | :----- | :---------------- | :------------ | :----------------------------- |
  | [SignUp](#1)     | POST   | /auth/user        | 0             | 注册                           |
  | [SignIn](#2)     | POST   | /auth/profile     | 0             | 登录                           |
  | [GetResetPasswordEmail](#3) | GET | /auth/password | 0 | 忘记密码-发送修改链接邮件 |
  | [GetProfile](#4) | GET    | /info/profile     | 1             | 获取个人信息                   |
  | [PutProfile](#5) | PUT    | /info/profile     | 1             | 修改个人信息                   |
  | [PutAvatar](#6) | PUT | /info/avatar | 1 | 更新个人头像 |
  | [PutPassword](#7) | PUT | /info/password | 1 | 修改密码 |
  | [ForgetPassword](#8) | POST | /info/password | 1 | 忘记密码-重置密码 |
  | [ParseQrCode](#9) | POST | /info/qr_code | 1 | 解析上传的二维码 |
  |  |  |  |  |  |
  | [SearchUser](#9) | GET | /relation/users | 1 | 搜索用户 |
  | [AddFriend](#10) | POST   | /relation/friend  | 1             | 添加好友                       |
  | [PutFriend](#11) | PUT    | /relation/friend  | 1             | 修改好友备注; 接受\拒绝好友申请; 加入\移出黑名单 |
  | [DeleteFriend](#12) | DELETE | /relation/friend  | 1             | 删除好友                       |
  | [GetFriendsInfo](#13) | GET | /relation/friends | 1 | 获取好友列表 |
  |  |  |  |  |  |
  | 其他接口文档待完善..... |  |  |  |  |
  |  |  |  |  |  |

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
| password | string   | 8到12位个字符               | 密码           |

```json
{
	"email":"nihao@qq.com",
	"password":"xixiixi123"
}
```

##### Response: ⚠️**与SignUp的完全一致**

----

- #### <span id="3">GetResetPasswordEmail 忘记密码-发送修改链接邮件</span> [Top](#0)

##### Request: 

Path: `/info/password`	Method: `GET`

QueryStringParams: `查询字符串参数, 必传`

| Columns | DataType | Constraints                 | Description  |
| ------- | -------- | --------------------------- | ------------ |
| email   | string   | 符合邮箱格式,最多100个字符; | 注册的邮箱号 |

##### Response: 成功则只返回200状态码, 失败返回错误状态码及错误信息.

----

----

- #### <span id="4">GetProfile 获取用户个人详细信息</span> [Top](#0)

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

- #### <span id="5">PutProfile 修改用户的个人信息</span> [Top](#0)

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

##### Response: ⚠️**与GetProfile的完全一致**

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

---

- #### <span id="7">PutPassword 修改密码</span> [Top](#0)

##### Request:

Path: `/info/password`	Method: `PUT`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

​		    `Content-Type: application/json;`

JsonBodyParams: `所有参数必传`

| Columns          | DataType | Constraints   | Description |
| ---------------- | -------- | ------------- | ----------- |
| old_password     | string   | 8到12位个字符 | 旧密码      |
| password         | string   | --            | 新密码      |
| confirm_password | string   | --            | 确认密码    |

##### Response: 成功则只返回200状态码, 失败返回错误状态码及错误信息.

---

- #### <span id="8">ForgetPassword 忘记密码-重置密码</span> [Top](#0)

##### Request:

Path: `/info/password` 	Method: `POST`

Headers: `Auth-Token: "auth token value got from email"`

​		    `Content-Type: application/json;`

JsonBodyParams: `所有参数必传`

| Columns          | DataType | Constraints   | Description |
| ---------------- | -------- | ------------- | ----------- |
| password         | string   | 8到12位个字符 | 新密码      |
| confirm_password | string   | --            | 确认密码    |

##### Response: 成功则只返回200状态码, 失败返回错误状态码及错误信息.

---

- #### <span id="9"> ParseQrCode 解析上传的二维码</span> [Top](#0)

##### Requset:

Path: `/info/qr_code`		Method: `POST`

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

----

- #### <span id="9">SearchUsers 搜索用户/获取单个好友信息</span> [Top](#0) 

##### Request:

Path: `/relation/users`		Method: `GET`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

JsonBodyParams: `kw参数必传`

| Column | DataType | Constraints  | Descritption   |
| ------ | -------- | ------------ | -------------- |
| kw     | string   | 不为空字符串 | 搜索关键       |
| page   | int      | 非负数       | 页码           |
| size   | int      | 非负数       | 每页需要的数据 |

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult: `ElasticSearch提供的原始搜索结果, 当结果中max_score的值为0时, 表示没有匹配时搜索结果`

```json
{
    "took": 8,
    "timed_out": false,
    "_shards": {
        "total": 3,
        "successful": 3,
        "skipped": 0,
        "failed": 0
    },
    "hits": {
        "total": {
            "value": 4,
            "relation": "eq"
        },
        "max_score": 0.2706483,
        "hits": [
            {
                "_index": "private_im_user",
                "_type": "_doc",
                "_id": "1183956383159029784",
                "_score": 0.2706483,
                "_source": {
                    "id": 1183956383159029784,
                    "name": "wang王@123",
                    "email": "test@demo.com",
                    "avatar": "<test avatar url>",
                    "gender": "2",
                    "is_delete": false
                }
            },
            {
                "_index": "private_im_user",
                "_type": "_doc",
                "_id": "1183956383159029780",
                "_score": 0,
                "_source": {
                    "id": 1183956383159029780,
                    "name": "cpp",
                    "email": "test@qq.com",
                    "avatar": "<avatar pic url>",
                    "gender": 1,
                    "is_delete": false
                }
            },
            {
                "_index": "private_im_user",
                "_type": "_doc",
                "_id": "1183956383159029786",
                "_score": 0,
                "_source": {
                    "id": 1183956383159029786,
                    "name": "golang",
                    "email": "python@qq.com",
                    "avatar": "<avatar pic url>",
                    "gender": 1,
                    "is_delete": false
                }
            },
            {
                "_index": "private_im_user",
                "_type": "_doc",
                "_id": "1183956383159029781",
                "_score": 0,
                "_source": {
                    "id": 1183956383159029781,
                    "name": "python",
                    "email": "demo@qq.com",
                    "avatar": "<avatar pic url>",
                    "gender": 1,
                    "is_delete": false
                }
            }
        ]
    }
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
    "message": "eg.: move friend into blacklist successful"
}
```

---

- #### <span id="12">DeleteFriend 删除好友</span> [Top](#0)

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
{
    "message": "the record has been deleted and will not be notified to your friend."
}
```

---

- #### <span id="13">GetFriendsInfo 获取自己的好友列表</span> [Top](#0)

##### Request:

Path: `/relation/friends` 	Method: `GET`

Headers: `Auth-Token: "auth token value from SignUp or SignIn"`

##### Response:

Headers: `Content-Type: application/json;`

JsonBodyResult:

| Column    | DataType | Description                        |
| --------- | -------- | ---------------------------------- |
| friend_id | int64    | 好友的Id                           |
| name      | string   | 好友的呢称                         |
| email     | string   | 好友的邮箱                         |
| mobile    | string   | 好友的手机号, 可能为空             |
| gender    | int      | 好友的性别(0:未知, 1:女性, 2:男性) |
| note      | string   | 给好友设置的备注                   |
| is_black  | bool     | 是否在黑名单内                     |

```json
{
    "friends": [
        {
            "id": 1162262948794597376,
            "note": "noteName",
            "email": "test@test.com",
            "name": "test",
            "mobile": "13122222223",
            "gender": 1,
            "avatar": "<avatar pic url>",
           	"is_accept": true,
            "is_black": false,
            "is_delete": false
        }
    ]
}
```

