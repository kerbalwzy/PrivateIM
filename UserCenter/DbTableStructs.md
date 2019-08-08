# UserCenter

- #### Tables

  | Name               | Description  |
  | ------------------ | ------------ |
  | tb_user_basic      | 用户基础信息 |
  | tb_user_more       | 用户更多信息 |
  | tb_friend_relation | 好友关系表   |
  |                    |              |

---

- #### tb_user_basic

  | Column      | DataType     | Constraints               | Description                    |
  | ----------- | ------------ | ------------------------- | ------------------------------ |
  | id          | BIGINT(20)   | UNSIGNED PRIMARY KEY      | 主键ID, 雪花算法生成uint64数值 |
  | name        | VARCHAR(10)  | NOT NULL INDEX            | 昵称, 最多10个字符, 不允许为空 |
  | mobile      | CHAR(11)     | UNIQUE NULL INDEX         | 手机号, 最多11个字符           |
  | email       | VARCHAR(100) | UNIQUE NULL INDEX         | 邮箱号, 最多100个字符          |
  | password    | VARCHAR(100) | NOT NULL                  | 密码hash值, 最多100个字符      |
  | gender      | TINYINT(1)   | NULL                      | 性别(0:女,1:男, null:未知)     |
  | create_time | DATETIME     | DEFAULT CURRENT_TIMESTAMP | 创建时间, 默认值自动生成       |
  | update_time | DATETIME     | DEFAULT CURRENT_TIMESTAMP | 最后活跃时间, 默认值自动生成   |

  ```mysql
  CREATE TABLE `tb_user_basic` (
    `id` bigint(20) unsigned NOT NULL,
    `name` varchar(10) NOT NULL,
    `mobile` varchar(11) DEFAULT '',
    `email` varchar(100) DEFAULT '',
    `password` varchar(100) NOT NULL DEFAULT '',
    `gender` tinyint(1) NOT NULL DEFAULT -1,
    `create_time` datetime DEFAULT current_timestamp(),
    `update_time` datetime DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`),
    KEY `find_by_name` (`name`),
    KEY `find_by_email` (`email`),
    KEY `find_by_mobile` (`mobile`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8

  ```

- #### tb_user_more

  | Column  | DataType     | Constraints          | Description               |
  | ------- | ------------ | -------------------- | ------------------------- |
  | user_id | BIGINT(20)   | UNSIGNED PRIMARY KEY | tb_user_basic.id          |
  | avatar  | VARCHAR(100) | NULL                 | 头像地址, 最多100个字符   |
  | qr_code | VARCHAR(100) | UNIQUE NULL          | 二维码地址, 最多100个字符 |

  ```mysql
  CREATE TABLE `tb_user_more` (
    `usr_id` bigint(20) unsigned NOT NULL,
    `avatar` varchar(100) DEFAULT NULL,
    `qr_code` varchar(100) DEFAULT NULL,
    PRIMARY KEY (`usr_id`),
    UNIQUE KEY `qr_code` (`qr_code`)
  ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
  ```

- #### tb_friend_relation

  | Column   | DataType   | Constraints          | Description                    |
  | -------- | ---------- | -------------------- | ------------------------------ |
  | id       | BIGINT(20) | UNSIGNED PRIMARY KEY | 主键ID, 雪花算法生成uint64数值 |
  | src_id   | BIGINT(20) | NOT NULL INDEX       | 用户自己的id,                  |
  | dst_id   | BIGINT(20) | NOT NULL INDEX       | 用户好友的id                   |
  | isActive | TINYINT(1) | NOT NULL DEFAULT 1   | 是否激活, 控制黑名单           |
  | isDelete | TINYINT(1) | NOT NULL DEFAULT 0   | 是否删除                       |

  ```mysql
  CREATE TABLE `tb_friend_relation` (
    `id` bigint(20) unsigned NOT NULL,
    `src_id` bigint(20) NOT NULL,
    `dst_id` bigint(20) NOT NULL,
    `isActive` tinyint(1) NOT NULL DEFAULT 1,
    `isDelete` tinyint(1) NOT NULL DEFAULT 1,
    PRIMARY KEY (`id`),
    KEY `src_user` (`src_id`),
    KEY `dst_user` (`dst_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```


