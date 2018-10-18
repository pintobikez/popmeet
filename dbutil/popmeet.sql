CREATE DATABASE popmeet;

USE popmeet;

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(100) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `event` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `start_datetime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end_datetime` datetime NOT NULL,
  `location` varchar(255) NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT 1,
  `fk_created_by` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`fk_created_by`) REFERENCES user(`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `event_users` (
  `fk_event` int(11) unsigned NOT NULL,
  `fk_user` int(11) unsigned NOT NULL,
  FOREIGN KEY (`fk_event`) REFERENCES event(`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
  FOREIGN KEY (`fk_user`) REFERENCES user(`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
  UNIQUE KEY unique_keys (fk_event,fk_user)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `login_provider` (
  `id` int(2) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `web_clientid` varchar(255) NOT NULL,
  `web_secret` varchar(255) NOT NULL,
  `android_clientid` varchar(255) NOT NULL,
  `android_secret` varchar(255) NOT NULL,
  `iphone_clientid` varchar(255) NOT NULL,
  `iphone_secret` varchar(255) NOT NULL,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `user_security` (
  `id` int(2) unsigned NOT NULL AUTO_INCREMENT,
  `fk_user` int(11) unsigned NOT NULL,
  `fk_login_provider` int(11) unsigned NULL,
  `hash` varchar(255) NULL,
  `last_machine` varchar(255) NOT NULL,
  `last_login_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`fk_user`) REFERENCES user(`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
  FOREIGN KEY (`fk_login_provider`) REFERENCES login_provider(`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `language` (
  `id` int(2) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `name_iso2` varchar(2) NOT NULL,
  `name_iso3` varchar(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `interest` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `users_profile_interests` (
  `fk_interest` int(11) unsigned NOT NULL,
  `fk_user_profile` int(11) unsigned NOT NULL,
  FOREIGN KEY (`fk_interest`) REFERENCES interest(`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
  FOREIGN KEY (`fk_user_profile`) REFERENCES user_profile(`id`) ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	CREATE TABLE IF NOT EXISTS `user_profile` (
	  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	  `fk_user` int(11) unsigned NOT NULL,
	  `fk_language` int(11) unsigned NOT NULL DEFAULT 1,
    `age_range` enum('18-25','26-32','33-39','40-46','47-53','54-60','61-70','+70') NOT NULL,
	  `sex` enum('male','female') NOT NULL,
	  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
	  PRIMARY KEY (`id`),
	  FOREIGN KEY (`fk_language`) REFERENCES language(`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
	  FOREIGN KEY (`fk_user`) REFERENCES user(`id`) ON UPDATE CASCADE ON DELETE RESTRICT
	) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;


INSERT INTO language VALUES(null, 'English', 'EN', 'ENG');
INSERT INTO login_provider VALUES(null, 'Api', 'CLIENTID-WEB', 'SECRET-WEB', 'CLIENTID-ANDROID', 'SECRET-ANDROID', 'CLIENTID-IPHONE', 'SECRET-IPHONE', NOW());
INSERT INTO login_provider VALUES(null, 'Google', 'CLIENTID-WEB', 'SECRET-WEB', 'CLIENTID-ANDROID', 'SECRET-ANDROID', 'CLIENTID-IPHONE', 'SECRET-IPHONE', NOW());
INSERT INTO interest VALUES(null, 'internet');
INSERT INTO interest VALUES(null, 'cars');
INSERT INTO interest VALUES(null, 'rugby');
INSERT INTO interest VALUES(null, 'football');