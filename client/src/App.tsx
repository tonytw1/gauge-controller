import './App.css'
import {Routes} from './Routes.tsx';
import {AddRoute} from './AddRoute.tsx';
import {useEffect, useState} from "react";

function App() {

    const apiUrl = 'http://10.0.46.10:32100';

    const [routes, setRoutes] = useState<Route[]>(new Array<Route>());

    function updateRoutes(routes: Route[]): void {
        setRoutes(routes)
    };

    function getRoutesAsync() {
        return fetch(apiUrl + '/routes')
            .then((response) => response.text())
            .then((responseJson) => {
                const routes: Route[] = JSON.parse(responseJson);
                updateRoutes(routes);
            })
            .catch((error) => {
                console.error(error);
            });
    }

    useEffect(() => {
        getRoutesAsync()
    }, []);

    return (
        <>
            <h1>Routes</h1>
            <AddRoute updateRoutes={updateRoutes} apiUrl={apiUrl}/>
            <hr/>
            <Routes routes={routes} updateRoutes={updateRoutes} apiUrl={apiUrl}/>
            <hr/>
        </>
    )
}

export default App
