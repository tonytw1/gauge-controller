import {MetricsDropdown} from "./MetricsDropdown.tsx";
import {GaugesDropdown} from "./GaugesDropdown.tsx";
import {TransformsDropdown} from "./TransformsDropdown.tsx";
import React, {useEffect, useState} from "react";
import { Button } from '@headlessui/react'

export function AddRoute({updateRoutes, apiUrl} : {updateRoutes: (routes: Route[]) => void, apiUrl: string}) {

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
        getTransformsAsync();
        getGaugesAsync();
        getMetricsAsync();
    }, []);

    const [gauges, setGauges] = useState<Gauge[]>([]);

    function getGaugesAsync() {
        return fetch(apiUrl + '/gauges')
            .then((response) => response.text())
            .then((responseJson) => {
                const gauges = JSON.parse(responseJson);
                setGauges(gauges);
            })
            .catch((error) => {
                console.error(error);
            });
    }

    const [metrics, setMetrics] = useState<Metric[]>([]);

    function getMetricsAsync() {
        return fetch(apiUrl + '/metrics')
            .then((response) => response.text())
            .then((responseJson) => {
                const metrics = JSON.parse(responseJson);
                setMetrics(metrics);

            })
            .catch((error) => {
                console.error(error);
            });
    }

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

        fetch(apiUrl + '/routes', requestOptions)
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
            {gauges.length > 0 && metrics.length > 0 && transforms.length > 0
                ? <>
                    <MetricsDropdown metrics={metrics}/> to <TransformsDropdown
                    transforms={transforms}/> to <GaugesDropdown gauges={gauges}/>
                    <Button type={"submit"}>Add route</Button>
                </>
                : <>
                    <p>Missing metrics, transforms or gauges; cannot add a new route</p>
                </>
            }
        </form>
    )
}
