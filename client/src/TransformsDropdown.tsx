import {useEffect, useState} from 'react';

export function TransformsDropdown({apiUrl}: { apiUrl: string }) {

    const [transforms, setTransforms] = useState<Transform[]>([]);

    function getTransformsAsync() {
        return fetch(apiUrl + '/transforms')
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
    }

    const listItems = transforms.map(transform => <TransformOption transform={transform}/>);

    return (
        <>
            <select name="Transform">{listItems}</select>
        </>
    )
}
