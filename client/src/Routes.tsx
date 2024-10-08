import { Button } from '@headlessui/react'

export function Routes({routes, updateRoutes, apiUrl}: { routes: Route[], updateRoutes: any, apiUrl: string }) {

    function DeleteRoute({route}: { route: Route }) {

        function handleDeleteRoute(e: React.FormEvent<HTMLButtonElement>) {
            e.preventDefault();
            const requestOptions = {
                method: 'DELETE',
            }
            return fetch(apiUrl + '/routes/' + route.Id, requestOptions)
                .then((response) => response.text())
                .then((responseJson) => {
                    routes = JSON.parse(responseJson);
                    updateRoutes(routes);
                })
                .catch((error) => {
                    console.error(error);
                });
        }

        return (
            <Button onClick={handleDeleteRoute}>Delete</Button>
        )
    }

    function RouteRow({route}: { route: Route }) {
        return (
            <>
                <h4>{route.FromMetric} to {route.Transform} to {route.ToGauge} - <DeleteRoute route={route}/></h4>
            </>
        )
    };

    const listItems = routes.map(route =>
        <li key={route.Id}><RouteRow route={route}/></li>
    );

    return (
        <>
            {routes.length > 0
            ? <>
                <ul>{listItems}</ul>
                </>
            : <>
                    <p>No routes</p>
                </>
            }
        </>
    );
}
