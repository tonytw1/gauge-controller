import {useEffect, useState} from 'react';

export function GaugesDropdown() {

    const [gauges, setGauges] = useState<Gauge[]>([]);

    function getGaugesAsync() {
        return fetch('http://10.0.46.10:32100/gauges')
            .then((response) => response.text())
            .then((responseJson) => {
                const gauges = JSON.parse(responseJson);
                setGauges(gauges);
            })
            .catch((error) => {
                console.error(error);
            });
    }

    useEffect(() => {
        getGaugesAsync()
    }, []);

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
