ALTER TABLE file_stat_data
ADD COLUMN (
    churn_cnt INT(11) NULL DEFAULT 0,
    refactoring_cnt INT(11) NULL DEFAULT 0
)
