import {MetricsDropdown} from "./MetricsDropdown.tsx";
import {GaugesDropdown} from "./GaugesDropdown.tsx";
import {TransformsDropdown} from "./TransformsDropdown.tsx";
import React from "react";

export function AddRoute({updateRoutes} : {updateRoutes: (routes: Route[]) => void}) {

    function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        // Prevent the browser from reloading the page
        e.preventDefault();

        const form = e.currentTarget;
        const formData = new FormData(form);
        const formJson = Object.fromEntries(formData.entries());

        const requestOptions = {
            method: 'POST',
            //headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formJson)
        }

        fetch('http://10.0.46.10:32100/routes', requestOptions)
            .then((response) => response.text())
            .then((responseJson) => {
                const routes = JSON.parse(responseJson);
                console.log("Routes: " + routes);
                updateRoutes(routes);
            })
            .catch((error) => {
                console.error(error);
            });
    }

    return (
        <form method="post" onSubmit={handleSubmit}>
            <MetricsDropdown /> to <TransformsDropdown /> to <GaugesDropdown />
            <button type="submit">Add route</button>
        </form>
    )
}
