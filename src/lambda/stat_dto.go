package lambda

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
)

// commitDto data object for commit entity
type commitDto struct {
	RepositoryID int
	Hash         string
	Author       string
	Message      string
	NumFiles     int
	AdditionLOC  int
	DeletionLOC  int
	NumParents   int
	LOC          int
	Year         int
	Month        int
	Day          int
	Hour         int
	TimeStamp    string
}

func (dto *commitDto) getListValues() []interface{} {
	vals := []interface{}{
		dto.RepositoryID,
		dto.Hash,
		dto.Author,
		dto.Message,
		dto.NumFiles,
		dto.AdditionLOC,
		dto.DeletionLOC,
		dto.NumParents,
		dto.LOC,
		dto.Year,
		dto.Month,
		dto.Day,
		dto.Hour,
		dto.TimeStamp,
	}
	return vals
}

// getCommitDTO return dto object of commit
func getCommitDTO(c *object.Commit) commitDto {
	dto := commitDto{}
	dto.Hash = c.Hash.String()
	dto.Author = c.Author.Email
	dto.Message = c.Message
	dto.Year = c.Author.When.UTC().Year()
	dto.Month = int(c.Author.When.UTC().Month())
	dto.Day = c.Author.When.UTC().Day()
	dto.Hour = c.Author.When.UTC().Hour()
	dto.TimeStamp = c.Author.When.UTC().String()
	// dto.LOC = getLineOfCode(c)
	dto.LOC = 0 // temporary disable getting total loc, to impove perf
	fileStats, err := c.Stats()
	if err != nil {
		log.Panicln(err)
	}
	dto.NumFiles = len(fileStats)
	for _, file := range fileStats {
		dto.AdditionLOC += file.Addition
		dto.DeletionLOC += file.Deletion
	}
	dto.NumParents = c.NumParents()
	return dto
}

// getFileStatDTO return file statistic dto
func getFileStatDTO(c *object.Commit, rID int) []fileStatDTO {
	fileStats, err := c.Stats()
	if err != nil {
		log.Panicln(err)
	}
	dtos := make([]fileStatDTO, len(fileStats))
	for i, file := range fileStats {
		dto := fileStatDTO{}
		dto.RepositoryID = rID
		dto.Hash = c.Hash.String()
		dto.FileName = file.Name
		dto.AdditionLOC = file.Addition
		dto.DeletionLOC = file.Deletion
		dto.Year = c.Author.When.UTC().Year()
		dto.Month = int(c.Author.When.UTC().Month())
		dto.Day = c.Author.When.UTC().Day()
		dto.Hour = c.Author.When.UTC().Hour()
		dto.TimeStamp = c.Author.When.UTC().String()
		dtos[i] = dto
	}
	return dtos
}

func getLineOfCode(c *object.Commit) (loc int) {
	fileIter, err := c.Files()
	if err != nil {
		panic(err.Error())
	}
	err = fileIter.ForEach(func(f *object.File) error {
		lines, _ := f.Lines()
		loc += len(lines)
		return nil
	})
	return loc
}

// fileStatDTO data object for file statistic
type fileStatDTO struct {
	RepositoryID int
	Hash         string
	FileName     string
	AdditionLOC  int
	DeletionLOC  int
	Year         int
	Month        int
	Day          int
	Hour         int
	TimeStamp    string
}
