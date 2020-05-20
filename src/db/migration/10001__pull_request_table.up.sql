CREATE TABLE pull_request (
    repository_id INT NOT NULL,
    url varchar(256) NOT NULL,
    pr_no INT NOT NULL,
    title varchar(256) NOT NULL,
    body varchar(1000),
    head varchar(256),
    base varchar(256),
    state varchar(32) NOT NULL,
    created_year int NOT NULL,
    created_month int NOT NULL,
    created_date int NOT NULL,
    created_hour int NOT NULL,
    closed_year int,
    closed_month int,
    closed_date int,
    closed_hour int
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;