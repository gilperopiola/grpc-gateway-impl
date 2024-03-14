# 🐉 gRPC Gateway Implementation 🐉

### Need to build a Backend API in Go? 😱 Using both gRPC and HTTP? 🔥 
### Following a custom Clean Architecture pattern, shrinking the amount of code by a **staggering 33%**? 🤯🥳🥵🤪🤩

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

#### ✅ And so you can basically *`remove the Transport Layer`* altogether from your API's internal architecture.

## Is this a fully working example? 👀

**~Yeah, actually!**. We implement a simple **`gRPC & HTTP Backend API`** with 4 endpoints: **Signup** and **Login**, **GetUser** and **GetUsers**.

`We got TLS, JWT Auth with Roles, Rate Limiting, Logging. What else do you want, huh?`

## Commands ✍🏼

**`make all-fast`**: Generate files and run app.

**`make run`**: Hmm... What could this possible be?

**`make generate`**: Based on the .proto files, generate the .pb files and the swagger documentation.

**`make generate-pbs`**: Based on the .proto files, generate the .pb files.

**`make generate-swagger`**: Based on the .proto files, generate the swagger documentation.

For more commands, check the Makefile.

## Request lifecycle 🔄

This needs some love, it's pretty outdated:

 - **`RegisterUsersServiceHandlerClient`**
 - **`request_UsersService_Signup_0`**
 - **`func (c *usersServiceClient) Signup`**
 - **`func (s *Server) handleStream`**
 - **`func (s *Server) processUnaryRPC`**
 - **`_UsersService_Signup_Handler`**
 - **`NewValidationInterceptor`** (the inside function)
 - *If error*
   - **`handleHTTPError`**
 - *If no error*
	 - **`func (api *API) Signup`**
	 - **`func (s *service) Signup`**
	 - (Here the request starts backtracking)
	 - **`_UsersService_Signup_Handler`**
	 - **`httpResponseModifier`**
	 - **`func (s *ServeMux) ServeHTTP`**

## 🐿 @gilperopiola
