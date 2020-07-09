CREATE TABLE repository_user (
	user_email varchar(256) NOT NULL,
    repo_full_name varchar(256) NOT NULL,
    repo_id INT NOT NULL,
	name varchar(256) NOT NULL,
	access_token varchar(256) NULL,
	branches varchar(8182) NULL,
	PRIMARY KEY (user_email,repo_full_name)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE repository (
	id INT NOT NULL AUTO_INCREMENT,
    repo_full_name varchar(256) NOT NULL,
	name varchar(256) NOT NULL,
	status varchar(32) NOT NULL,
	url varchar(256) NOT NULL,
	access_token varchar(256) NULL,
	branches varchar(8182) NULL,
	num_ref int NOT NULL,
	ctl_created_date TIMESTAMP NOT NULL,
	ctl_created_by varchar(256) NOT NULL,
	ctl_modified_date TIMESTAMP NULL,
	ctl_modified_by varchar(256),
	ctl_last_metric_updated TIMESTAMP NULL,
	PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE INDEX idx_repo_full_name ON repository (repo_full_name);

CREATE TABLE metric (
	id BIGINT NOT NULL AUTO_INCREMENT,
	repository_id INT NOT NULL,
	branch varchar(256),
	type INT NOT NULL,
	year INT NOT NULL,
	month INT NOT NULL,
	day INT NOT NULL,
	hour INT NOT NULL,
	value BIGINT NOT NULL,
	PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE pull_request (
    repository_id INT NOT NULL,
    url varchar(256) NOT NULL,
    pr_no INT NOT NULL,
    title varchar(256) NOT NULL,
    body varchar(1000),
    head varchar(256),
    base varchar(256),
    state varchar(32) NOT NULL,
    created_by varchar(256) NOT NULL,
    created_year int NOT NULL,
    created_month int NOT NULL,
    created_day int NOT NULL,
    created_hour int NOT NULL,
    closed_year int,
    closed_month int,
    closed_day int,
    closed_hour int,
    additions int,
    deletions int,
    review_duration bigint,
    UNIQUE KEY pr_repo_no_idx (repository_id, pr_no)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE commit_data (
  repository_id int(11) NOT NULL,
  hash varchar(70) NOT NULL,
  author_email varchar(100) NOT NULL,
  author_name varchar(100) NOT NULL,
  message varchar(250) DEFAULT NULL,
  num_files int(11) DEFAULT NULL,
  addition_loc int(11) DEFAULT NULL,
  deletion_loc int(11) DEFAULT NULL,
  num_parents int(11) DEFAULT NULL,
  total_loc int(11) NOT NULL,
  year year(4) NOT NULL,
  month tinyint(1) NOT NULL,
  day tinyint(1) NOT NULL,
  hour tinyint(1) NOT NULL,
  commit_time_stamp timestamp NULL,
  PRIMARY KEY (repository_id,hash)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;


CREATE TABLE file_stat_data (
  repository_id int(11) NOT NULL,
  hash varchar(70) NOT NULL,
  author_email varchar(100) NOT NULL,
  author_name varchar(100) NOT NULL,
  file_name varchar(200) NOT NULL,
  addition_loc int(11) DEFAULT NULL,
  deletion_loc int(11) DEFAULT NULL,
  year year(4) NOT NULL,
  month tinyint(1) NOT NULL,
  day tinyint(1) NOT NULL,
  hour tinyint(1) NOT NULL,
  commit_time_stamp timestamp NULL,
  PRIMARY KEY (repository_id,hash,file_name)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

-- Calculate PR Open metric for repository
DROP PROCEDURE IF EXISTS calculate_metric_open_pr;

CREATE PROCEDURE calculate_metric_open_pr(
   IN repositoryId INT
)
    main: BEGIN
        DELETE FROM metric WHERE repository_id = repositoryId AND type = 8;

        SELECT @minPrOpen := MIN(created_hour) FROM pull_request WHERE repository_id = repositoryId;

        IF @minPrOpen IS NULL THEN
            LEAVE main;
        END IF;

        SET @hour := STR_TO_DATE(@minPrOpen, '%Y%m%d%H');
        SET @end := NOW();

        metric_loop: LOOP
        IF @hour < @end THEN
            INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
            SELECT repository_id, 'master', 8, COUNT(*), DATE_FORMAT(@hour, '%Y'), DATE_FORMAT(@hour, '%Y%m'), DATE_FORMAT(@hour, '%Y%m%d'), DATE_FORMAT(@hour, '%Y%m%d%H')
            FROM pull_request
            WHERE repository_id = repositoryId AND created_hour <= DATE_FORMAT(@hour, '%Y%m%d%H')
                AND (closed_hour = 0 OR closed_hour IS NULL OR closed_hour > DATE_FORMAT(@hour, '%Y%m%d%H'));

            SET @hour := DATE_ADD(@hour, INTERVAL 1 HOUR);
        ELSE
             LEAVE metric_loop;
        END IF;
    END LOOP metric_loop;
END;
-- Calculate for all repos
DROP PROCEDURE IF EXISTS calculate_metric_open_pr_all_repos;
CREATE PROCEDURE calculate_metric_open_pr_all_repos()
BEGIN
    DECLARE repositoryId INT DEFAULT NULL;
    DECLARE done TINYINT DEFAULT FALSE;

    # cursor over repository id
    DECLARE cur
    CURSOR FOR
    SELECT id
    FROM repository;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION

    OPEN cur;

    main_loop: LOOP
        FETCH NEXT FROM cur INTO repositoryId;

        IF done THEN
            LEAVE main_loop;
        ELSE
            CALL calculate_metric_open_pr(repositoryId);
        END IF;
    END LOOP;

    CLOSE cur;
END;