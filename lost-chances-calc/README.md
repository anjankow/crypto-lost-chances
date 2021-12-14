# Calculator of Lost Chances

This application provides the actual functionality of the system - calculating the possible investments and chosing the best one. It gets the input data from the API application and does the following:

 - starts listening on the historical-prices topic for this request ID
 - dispatches the tasks to the historical price getter
 - while waiting for the historical prices, it gets the current prices from the provider
 - collects the historical prices
 - calculates the best investment and returns the result.

## Local run

To run the application locally, google authorization credentials have to be provided.

```bash
GOOGLE_APPLICATION_CREDENTIALS=<path>/credentials.json go run main.go
```

## Infrastructure

The application uses Cloud Tasks to dispatch tasks for the historical price getter and two PubSub message queues:
 - progress-update - producer
 - historical-prices - consumer
