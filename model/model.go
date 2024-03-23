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
	ToGauge string
}
