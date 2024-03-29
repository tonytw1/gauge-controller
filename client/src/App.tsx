import './App.css'
import {Routes} from './Routes.tsx';
import {AddRoute} from './AddRoute.tsx';
import {useEffect, useState} from "react";

function App() {

    const [routes, setRoutes] = useState<Route[]>(new Array<Route>());

    function updateRoutes(routes: Route[]) {
        setRoutes(routes)
    };

    function getRoutesAsync() {
        return fetch('http://10.0.46.10:32100/routes')
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
            <AddRoute updateRoutes={updateRoutes}/>
            <hr/>
            <Routes routes={routes} updateRoutes={updateRoutes}/>
            <hr/>
        </>
    )
}

export default App
