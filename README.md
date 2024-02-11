# 🐉 gRPC Gateway Implementation 🐉

### Need to build a Backend API in Go? 😱 Using both gRPC and HTTP? 😌 
### Following a custom Clear Architecture pattern, shrinking the amount of code by a staggering 33%? 🤯

...

#### > No. 

## About

The best thing about using `.proto` files is that you clearly define you service's specs, requests and responses. Combining this with custom annotations on the protofiles and using gRPC-Gateway we can:

* Automatically generate an HTTP handler for each gRPC method.
* Automatically handle requests' input values (e.g. An HTTP request's body) and map them to our data structures.
* Assert validation rules for each request.
* Automatically generate a swagger spec.

And so you can basically remove the Transport Layer altogether from your API's architecture, which in turn calls the Service Layer directly.

## Commands

`make run`: Hmm... What could this possible be?

`make gen`: Based on the .proto files, generate the .pb files and the swagger documentation.

`make protoc-gen`: Based on the .proto files, generate the .pb files.

`make swagger-gen`: Based on the .proto files, generate the swagger documentation.
