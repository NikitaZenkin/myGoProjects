SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions`
(
    `ID`       varchar(200) NOT NULL,
    `userID`   varchar(200) NOT NULL,
    `userName` varchar(200) NOT NULL,
    PRIMARY KEY (`ID`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `ID`       varchar(200) NOT NULL,
    `userName` varchar(200) NOT NULL,
    `password` varchar(200) NOT NULL,
    PRIMARY KEY (`ID`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

