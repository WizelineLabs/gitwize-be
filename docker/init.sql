USE gitwize;
CREATE TABLE repository (
	id INT NOT NULL AUTO_INCREMENT,
	name varchar(256) NOT NULL,
	status varchar(32) NOT NULL,
	url varchar(256) NOT NULL,
	username varchar(256) NULL,
	password varchar(256) NULL,
	ctl_created_date TIMESTAMP NOT NULL,
	ctl_created_by varchar(256) NOT NULL,
	ctl_modified_date TIMESTAMP NULL,
	ctl_modified_by varchar(256),
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
	as_of_date TIMESTAMP NOT NULL,
	value BIGINT NOT NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (repository_id) REFERENCES repository(id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

INSERT INTO repository(name,status,url,username,ctl_created_date,ctl_created_by)
    VALUES ("gitwize","ONGOING","https://github.com/gitwize","tester1",now(),"tester1");
INSERT INTO repository(name,status,url,username,ctl_created_date,ctl_created_by)
    VALUES ("gitwize2","ONGOING","https://github.com/gitwize2","tester2",now(),"tester2");
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",1,now(),100);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",2,now(),110);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",3,now(),120);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",4,now(),130);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",5,now(),140);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",6,now(),150);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",7,now(),160);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",1,now(),105);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",2,now(),115);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",3,now(),125);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",4,now(),135);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",5,now(),145);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",6,now(),155);
INSERT INTO metric(repository_id,branch,type,as_of_date,value) VALUES (1,"master",7,now(),165);
