# Historical Price Saver

It is responsible for saving the messages from the queue into the Spanner DB.
It doesn't duplicate the entries.

## Deployment

This is a GCP Function triggered by the `historical-prices` PubSub queue.

Entrypoint: `SavePrice`
