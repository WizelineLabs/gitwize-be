package db

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func GetMetricBaseOnType(idRepository string, metricTypeVal MetricsType, epochFrom int64, epochTo int64) (map[string][]MetricDTO, error) {
	result := make(map[string][]MetricDTO)
	var metricTypes map[MetricsType]string
	if metricTypeVal == ALL {
		metricTypes = MapTypeMetricToName
	} else {
		metricTypes = make(map[MetricsType]string)
		metricTypes[metricTypeVal] = MapTypeMetricToName[metricTypeVal]
	}

	dateFrom := time.Unix(epochFrom, 0)
	dateTo := time.Unix(epochTo, 0)
	yearFrom, monthFrom, dayFrom, hourFrom := dateFrom.Year(), int(dateFrom.Month()), dateFrom.Day(), dateFrom.Hour()
	from := ((yearFrom*100+monthFrom)*100+dayFrom)*100 + hourFrom
	yearTo, monthTo, dayTo, hourTo := dateTo.Year(), int(dateTo.Month()), dateTo.Day(), dateTo.Hour()
	to := ((yearTo*100+monthTo)*100+dayTo)*100 + hourTo

	for metricTypeInt, metricTypeName := range metricTypes {
		var metrics []Metric
		err := gormDB.Where("repository_id = ? AND type = ? AND hour >= ? AND hour <= ?",
			idRepository, metricTypeInt, from, to).Select("repository_id, " +
			"branch, type, sum(value) as value, year, month, day").Group("repository_id, branch, type, " +
			"year, month, day").Find(&metrics).Error

		if err != nil {
			return nil, err
		}
		var metricDTOs []MetricDTO
		for _, metric := range metrics {
			metricDTOs = append(metricDTOs, MetricDTO{
				BranchName: metric.BranchName,
				Type:       metric.Type,
				Value:      metric.Value,
				AsOfDate:   strconv.Itoa(metric.Month%100) + "/" + strconv.Itoa(metric.Day%100) + "/" + strconv.Itoa(metric.Year),
			})
		}
		result[metricTypeName] = metricDTOs
	}
	return result, nil
}

func UpdateMetricTable() {
	file, err := ioutil.ReadFile("src/db/update_metric_table.sql")

	if err != nil {
		log.Fatal("Failed to read sql script: " + err.Error())
	}

	requests := strings.Split(string(file), ";\n")

	log.Println("=====UpdateMetricTable=====")
	for _, request := range requests {
		log.Println(request)
		err := gormDB.Exec(request).Error
		if err != nil {
			log.Fatal("Failed to update database: " + err.Error())
		}
	}
}
