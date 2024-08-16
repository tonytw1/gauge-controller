## Gauge Controller

Listens for metrics on the local MQTT message queue.
Provides a UI to select interesting metrics and forward them to gauges which announce themselves via MQTT.

React client with a Golang backend.
Used as my sandbox project for staying current in React front ends, Golang and building non Intel container images.



### Client local dev

```
cd client
npm run dev
```

UI is on `http://localhost:5173/`


### Client build

```
cd client
npm install
npm run build
```

Produces files in `client/dist`

### Cloud Build

Test locally with:
```
gcloud components install cloud-build-local
cloud-build-local --config=cloudbuild.yaml --dryrun=false --push=false .
```
