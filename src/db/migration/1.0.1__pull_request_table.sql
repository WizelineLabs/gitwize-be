CREATE TABLE pull_request (
    repository_id INT NOT NULL,
    pr_no INT NOT NULL,
    title varchar(256) NOT NULL,
    body varchar(1000),
    head varchar(256),
    base varchar(256),
    state varchar(32) NOT NULL,
    created_by varchar(32) NOT NULL,
    created_at TIMESTAMP,
    merged_at TIMESTAMP
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;