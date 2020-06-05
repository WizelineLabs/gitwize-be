package lambda

import (
	"database/sql"
	"log"
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
	result, err := conn.Exec(statement, valArgs...)
	if err != nil {
		log.Panicln(err.Error())
	}
	rows, _ := result.RowsAffected()
	log.Println("number rows affected ", rows)
}
