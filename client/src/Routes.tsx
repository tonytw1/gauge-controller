import { useEffect, useState } from 'react';

export function Routes({routes, updateRoutes}) {

    function DeleteRoute({row}) {

        function handleDeleteRoute(e) {
            e.preventDefault();
            const requestOptions = {
                method: 'DELETE',
            }
            return fetch('http://10.0.46.10:32100/routes/' + row.Id, requestOptions)
                .then((response) => response.json())
                .then((responseJson) => {
                    updateRoutes(responseJson);
                })
                .catch((error) => {
                    console.error(error);
                });
        }

        return (
            <button onClick={handleDeleteRoute}>Delete</button>
        )
    }

    function Route({row}) {
        return (
            <>
                <h4>{row.FromMetric} to {row.ToGauge} - <DeleteRoute row={row} /></h4>
            </>
        )
    };

    const listItems = routes.map(row => <li key={row.Id}><Route row={row} /></li>);

    return (
        <>
        <ul>{listItems}</ul>
        </>
    )
}
