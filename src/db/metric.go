package db

func GetMetricBaseOnType(idRepository string, metricTypeVal MetricsType) (map[string][]Metric, error) {
	result := make(map[string][]Metric)
	var metricTypes map[MetricsType]string
	if metricTypeVal == ALL {
		metricTypes = MapTypeMetricToName
	} else {
		metricTypes = make(map[MetricsType]string)
		metricTypes[metricTypeVal] = MapTypeMetricToName[metricTypeVal]
	}

	for metricTypeInt, metricTypeName := range metricTypes {
		var metrics []Metric
		err := gormDB.Where("repository_id = ? AND type = ?", idRepository, metricTypeInt).Find(&metrics).Error
		if err != nil {
			return nil, err
		}
		result[metricTypeName] = metrics
	}
	return result, nil
}
