SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=4;
-- commit
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 4, COUNT(*), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM gitwize.commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=3;
-- line added
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 3, SUM(addition_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM gitwize.commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=2;
-- line removed
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 2, SUM(deletion_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM gitwize.commit_data
GROUP BY repository_id, `year`, `month`, `day`, `hour`
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=5;
-- pr created
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 5 as type, created_year as year,
	created_month as month, created_day as day, created_hour as hour, COUNT(*) as value
FROM gitwize.pull_request
WHERE state = 'open'
GROUP BY repository_id, created_year, created_month, created_day, created_hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=6;
-- pr merged
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 6 as type, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour, COUNT(*) as value
FROM gitwize.pull_request
WHERE state = 'merged'
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=7;
-- pr rejected
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 7 as type, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour, COUNT(*) as value
FROM gitwize.pull_request
WHERE state = 'rejected'
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;
