DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name` char(80) COLLATE utf8mb4_bin NOT NULL,
    `login` char(254) COLLATE utf8mb4_bin NOT NULL,
    `password` char(80) COLLATE utf8mb4_bin NOT NULL,
    `type` enum('user','expert','moderator') COLLATE utf8mb4_bin NOT NULL,
    `contacts` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `login` (`login`),
    KEY `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Accounts';

INSERT INTO `account` VALUES (1,'Владимир Владимирович','user1','demo','user','vv@yandex.ru');
INSERT INTO `account` VALUES (2,'Сергей Николаевич','user2','demo','user','sn@yandex.ru');
INSERT INTO `account` VALUES (3,'Елена Михайловна','user3','demo','user','em@yandex.ru');
INSERT INTO `account` VALUES (4,'Михаил Иванович','user4','demo','user','mi@yandex.ru');
INSERT INTO `account` VALUES (5,'Ксения Владимировна','user5','demo','user','kv@yandex.ru');
INSERT INTO `account` VALUES (6,'Виктор Владимирович','expert1','demo','expert','viv@yandex.ru');
INSERT INTO `account` VALUES (7,'Екатерина Вячеславовна','expert2','demo','expert','ev@yandex.ru');
INSERT INTO `account` VALUES (8,'Марат Анатольевич','moderator1','demo','moderator','ma@yandex.ru');

DROP TABLE IF EXISTS `obr`;
CREATE TABLE `obr` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `title` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `content` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `file` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `account_id` int(10) unsigned NOT NULL,
    `public` tinyint(1) unsigned NOT NULL DEFAULT 1,
    `state` enum('draft','sign','post') COLLATE utf8mb4_bin NOT NULL,
    `address` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `dtreg` datetime NOT NULL DEFAULT utc_timestamp(),
    `dtlast` datetime NOT NULL DEFAULT utc_timestamp(),
    PRIMARY KEY (`id`),
    KEY `account_id` (`account_id`),
    KEY `public` (`public`),
    KEY `state` (`state`),
    KEY `dtreg` (`dtreg`),
    KEY `dtlast` (`dtlast`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Obrashenie';

DROP TABLE IF EXISTS `obr_account`;
CREATE TABLE `obr_account` (
    `obr_id` int(10) unsigned NOT NULL,
    `account_id` int(10) unsigned NOT NULL,
    PRIMARY KEY (`obr_id`, `account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Obr Users';

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `obr_id` int(10) unsigned NOT NULL,
    `account_id` int(10) unsigned NOT NULL,
    `content` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `type` enum('draft','sign','post') COLLATE utf8mb4_bin NOT NULL,
    `dt` datetime NOT NULL DEFAULT utc_timestamp(),
    PRIMARY KEY (`id`),
    KEY `obr_id` (`obr_id`),
    KEY `account_id` (`account_id`),
    KEY `type` (`type`),
    KEY `dt` (`dt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Comments';

DROP TABLE IF EXISTS `obr_sign`;
CREATE TABLE `obr_sign` (
    `obr_id` int(10) unsigned NOT NULL,
    `account_id` int(10) unsigned NOT NULL,
    `dt` datetime NOT NULL DEFAULT utc_timestamp(),
    PRIMARY KEY (`obr_id`, `account_id`),
    KEY `dt` (`dt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Obr Signers';

DROP TABLE IF EXISTS `answer`;
CREATE TABLE `answer` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `obr_id` int(10) unsigned NOT NULL,
    `file` text COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `dt` datetime NOT NULL DEFAULT utc_timestamp(),
    PRIMARY KEY (`id`),
    KEY `obr_id` (`obr_id`),
    KEY `file` (`file`),
    KEY `dt` (`dt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='answer';
