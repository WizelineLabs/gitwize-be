CREATE TABLE repository (
	id INT NOT NULL AUTO_INCREMENT,
	name varchar(256) NOT NULL,
	status varchar(32) NOT NULL,
	url varchar(256) NOT NULL,
	username varchar(256) NULL,
	password varchar(256) NULL,
	ctl_created_date TIMESTAMP NOT NULL,
	ctl_created_by varchar(256) NOT NULL,
	ctl_modified_date TIMESTAMP,
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
	type varchar(32) NOT NULL,
	year INT NOT NULL,
	month INT NOT NULL,
	day INT NOT NULL,
	hour INT NOT NULL,
	value BIGINT NOT NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (repository_id) REFERENCES repository(id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;
