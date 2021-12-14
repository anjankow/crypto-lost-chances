# crypto-lost-chances

How much money could you earn if you invested in cryptocurrencies just in a right time?
Let's find out...
Choose a month and check which currency could you buy then and how rich would you be today.

## Components

![alt text](https://gitlab.com/anjankow/crypto-lost-chances/-/blob/main/components.png)

## Deployment

The application is deployed on the Google Cloud Platform using following infrastructure:

 - App Engine: api and lost-chances-calc
 - Functions: historical-price-getter, historical-price-saver, request-saver
 - PubSub queues: progress-updates and historical-prices
 - Cloud Tasks: historical-price-req
 - Spanner DB

## Local run

Applications api and lost-chances-calc can run locally. The local run instructions are in the README files located in the corresponding directories.
