# 🐉 GRPC Gateway Implementation 🐉

### **If an HTTP API is worth a dollar, its GRPC counterpart is like 20 cents more**.
<div style="text-align: right; margin-top: -4px">- Lionel Messi. </div>

---
(This README needs updating)

**And now** with _`GRPC Gateway`_ we get both of them for only... _**$0.75!!!**_ 😱 With extra fries and a refreshing _Monster Mango Loco™_ free of charge. 🍟🥤 

... _**What?**_ 🤨 You're not making any sense, why do you use emojis on a _**readme**_ ❗❓ 

**Let's say** we define a `UsersService` in a `.proto` file, with some endpoints. We can _leverage_ `GRPC Gateway` to auto-generate a **GRPC Server** with handlers mapping each endpoint on our previously defined `UsersService`. It also auto-generates a**n** **HTTP Gateway** that points to the server and translates _HTTP_ to _gRPC_ and viceversa. 🤯 _**For free!**_

... ... ... But wasn't it... _$0.75_? 🤔

**Yes, _BUUUT_** if you buy _~now~_ we'll throw in the **BEAUTIFUL. CHARMING. BREATHTAKING.........** **`HTTP Swagger Spec`** for your API. And I think I don't even need to say this, but I'll say it anyways... Th-The Swagger... The Swagger Spec. The Swagger Spec... is... au-tO-**GE-nE-RRA-ttTTTEDD** from annotations on the `.proto`. 🎉 Even request validations are configured on the protofile, so there is practically _no Transport Lyo ayer~_.

  ️️👁️👁️  **...** 🤔 **...** OK **...** 👁️🤔 ... That's actually... cool? ... Take my money.


# So what's in here? 👀

_**~I'm glad you ask!**_ - We have **`two simple APIs: 1 GRPC & 1 HTTP`**, each of them with **4** endpoints --> **Signup**, **Login**, **GetUser** and **GetUsers**.

**_It leverages_** --> Clean, Hexagonal Architecture 🔷 / MySQL 🐬 / Patterns and Good Practices 📐 / Excellent Documentation 📚 / Gorm 🌱 / Centralized Error Handling 🎯 / JWT 🔑 / TLS 🔒 / RBAC 👑 / GCP ❌ / JJR ❓ / Y2K 🤔 / Swagger 📜 / BRB 🤦‍♂️ / LOL 😂 / Postman Automation 📬 / AFK 🏃‍♀️.

# Request lifecycle 🔄

➡️ When a **Signup _HTTP_ Request** hits the Gateway, the first file to be called is **`users.pb.gw.go`** on:

* **RegisterUsersServiceHandlerClient** > **request_UsersService_Signup_0**

➡️ Then, when it needs to go through **`google.golang.org/grpc/server.go`**, it calls **`users_grpc.pb.go`** on:

* **usersServiceClient.Signup** > **_UsersService_Signup_Handler**

➡️ Followed by our interceptors in **`grpc_interceptors.go`**:

* **rateLimiterInterc** > **loggerInterc** > **tokenValidationInterc** > **inputValidationInterc** > **etc...**

➡️ To finally reach our beloved **`service_users.go`** on **service.Signup**.

# Useful Commands ✍🏼

**`make all`**: Makes all.

🤪 It cleans the env, generates code, runs tests, and runs the app.

**`make all fast=1`**: Makes all, but faster. Skips cleaning and testing.

**`make help`**: Shows help message. 

**`make run`**: Updates `go.mod` and runs app.

**`make generate`**: Based on the `.proto` files, generates the `.pb.go` files and Swagger Spec.

For more commands, check the `Makefile`. 🌈

# Code Generation 🖥️

With **`.protos`**: your _API_ gets defined before it's implemented. 

Using custom annotations on the **`.proto`** file and tools like `GRPC Gateway` you get an _HTTP Handler_ for each _gRPC Method_, each Handler decoding _HTTP_ Requests into _gRPC_ ones, calling their designated method on the _gRPC_ server and encoding the _gRPC_ Response back into _HTTP_.

<div style="margin-bottom: -16px">
You also get a <i>Validation Layer</i> based on the protofile. And did I mention the free Swagger? 😁 It's <i>$11.99</i>.
 <p style="display: inline-block;font-size:8px">Plus taxes.</p>
</div>

## Auto Generated Files 🕸

### **users.pb.go**
**Protobuf Types and their methods, as defined in the .proto.**

SignupRequest / SignupResponse / LoginRequest / LoginResponse / UserInfo / PaginationInfo. 

### **users_grpc.pb.go**

**GRPC Server and Client, endpoints registration.**

<div>type UsersServiceClient interface</div>
<div>type UsersServiceServer interface</div>

<div style='margin-top: 16px'>func RegisterUsersServiceServer(...)</div>

### **users.pb.gw.go**
**Reverse proxy, decodes HTTP into GRPC and viceversa.**

<div>request_UsersService_Signup_0(...)</div>
<div>request_UsersService_Login_0(...)</div>
<div>request_UsersService_GetUser_0(...)</div>
<div>request_UsersService_GetUsers_0(...)</div>

<div style='margin-top: 16px'>RegisterUsersServiceHandler(...)</div>
<div>RegisterUsersServiceHandlerFromEndpoint(...)</div>


## 🐿~Y ya que estás acá...
 [_@gilperopiola_](https://www.instagram.com/gilperopiola/) 🚀