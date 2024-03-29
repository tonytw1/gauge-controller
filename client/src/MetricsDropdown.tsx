import {useEffect, useState} from 'react';

export function MetricsDropdown() {

    const [metrics, setMetrics] = useState<Metric[]>([]);

    function getMetricsAsync() {
        return fetch('http://10.0.46.10:32100/metrics')
            .then((response) => response.text())
            .then((responseJson) => {
                const metrics = JSON.parse(responseJson);
                setMetrics(metrics);

            })
            .catch((error) => {
                console.error(error);
            });
    }

    useEffect(() => {
        console.log("get metrics");
        getMetricsAsync()
    }, []);

    function MetricOption({row}: { row: Metric }) {
        return (
            <>
                <option value={row.Name}>{row.Name} ({row.Value})</option>
            </>
        )
    };

    const listItems = metrics.map(row => <MetricOption key={row.Name} row={row}/>);

    return (
        <>
            <select name="Metric">{listItems}</select>
        </>
    )
}
