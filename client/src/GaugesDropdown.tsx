import { useEffect, useState } from 'react';

export function GaugesDropdown() {

    const [gauges, setGauges] = useState([]);

    function getGaugesAsync() {
       return fetch('http://10.0.46.10:32100/gauges')
       .then((response) => response.json())
       .then((responseJson) => {
        setGauges(responseJson);
       })
       .catch((error) => {
         console.error(error);
       });
    }

    useEffect(() => getGaugesAsync, []);

    function GaugeOption({row}) {
        return (
            <>
                <option value={row.Name}>{row.Name}</option>
            </>
        )
    };

    const listItems = gauges.map(row => <GaugeOption row={row}/> );

    return (
        <>
        <select name="Gauge">{listItems}</select>
        </>
    )
}
