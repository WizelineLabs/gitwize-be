CREATE TABLE `commit_data` (
  `repository_id` int(11) NOT NULL,
  `hash` varchar(70) NOT NULL,
  `author_email` varchar(100) NOT NULL,
  `message` varchar(250) DEFAULT NULL,
  `num_files` int(11) DEFAULT NULL,
  `addition_loc` int(11) DEFAULT NULL,
  `deletion_loc` int(11) DEFAULT NULL,
  `num_parents` int(11) DEFAULT NULL,
  `total_loc` int(11) NOT NULL,
  `year` year(4) NOT NULL,
  `month` tinyint(1) NOT NULL,
  `day` tinyint(1) NOT NULL,
  `hour` tinyint(1) NOT NULL,
  `commit_time_stamp` timestamp NOT NULL,
  PRIMARY KEY (`repository_id`,`hash`),
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;


CREATE TABLE `file_stat_data` (
  `repository_id` int(11) NOT NULL,
  `hash` varchar(70) NOT NULL,
  `author_email` varchar(100) NOT NULL,
  `file_name` varchar(200) NOT NULL,
  `addition_loc` int(11) DEFAULT NULL,
  `deletion_loc` int(11) DEFAULT NULL,
  `year` year(4) NOT NULL,
  `month` tinyint(1) NOT NULL,
  `day` tinyint(1) NOT NULL,
  `hour` tinyint(1) NOT NULL,
  `commit_time_stamp` timestamp NOT NULL,
  PRIMARY KEY (`repository_id`,`hash`),
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;
