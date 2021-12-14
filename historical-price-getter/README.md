# Historical Price Getter

It is responsible for getting the historical price of a certain cryptocurrency and write it to the `historical-prices` PubSub queue.
As an input the function gets following data:

    - request ID
    - cryptocurrency name to lookup
    - fiat name
    - date

At first, the function checks if this price already exists in the DB. If so, then it simply reads it from there and wirtes to the message queue.

If the price doesn't exist in the DB yet, it needs to be requested from the historical price provider.
The function requests the prices information agreggated over the given period (month and year) and gets the maximal and minimal value.
Then it writes the message to the PubSub queue.


## Deployment

This is a GCP Function triggered by a HTTP request.

Entrypoint: `GetHistoricalPrice`
