CREATE TABLE IF NOT EXISTS `test_database`.`users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `age` int(11) UNSIGNED,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `test_database`.`users` (name, age) VALUES ('bob', 15);
INSERT INTO `test_database`.`users` (name, age) VALUES ('alice', 16);