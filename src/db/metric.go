package db

import (
	"github.com/jinzhu/gorm"
	"strconv"
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
		metrics := make([]Metric, 0)
		if metricTypeInt == PRS_OPENED {
			if err := gormDB.Debug().Where("repository_id = ? AND type = ? AND hour >= ? AND hour <= ? AND (hour%100=23 OR hour=?)",
				idRepository, metricTypeInt, from, to, to).Select("repository_id, " +
				"branch, type, value, year, month, day").Find(&metrics).Error; err != nil {
				return nil, err
			}
		} else {
			if err := gormDB.Debug().Where("repository_id = ? AND type = ? AND hour >= ? AND hour <= ?",
				idRepository, metricTypeInt, from, to).Select("repository_id, " +
				"branch, type, sum(value) as value, year, month, day").Group("repository_id, branch, type, " +
				"year, month, day").Find(&metrics).Error; err != nil {
				return nil, err
			}
		}
		metricDTOs := make([]MetricDTO, 0)
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

func DeleteMetricsInOneRepo(idRepository int) error {
	if err := gormDB.Where("repository_id = ?", idRepository).Delete(Metric{}).Error; err != nil {
		return err
	}
	return nil
}

func GetQuarterlyTrends(idRepository string, epochFrom int64, epochTo int64) (map[string]int, error) {
	dateFrom := time.Unix(epochFrom, 0)
	dateTo := time.Unix(epochTo, 0)
	yearFrom, monthFrom, dayFrom, hourFrom := dateFrom.Year(), int(dateFrom.Month()), dateFrom.Day(), dateFrom.Hour()
	from := ((yearFrom*100+monthFrom)*100+dayFrom)*100 + hourFrom
	yearTo, monthTo, dayTo, hourTo := dateTo.Year(), int(dateTo.Month()), dateTo.Day(), dateTo.Hour()
	to := ((yearTo*100+monthTo)*100+dayTo)*100 + hourTo

	rejectedPRs := make([]RejectedMergedPR, 0)
	if err := gormDB.Debug().
		Select("(month%100) as month, SUM(value) as value").
		Where("repository_id = ? AND type = ? AND hour >= ? AND hour <= ? ", idRepository, PRS_REJECTED, from, to).
		Group("year, month").
		Order("year, month").
		Find(&rejectedPRs).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	mergedPRs := make([]RejectedMergedPR, 0)
	if err := gormDB.Debug().
		Select("(month%100) as month, SUM(value) as value").
		Where("repository_id = ? AND type = ? AND hour >= ? AND hour <= ? ", idRepository, PRS_MERGED, from, to).
		Group("year, month").
		Order("year, month").
		Find(&mergedPRs).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	mergedPRsMap := make(map[int]int)
	for _, pr := range mergedPRs {
		mergedPRsMap[pr.Month] = pr.Value
	}

	result := make(map[string]int)
	for _, v := range rejectedPRs {
		result[time.Month(v.Month).String()] = v.Value * 100 / (v.Value + mergedPRsMap[v.Month])
	}
	return result, nil
}
