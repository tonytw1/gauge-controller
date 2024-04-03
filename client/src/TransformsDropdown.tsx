import {useEffect, useState} from 'react';

export function TransformsDropdown() {

    const [transforms, setTransforms] = useState<Transform[]>([]);

    function getTransformsAsync() {
        return fetch('http://10.0.46.10:32100/transforms')
            .then((response) => response.text())
            .then((responseJson) => {
                const gauges = JSON.parse(responseJson);
                setTransforms(gauges);
            })
            .catch((error) => {
                console.error(error);
            });
    }

    useEffect(() => {
        getTransformsAsync()
    }, []);

    function TransformOption({transform}: { transform: Transform }) {
        return (
            <>
                <option value={transform.Name}>{transform.Name}</option>
            </>
        )
    };

    const listItems = transforms.map(transform => <TransformOption transform={transform}/>);

    return (
        <>
            <select name="Transform">{listItems}</select>
        </>
    )
}
