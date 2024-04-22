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
