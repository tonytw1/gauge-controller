import { useEffect, useState } from 'react';

export function MetricsDropdown() {

    const [metrics, setMetrics] = useState([]);

    function getMetricsAsync() {
       return fetch('http://10.0.46.10:32100/metrics')
       .then((response) => response.json())
       .then((responseJson) => {
         setMetrics(responseJson);
       })
       .catch((error) => {
         console.error(error);
       });
    }

    useEffect(() => getMetricsAsync, []);

    function MetricOption({row}) {
        return (
            <>
                <option value={row.Name}>{row.Name} ({row.Value})</option>
            </>
        )
    };

    const listItems = metrics.map(row => <MetricOption key={row.Name} row={row}/> );

    return (
        <>
        <select name="Metric">{listItems}</select>
        </>
    )
}
