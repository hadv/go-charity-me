CREATE TABLE `users` (
  `id` CHAR(36) NOT NULL,
  `firstname` VARCHAR(100),
  `lastname` VARCHAR(100),
  `email` VARCHAR(40) NOT NULL,
  `email_verification` BOOLEAN NOT NULL DEFAULT 0,
  `password` VARCHAR(256) NOT NULL,
  `token` VARCHAR(1024),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY(`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
