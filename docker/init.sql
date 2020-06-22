USE gitwize;
CREATE TABLE repository (
	id INT NOT NULL AUTO_INCREMENT,
	name varchar(256) NOT NULL,
	status varchar(32) NOT NULL,
	url varchar(256) NOT NULL,
	username varchar(256) NULL,
	password varchar(256) NULL,
	branches varchar(8182) NULL,
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
    closed_hour int
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE commit_data (
  repository_id int(11) NOT NULL,
  hash varchar(70) NOT NULL,
  author_email varchar(100) NOT NULL,
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

INSERT INTO repository(name,status,url,username,password,branches,ctl_created_date,ctl_created_by,ctl_modified_date,ctl_modified_by,ctl_last_metric_updated)
    VALUES ("go-git","ONGOING","https://github.com/go-git/go-git.git","tester1","L4ug7bs3myyxTR7Zmj3qKXi+SR6NqUwXHi+MksVmNIuYKzlR5IjzPls2j+ck6n2Pz1tV3PGyqYezQgeq5ED43PuV0Bs=","",now(),"tester1",now(),"tester1","1970-01-01 00:00:00");
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
