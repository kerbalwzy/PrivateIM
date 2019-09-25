# Project MySQL DB Tables Introduction  [TOP](#0)

- #### Tables

  | Name                       | Description        |
  | -------------------------- | ------------------ |
  | [tb_user_basic](#1)        | 用户基础信息       |
  | [tb_friendship](#2)        | 用户好友关系表     |
  | [tb_group_chat](#3)        | 群聊表             |
  | [tb_user_group_chat](#4)   | 用户与群聊关系表   |
  | [tb_subscription](#5)      | 订阅号表           |
  | [tb_user_subscription](#6) | 用户与订阅号关系表 |

---

- #### <span id="1">tb_user_basic</span> [TOP](#0)

  | Column    | DataType     | Constraints           | Description                    |
  | --------- | ------------ | --------------------- | ------------------------------ |
  | id        | BIGINT(20)   | UNSIGNED PRIMARY KEY  | 主键ID, 雪花算法生成uint64数值 |
  | email     | VARCHAR(100) | NOT NULL UNIQUE INDEX | 邮箱号                         |
  | name      | VARCHAR(10)  | NOT NULL INDEX        | 昵称, 也不允许为空字符串       |
  | password  | VARCHAR(100) | NOT NULL              | 密码加盐HASH值                 |
  | mobile    | CHAR(11)     | NOT NULL DEFAULT ''   | 手机号                         |
  | gender    | TINYINT(1)   | NOT NULL DEFAULT 0    | 性别(0:未知 1:女, 2:男)        |
  | avatar    | VARCHAR(100) | NOT NULL DEFAULT ''   | 个人头像文件名称               |
  | qr_code   | VARCHAR(100) | NOT NULL UNIQUE       | 个人二维码文件名称             |
  | is_delete | TINYINT(1)   | NOT NULL DEFAULT 0    | 是否注销                       |

  ```mysql
  CREATE TABLE `tb_user_basic` (
    `id` bigint(20) unsigned NOT NULL,
    `email` varchar(100) NOT NULL,
    `name` varchar(10) NOT NULL,
    `password` varchar(100) NOT NULL,
    `mobile` char(11) NOT NULL DEFAULT '',
    `gender` tinyint(1) NOT NULL DEFAULT 0,
    `avatar` varchar(100) NOT NULL DEFAULT '',
    `qr_code` varchar(100) NOT NULL,
    `is_delete` tinyint(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email_unique` (`email`),
    UNIQUE KEY `qr_code_unique` (`qr_code`),
    KEY `name_index` (`name`),
    KEY `email_index` (`email`),
    KEY `mobile_index` (`mobile`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  
  ```

----

- #### <span id="2">tb_friendship</span> [TOP](#0)

  | Column      | DataType    | Constraints         | Description                              |
  | ----------- | ----------- | ------------------- | ---------------------------------------- |
  | self_id     | BIGINT(20)  | NOT NULL INDEX      | 用户自己的ID                             |
  | friend_id   | BIGINT(20)  | NOT NULL INDEX      | 用户好友的ID                             |
  | friend_note | VARCHAR(10) | NOT NULL DEFALUT '' | 用户给好友设定备注名称,默认为空字符串    |
  | is_accept   | TINYINT(1)  | NOT NULL DEFAULT 0  | 好友是否接受了添加请求, 默认未接受       |
  | is_black    | TINYINT(1)  | NOT NULL DEFAULT 0  | 是否将好友设置为黑名单内                 |
  | is_delete   | TINYINT(1)  | NOT NULL DEFAULT 0  | 是否删除好友关系,用户注销时,删除所有关系 |

  ```mysql
  CREATE TABLE `tb_friendship` (
    `self_id` bigint(20) NOT NULL,
    `friend_id` bigint(20) NOT NULL,
    `friend_note` varchar(10) NOT NULL DEFAULT '',
    `is_accept` tinyint(1) NOT NULL DEFAULT 0,
    `is_black` tinyint(1) NOT NULL DEFAULT 0,
    `is_delete` tinyint(1) NOT NULL DEFAULT 0,
    UNIQUE KEY `self_friend_id_unique` (`self_id`,`friend_id`),
    KEY `self_id_index` (`self_id`),
    KEY `friend_id_index` (`friend_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```

- #### <span id="3">tb_group_chat</span> [TOP](#0)

  | Column     | DataType     | Constraints         | Description                   |
  | ---------- | ------------ | ------------------- | ----------------------------- |
  | id         | BIGINT(20)   | NOT NULL INDEX      | 群聊ID,雪花算法生成uint64数值 |
  | name       | VARCHAR(10)  | NOT NULL            | 群聊名称, 也不允许为空字符串  |
  | manager_id | BIGINT(20)   | NOT NULL INDEX      | 群主用户ID                    |
  | avatar     | VARCHAR(100) | NOT NULL DEFAULT '' | 群聊头像文件名称              |
  | qr_code    | VARCHAR(100) | NOT NULL UNIQUE     | 群聊二维码文件名称            |
  | is_delete  | TINYINT(1)   | NOT NULL DEFAULT 0  | 是否解散                      |

  ```mysql
  CREATE TABLE `tb_group_chat` (
    `id` bigint(20) NOT NULL,
    `name` varchar(10) NOT NULL,
    `manager_id` bigint(20) NOT NULL,
    `avatar` varchar(100) NOT NULL DEFAULT '',
    `qr_code` varchar(100) NOT NULL,
    `is_delete` tinyint(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `qr_code_unique` (`qr_code`),
    KEY `manager_id_index` (`manager_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```

- #### <span id="4">tb_user_group_chat</span> [TOP](#0)

  | Column    | DataType    | Constraints        | Description                                |
  | --------- | ----------- | ------------------ | ------------------------------------------ |
  | group_id  | BIGINT(20)  | NOT NULL INDEX     | 群聊ID                                     |
  | user_id   | BIGINT(20)  | NOT NULL INDEX     | 用户ID                                     |
  | user_note | VARCHAR(10) | NOT NULL           | 用户在此群的昵称,默认为用户昵称            |
  | is_delete | TINYINT(1)  | NOT NULL DEFAULT 0 | 是否退出群聊, 当群聊被解散时, 所有人都退出 |

  ```mysql
  CREATE TABLE `tb_user_group_chat` (
    `group_id` bigint(20) NOT NULL,
    `user_id` bigint(20) NOT NULL,
    `user_note` varchar(10) NOT NULL,
    `is_delete` tinyint(1) NOT NULL DEFAULT 0,
    UNIQUE KEY `group_user_id_unique` (`group_id`,`user_id`),
    KEY `group_id_index` (`group_id`),
    KEY `user_id_index` (`user_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```

- #### <span id="5">tb_subscription</span> [TOP](#0)

  | Column     | DataType     | Constraints         | Description                        |
  | ---------- | ------------ | ------------------- | ---------------------------------- |
  | id         | BIGINT(20)   | NOT NULL INDEX      | 订阅号ID,雪花算法生成uint64数值    |
  | name       | VARCHAR(10)  | NOT NULL UNIQUE     | 订阅号名称, 唯一, 不允许为空字符串 |
  | manager_id | BIGINT(20)   | NOT NULL INDEX      | 订阅号拥有者ID                     |
  | intro      | VARCHAR(150) | NOT NULL            | 订阅号简介                         |
  | avatar     | VARCHAR(100) | NOT NULL DEFALUT '' | 订阅号头像文件名称                 |
  | qr_code    | VARCHAR(100) | NOT NULL UNIQUE     | 订阅号二维码文件名称               |
  | is_delete  | TINYINT(1)   | NOT NULL DEFAULT 0  | 是否注销                           |

  ```mysql
  CREATE TABLE `tb_subscription` (
    `id` bigint(20) NOT NULL,
    `name` varchar(10) NOT NULL,
    `manager_id` bigint(20) NOT NULL,
    `intro` varchar(150) NOT NULL,
    `avatar` varchar(100) NOT NULL DEFAULT '',
    `qr_code` varchar(100) NOT NULL,
    `is_delete` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `qr_code_unique` (`qr_code`),
    UNIQUE KEY `name_unique` (`name`),
    KEY `manager_id_index` (`manager_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```

- #### <span id="6">tb_user_subscription</span> [TOP](#0)

  | Column    | DataType   | Constraints        | Description                                 |
  | --------- | ---------- | ------------------ | ------------------------------------------- |
  | subs_id   | BIGINT(20) | NOT NULL INDEX     | 订阅号ID                                    |
  | user_id   | BIGINT(20) | NOT NULL INDED     | 用户ID                                      |
  | is_delete | TINYINT    | NOT NULL DEFAULT 0 | 是否取消关注, 当订阅号注销时,所有人取消关注 |

  ```mysql
  CREATE TABLE `tb_user_subscription` (
    `subs_id` bigint(20) NOT NULL,
    `user_id` bigint(20) NOT NULL,
    `is_delete` tinyint(1) NOT NULL DEFAULT 0,
    UNIQUE KEY `subs_user_id_index` (`subs_id`,`user_id`),
    KEY `subs_id_index` (`subs_id`),
    KEY `user_id_index` (`user_id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```
