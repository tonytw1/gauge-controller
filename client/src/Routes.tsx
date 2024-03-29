export function Routes({routes, updateRoutes} : {routes: Route[], updateRoutes:any}) {

    function DeleteRoute({route} : {route:Route}) {

        function handleDeleteRoute(e: React.FormEvent<HTMLButtonElement>) {
            e.preventDefault();
            const requestOptions = {
                method: 'DELETE',
            }
            return fetch('http://10.0.46.10:32100/routes/' + route.Id, requestOptions)
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
            <button onClick={handleDeleteRoute}>Delete</button>
        )
    }

    function RouteRow({route} : {route:Route}) {
        return (
            <>
                <h4>{route.FromMetric} to {route.ToGauge} - <DeleteRoute route={route} /></h4>
            </>
        )
    };

    console.log("MEH: " + routes);
    const listItems = routes.map(route =>
        <li key={route.Id}><RouteRow route={route} /></li>
    );

    return (
        <>
        <ul>{listItems}</ul>
        </>
    )
}
