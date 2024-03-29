type Gauge = {
    Name: string;
    MaxValue: string;
};

type Metric = {
    Name: string;
    Value: string;
};

type Route = {
    Id: string
    FromMetric: string;
    ToGauge: string;
};