# 🐉 gRPC Gateway Implementation 🐉

### Need to build a Backend API in Go? 😱 Using both gRPC and HTTP? 🔥 
### Following a custom Clear Architecture pattern, shrinking the amount of code by a **staggering 33%**? 🤯🥳🥵🤪🤩

Hmmm...

...

...

### **`> No.`**

## Why? 🧐

The best thing about using **`.proto`** files is that you clearly define your service's specs, requests and responses. Combining this with custom annotations on the protofiles and using **`gRPC-Gateway`** we can:

* Automatically generate an **`HTTP handler`** for each **`gRPC method`**.
* Automatically handle requests' input values *(e.g. An HTTP request's body)* and map them to our data structures. 
* Assert validation rules for each request.
* Automatically generate a **`swagger`** spec.

#### ✅ And so you can basically remove the **`Transport Layer`** altogether from your API's internal architecture.

## Is this a fully working example? 👀

**~Kinda**. We implement a simple **`gRPC & HTTP backend API`** with 2 mock endpoints: **Signup** and **Login**.

* With **`gRPC-Gateway`** we expose our gRPC service as a **`RESTful HTTP API`**, defining routes and verbs with annotations on the **`.proto`** files.
* Then we just generate the gateway code and run it alongside the gRPC server. The gateway will translate HTTP requests to gRPC calls, handling input automatically.
* Although we have a **`Service Layer`** implementing the 2 API methods, they don't do much: **There is no DB or persistent storage**.
* **`Protovalidate`** is used to define input rules on the **`.proto`** files themselves for each request, which we enforce using an interceptor.

## Commands ✍🏼

**`make all`**: Generate and run.

**`make run`**: Hmm... What could this possible be?

**`make gen`**: Based on the .proto files, generate the .pb files and the swagger documentation.

**`make protoc-gen`**: Based on the .proto files, generate the .pb files.

**`make swagger-gen`**: Based on the .proto files, generate the swagger documentation.

## Request lifecycle 🔄

 - **`RegisterUsersServiceHandlerClient`**
 - **`request_UsersService_Signup_0`**
 - **`func (c *usersServiceClient) Signup`**
 - **`func (s *Server) handleStream`**
 - **`func (s *Server) processUnaryRPC`**
 - **`_UsersService_Signup_Handler`**
 - **`NewValidationInterceptor`** (the inside function)
 - If error
   - **`handleHTTPError`**
 - If no error
	 - **`func (api *API) Signup`**
	 - **`func (s *service) Signup`**
	 - (Here the request starts backtracking)
	 - **`_UsersService_Signup_Handler`**
	 - **`httpResponseModifier`**
	 - **`func (s *ServeMux) ServeHTTP`**

## 🐿 @gilperopiola
