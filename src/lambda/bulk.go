package lambda

import (
	"database/sql"
	"strings"
	"time"
)

func generateBulkCommitStatement(dtos []commitDto) (statement string, valArgs []interface{}) {
	defer timeTrack(time.Now(), "generateBulkCommitStatement")

	statement = "INSERT INTO " + commitTable + " (repository_id, hash, author_email, message, num_files, addition_loc, deletion_loc, num_parents, total_loc, year, month, day, hour, commit_time_stamp) "
	values := make([]string, len(dtos))
	valArgs = []interface{}{}
	for i, dto := range dtos {
		values[i] = "(" + strings.Repeat("?, ", 13) + "?)"
		args := dto.getListValues()
		valArgs = append(valArgs, args...)
	}
	statement = statement + "VALUES" + strings.Join(values, ",") + " ON DUPLICATE KEY UPDATE repository_id=repository_id"
	return statement, valArgs
}

func executeBulkStatement(dtos []commitDto, conn *sql.DB) {
	defer timeTrack(time.Now(), "executeBulkStatement")

	statement, valArgs := generateBulkCommitStatement(dtos)
	insertStmt, err := conn.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(valArgs...)
	if err != nil {
		panic(err.Error())
	}
}

func executeMulipleBulks(dtos []commitDto, conn *sql.DB) {
	defer timeTrack(time.Now(), "executeMulipleBulks")

	start, end, size := 0, batchSize, len(dtos)
	for {
		if end >= size {
			executeBulkStatement(dtos[start:size], conn)
			return
		}
		executeBulkStatement(dtos[start:end], conn)
		start, end = end, end+batchSize
	}
}
