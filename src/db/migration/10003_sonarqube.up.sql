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
