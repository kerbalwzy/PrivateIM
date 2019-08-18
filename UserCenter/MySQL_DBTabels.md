# UserCenter

- #### Tables

  | Name               | Description  |
  | ------------------ | ------------ |
  | tb_user_basic      | 用户基础信息 |
  | tb_user_more       | 用户更多信息 |
  | tb_friend_relation | 好友关系表   |

---

- #### tb_user_basic

  | Column      | DataType     | Constraints               | Description                    |
  | ----------- | ------------ | ------------------------- | ------------------------------ |
  | id          | BIGINT(20)   | UNSIGNED PRIMARY KEY      | 主键ID, 雪花算法生成uint64数值 |
  | name        | VARCHAR(10)  | NOT NULL INDEX            | 昵称, 最多10个字符, 不允许为空 |
  | mobile      | CHAR(11)     | UNIQUE NULL INDEX         | 手机号, 最多11个字符           |
  | email       | VARCHAR(100) | UNIQUE NULL INDEX         | 邮箱号, 最多100个字符          |
  | password    | VARCHAR(100) | NOT NULL                  | 密码hash值, 最多100个字符      |
  | gender      | TINYINT(1)   | NULL                      | 性别(0:未知 1:女, 2:男)        |
  | create_time | DATETIME     | DEFAULT CURRENT_TIMESTAMP | 创建时间, 默认值自动生成       |

  ```mysql
  CREATE TABLE `tb_user_basic` (
    `id` bigint(20) unsigned NOT NULL,
    `name` varchar(10) NOT NULL,
    `mobile` varchar(11) DEFAULT '',
    `email` varchar(100) DEFAULT '',
    `password` varchar(100) NOT NULL DEFAULT '',
    `gender` tinyint(1) NOT NULL DEFAULT 0,
    `create_time` datetime DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`),
    KEY `find_by_name` (`name`),
    KEY `find_by_email` (`email`),
    KEY `find_by_mobile` (`mobile`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  
  ```

- #### tb_user_more

  | Column  | DataType     | Constraints          | Description                   |
  | ------- | ------------ | -------------------- | ----------------------------- |
  | user_id | BIGINT(20)   | UNSIGNED PRIMARY KEY | tb_user_basic.id              |
  | avatar  | VARCHAR(100) | NULL                 | 头像文件名称, 最多100个字符   |
  | qr_code | VARCHAR(100) | UNIQUE NULL          | 二维码文件名称, 最多100个字符 |

  ```mysql
   CREATE TABLE `tb_user_more` (
    `user_id` bigint(20) unsigned NOT NULL,
    `avatar` varchar(100) DEFAULT '',
    `qr_code` varchar(100) DEFAULT '',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `qr_code` (`qr_code`)
  ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
  ```

- #### tb_friend_relation

  | Column   | DataType    | Constraints          | Description                           |
  | -------- | ----------- | -------------------- | ------------------------------------- |
  | id       | BIGINT(20)  | UNSIGNED PRIMARY KEY | 主键ID, 雪花算法生成uint64数值        |
  | src_id   | BIGINT(20)  | NOT NULL INDEX       | 用户自己的id,                         |
  | dst_id   | BIGINT(20)  | NOT NULL INDEX       | 用户好友的id                          |
  | note     | VARCHAR(20) | NOT NULL             | 用户给好友设定备注名称,默认为好友昵称 |
  | isAccept | TINYINT(1)  | NOT NULL DEFAULT 0   | 好友是否接受了添加请求, 默认未接受    |
  | isBlack  | TINYINT(1)  | NOT NULL DEFAULT 0   | 是否在黑名单内                        |
  | isDelete | TINYINT(1)  | NOT NULL DEFAULT 0   | 是否删除                              |

  ```mysql
  CREATE TABLE `tb_friend_relation` (
    `id` bigint(20) unsigned NOT NULL,
    `src_id` bigint(20) NOT NULL,
    `dst_id` bigint(20) NOT NULL,
    `note` varchar(10) DEFAULT '',
    `is_accept` tinyint(1) NOT NULL DEFAULT 0,
    `is_refuse` tinyint(1) NOT NULL DEFAULT 0,
    `is_delete` tinyint(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `src_dst_id_index` (`src_id`,`dst_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8
  ```


