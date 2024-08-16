export function MetricsDropdown({metrics}: { metrics: Metric[] }) {

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
