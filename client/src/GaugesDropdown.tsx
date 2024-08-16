export function GaugesDropdown({gauges}: { gauges: Gauge[] }) {

    function GaugeOption({gauge}: { gauge: Gauge }) {
        return (
            <>
                <option value={gauge.Name}>{gauge.Name}</option>
            </>
        )
    };

    const listItems = gauges.map(gauge => <GaugeOption gauge={gauge}/>);

    return (
        <>
            <select name="Gauge">{listItems}</select>
        </>
    )
}
