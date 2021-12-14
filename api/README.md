# API

The entrypoint to the application. Gets the requests from user passes them to another components.

## Calculate request

To serve the /calculate request, the API app passes the investment information to the lost-chances-calc and waits for the response on one thread, while listening for the progress updates on another. The progress updates come from a PubSub message queue. 

The progress information and results are updated on the page via websockets.

## Local run

To run the application locally, google authorization credentials have to be provided.

```bash
ENV=prod GOOGLE_APPLICATION_CREDENTIALS=<path>/credentials.json go run main.go
```

Specifing ENV as production makes the application communicate with the lost-chances-calc running on the server.
If lost-chances-calc is running locally, the ENV variable should remain empty.

## Deployment

The application is to be deployed on the App Engine using the api.yaml file.
