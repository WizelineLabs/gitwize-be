USE gitwize;
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
  insertion_point int(11) DEFAULT 0,
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
  modification_loc int(11) DEFAULT NULL,
  PRIMARY KEY (repository_id,hash,file_name)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE sonarqube (
  user_email varchar(256) NOT NULL,
  repository_id int NOT NULL,
  token varchar(70) NOT NULL,
  project_key varchar(70) NOT NULL,
  branch varchar(256) DEFAULT 'master',
  quality_gates varchar(128) DEFAULT 'passed',
  bugs int DEFAULT 0,
  bugs_rating varchar(10) DEFAULT "",
  vulnerabilities int DEFAULT 0,
  vulnerabilities_rating varchar(10) DEFAULT "",
  code_smells int DEFAULT 0,
  coverage float DEFAULT 0.0,
  duplications float DEFAULT 0.0,
  duplication_blocks int DEFAULT 0,
  cognitive_complexity int DEFAULT 0,
  cyclomatic_complexity int DEFAULT 0,
  security_hotspots int DEFAULT 0,
  technical_debt int DEFAULT 0,
  technical_debt_rating varchar(10) DEFAULT "",
  last_updated timestamp NOT NULL,
  PRIMARY KEY (user_email,repository_id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

-- Calculate PR Open metric for repository
DROP PROCEDURE IF EXISTS calculate_metric_open_pr;
DELIMITER $$
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
END $$
DELIMITER ;
-- Calculate for all repos
DROP PROCEDURE IF EXISTS calculate_metric_open_pr_all_repos;
DELIMITER $$
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
    FETCH NEXT FROM cur INTO repositoryId;

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
END $$
DELIMITER ;

INSERT INTO repository_user(user_email,repo_full_name,repo_id,name,access_token,branches)
    VALUES ("tester@wizeline.com","wizeline/gitwize-be",1,"gitwize-be","L4ug7bs3myyxTR7Zmj3qKXi+SR6NqUwXHi+MksVmNIuYKzlR5IjzPls2j+ck6n2Pz1tV3PGyqYezQgeq5ED43PuV0Bs=","master");
INSERT INTO repository(repo_full_name,name,status,url,access_token,branches,num_ref,ctl_created_date,ctl_created_by,ctl_modified_date,ctl_modified_by,ctl_last_metric_updated)
    VALUES ("wizeline/gitwize-be", "gitwize-be","AVAILABLE","https://github.com/wizeline/gitwize-be","L4ug7bs3myyxTR7Zmj3qKXi+SR6NqUwXHi+MksVmNIuYKzlR5IjzPls2j+ck6n2Pz1tV3PGyqYezQgeq5ED43PuV0Bs=","master",1,now(),"tester@wizeline.com",now(),"tester@wizeline.com","1970-01-01 00:00:00");

INSERT INTO repository_user(user_email,repo_full_name,repo_id,name,access_token,branches)
    VALUES ("tester@wizeline.com","wizeline/gitwize",2,"gitwize","L4ug7bs3myyxTR7Zmj3qKXi+SR6NqUwXHi+MksVmNIuYKzlR5IjzPls2j+ck6n2Pz1tV3PGyqYezQgeq5ED43PuV0Bs=","master");
INSERT INTO repository(repo_full_name,name,status,url,access_token,branches,num_ref,ctl_created_date,ctl_created_by,ctl_modified_date,ctl_modified_by,ctl_last_metric_updated)
    VALUES ("wizeline/gitwize", "gitwize","LOADING","https://github.com/wizeline/gitwize","L4ug7bs3myyxTR7Zmj3qKXi+SR6NqUwXHi+MksVmNIuYKzlR5IjzPls2j+ck6n2Pz1tV3PGyqYezQgeq5ED43PuV0Bs=","master",1,now(),"tester@wizeline.com",now(),"tester@wizeline.com","1970-01-01 00:00:00");

INSERT INTO commit_data (repository_id, hash, author_email, author_name, message, num_files, addition_loc, deletion_loc, num_parents, insertion_point, total_loc, year, month, day, hour, commit_time_stamp)
    VALUES ('1', 'testhash0001', 'test@wizeline.com', 'test', 'test message', '1', '32', '20', '1', '3', '1000', 2020, '6', '20', '0', "2020-06-20 00:00:00");
INSERT INTO commit_data (repository_id, hash, author_email, author_name, message, num_files, addition_loc, deletion_loc, num_parents, insertion_point, total_loc, year, month, day, hour, commit_time_stamp)
    VALUES ('1', 'testhash0002', 'test@wizeline.com', 'test', 'test message', '3', '100', '90', '1', '14', '1010', 2020, '6', '22', '8', "2020-06-22 08:00:00");

INSERT INTO file_stat_data (repository_id, hash, author_email, author_name, file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp, modification_loc)
    VALUES ('1', 'testhash0001', 'test@wizeline.com', 'test', 'file1', '1', '2', 2020, '6', '20', '0', "2020-06-20 00:00:00", 3);
INSERT INTO file_stat_data (repository_id, hash, author_email, author_name, file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp, modification_loc)
    VALUES ('1', 'testhash0001', 'test@wizeline.com', 'test', 'file2', '10', '5', 2020, '6', '22', '8', "2020-06-22 08:00:00", 4);
INSERT INTO file_stat_data (repository_id, hash, author_email, author_name, file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp, modification_loc)
    VALUES ('1', 'testhash0002', 'test@wizeline.com', 'test', 'unusualfile', '1001', '5', 2020, '7', '7', '8', "2020-07-07 08:00:00", 100);

INSERT INTO pull_request(repository_id, url, pr_no, title, body, head, base, state, created_by, created_year, created_month, created_day, created_hour, closed_year, closed_month, closed_day, closed_hour, additions, deletions, review_duration)
    VALUES ('1', 'https://github.com/wizeline/gitwize-be/pull/1', '1', 'GWZ-23 verifies access token','','', 'master', 'merged', 'ltvcuong', 2020, 202005, 20200507, 2020050711, 2020, 202005, 20200511, 2020051104, 113, 7, 322766);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,100);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,110);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,120);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,130);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,140);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,150);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 1,160);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,105);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,115);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,125);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,135);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,145);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,155);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 6,(2020*100 + 6)* 100 + 2,((2020*100 + 6)* 100 + 2) * 100 + 3,165);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,200);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,210);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,220);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,230);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,240);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,250);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 1,260);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,205);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,215);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,225);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,235);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,245);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,255);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 6,(2020*100 + 6)* 100 + 1,((2020*100 + 6)* 100 + 1) * 100 + 4,265);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,300);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,310);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,320);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,330);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,340);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,350);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 1,360);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,305);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,315);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,325);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,335);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,345);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,355);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 5,(2020*100 + 5)* 100 + 31,((2020*100 + 5)* 100 + 31) * 100 + 3,365);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,400);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,410);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,420);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,430);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,440);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,450);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 2,460);

INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",1,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,405);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",2,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,415);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",3,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,425);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",4,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,435);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",5,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,445);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",6,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,455);
INSERT INTO metric(repository_id,branch,type,year,month,day,hour,value) VALUES (1,"master",7,2020,2020*100 + 5,(2020*100 + 5)* 100 + 30,((2020*100 + 5)* 100 + 30) * 100 + 5,465);
