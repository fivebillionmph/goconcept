DROP TABLE IF EXISTS `base_users`;
drop table if exists `base_concept_relationships`;
drop table if exists `base_concepts_history`;
drop table if exists `base_concept_data`;
drop table if exists `base_concepts`;


CREATE TABLE `base_users` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`timestamp` INT(11) NOT NULL,
	`email` VARCHAR(128) NOT NULL,
	`password` VARCHAR(64) NOT NULL,
	`username` VARCHAR(128) NOT NULL,
	`level` TINYINT(1) NOT NULL,
	`active` TINYINT(1) NOT NULL,
	PRIMARY KEY(`id`),
	UNIQUE KEY(`email`),
	UNIQUE KEY(`username`)
) ENGINE=InnoDB;

create table `base_api_keys` (
	`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` int(11) UNSIGNED NOT NULL,
	`timestamp` int(11) NOT NULL,
	`active` tinyint(1) NOT NULL,
	`key` varchar(32) NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY(`key`),
	FOREIGN KEY(`user_id`) REFERENCES base_users(`id`)
) ENGINE=InnoDB;

create table `base_concepts` (
	`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`timestamp` int(11) NOT NULL,
	`active` tinyint(1) NOT NULL,
	`type` varchar(64) NOT NULL,
	`name` varchar(128) NOT NULL,
	PRIMARY KEY(`id`),
	UNIQUE KEY(`type`, `name`)
) ENGINE=InnoDB;

create table `base_concept_data` (
	`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`timestamp` int(11) NOT NULL,
	`concept_id` int(11) UNSIGNED NOT NULL,
	`active` tinyint(1) NOT NULL,
	`key` varchar(64) NOT NULL,
	`value` LONGTEXT,
	PRIMARY KEY(`id`),
	FOREIGN KEY(`concept_id`) REFERENCES concepts(`id`)
) ENGINE=InnoDB;

create table `base_concepts_history` (
	`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`timestamp` int(11) NOT NULL,
	`history_type` varchar(32) NOT NULL,
	`type` varchar(64) NOT NULL,
	`name` varchar(128) NOT NULL,
	`data` LONGTEXT,
	PRIMARY KEY(`id`)
) ENGINE=InnoDB;

create table `base_concept_relationships` (
	`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	`timestamp` int(11) NOT NULL,
	`id1` int(11) UNSIGNED NOT NULL,
	`id2` int(11) UNSIGNED NOT NULL,
	`string1` varchar(64),
	`string2` varchar(64),
	PRIMARY KEY(`id`),
	FOREIGN KEY(`id1`) REFERENCES concepts(`id`),
	FOREIGN KEY(`id2`) REFERENCES concepts(`id`),
	UNIQUE KEY(`id1`, `id2`, `string1`, `string2`)
) ENGINE=InnoDB;
