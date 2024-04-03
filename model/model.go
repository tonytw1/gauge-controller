package model

type Metric struct {
	Name  string
	Value string
}

type Gauge struct {
	Name     string
	MaxValue string
}

type Route struct {
	Id         string
	FromMetric string
	Transform  string
	ToGauge    string
}

type Transform struct {
	Name      string
	Transform func(string) int
}
