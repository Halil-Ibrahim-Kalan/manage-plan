SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

CREATE TABLE `plans` (
  `id` int(255) NOT NULL,
  `uid` longtext DEFAULT NULL,
  `plan` longtext DEFAULT NULL,
  `date` longtext DEFAULT NULL,
  `start` longtext DEFAULT NULL,
  `end` longtext DEFAULT NULL,
  `status` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE `users` (
  `id` int(255) NOT NULL,
  `username` longtext DEFAULT NULL,
  `password` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

ALTER TABLE `plans`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `plans`
  MODIFY `id` int(255) NOT NULL AUTO_INCREMENT;

 TABLE `users`
  MODIFY `id` int(255) NOT NULL AUTO_INCREMENT;
COMMIT;