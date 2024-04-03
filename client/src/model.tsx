type Gauge = {
    Name: string
    MaxValue: string
};

type Metric = {
    Name: string;
    Value: string;
};

type Route = {
    Id: string
    FromMetric: string
    Transform: string
    ToGauge: string
};

type Transform = {
    Name: string
};
