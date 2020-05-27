
DROP table metric_new;
CREATE TABLE metric_new AS
SELECT repository_id, 'master' as branch, 5 as type, COUNT(*) as value, created_year as year,
	created_month as month, created_day as day, created_hour as hour
FROM gitwize.pull_request
WHERE state = 'open'
GROUP BY repository_id, created_year, created_month, created_day, created_hour
;

-- pr merged
INSERT INTO metric_new
SELECT repository_id, 'master' as branch, 6 as type, COUNT(*) as value, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour
FROM gitwize.pull_request
WHERE state = 'merged'
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

-- pr rejected
INSERT INTO metric_new
SELECT repository_id, 'master' as branch, 7 as type, COUNT(*) as value, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour
FROM gitwize.pull_request
WHERE state = 'rejected'
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

-- commit
INSERT INTO metric_new
SELECT repository_id, 'master', 4, COUNT(*), `year`, `month`, `day`, `hour`
FROM commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

-- line added
INSERT INTO metric_new
SELECT repository_id, 'master', 3, SUM(addition_loc), `year`, `month`, `day`, `hour`
FROM commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

-- line removed
INSERT INTO metric_new
SELECT repository_id, 'master', 2, SUM(deletion_loc), `year`, `month`, `day`, `hour`
FROM commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

-- loc
INSERT INTO metric_new
SELECT repository_id, 'master', 1, ROUND(AVG(total_loc)), `year`, `month`, `day`, `hour`
FROM commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

-- fetch metric from 2020-04 - 2020-05
SELECT repository_id, branch, type, year, month, `day`, SUM(value) as value
FROM metric_new
WHERE repository_id = 1 AND year = 2020 AND month = 5
GROUP BY repository_id, branch, type, year, `month`, `day`
;
