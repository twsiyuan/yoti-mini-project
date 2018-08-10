-- Adminer 4.6.3 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `user_id` varchar(100) CHARACTER SET ascii NOT NULL,
  `full_name` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `email_address` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `phone_number` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `selfie` mediumblob NOT NULL,
  `selfie_format` varchar(3) CHARACTER SET ascii NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;



DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `comment_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(100) CHARACTER SET ascii NOT NULL,
  `comment` mediumtext COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `anonymous` tinyint(4) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`comment_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `comments_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
