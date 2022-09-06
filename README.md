# Autocomplete Proof of Concept.

Takes some data in the form of sentences and generates a trigram file of suggestions using disk based processing.  

Also spins up a reference app to see the data.


## Get it running

1. Install air 
    `go install github.com/cosmtrek/air@latest`
2. Run the service  
    `air`   

This will open two ports, 8080 for the app, and 8081 for admin type services.

## Use it

Curl some sentence data at the service for a specific site ID.

`curl localhost:8080/data/aoeu12 -d @./test-data/example-products-1.txt`

Then hit the reference app

http://localhost:8080/?siteId=aoeu12

