package plugin

import "time"

type MetricModifier interface {
	UpdateMetric(mt MetricSetter)
}

func MetricTag(key string, value string) MetricModifier {
	return &metricTags{
		tags: map[string]string{key: value},
	}
}

func MetricTags(tags map[string]string) MetricModifier {
	return &metricTags{
		tags: tags,
	}
}

func MetricTimestamp(timestamp time.Time) MetricModifier {
	return &metricTimestamp{
		timestamp: timestamp,
	}
}

func MetricDescription(description string) MetricModifier {
	return &metricDescription{
		description: description,
	}
}

func MetricUnit(unit string) MetricModifier {
	return &metricUnit{
		unit: unit,
	}
}

///////////////////////////////////////////////////////////////////////////////

type metricTags struct {
	tags map[string]string
}

func (m metricTags) UpdateMetric(mt MetricSetter) {
	mt.AddTags(m.tags)
}

type metricTimestamp struct {
	timestamp time.Time
}

func (m metricTimestamp) UpdateMetric(mt MetricSetter) {
	mt.SetTimestamp(m.timestamp)
}

type metricDescription struct {
	description string
}

func (m metricDescription) UpdateMetric(mt MetricSetter) {
	mt.SetDescription(m.description)
}

type metricUnit struct {
	unit string
}

func (m metricUnit) UpdateMetric(mt MetricSetter) {
	mt.SetUnit(m.unit)
}
