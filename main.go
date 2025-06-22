package main

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app"
)

// 🔻 - --------------------------------------- - 🔻
// 🔻 - ─── ─ GRPC Gateway Implementation ─ ─── - 🔻
// 🔻 - --------------------------------------- - 🔻

func main() {
	runApp, cleanUp := app.Setup()
	defer cleanUp()
	runApp()
	god.WaitForSignal()
}

// —               ╭───────────────╮                         ╭───────────╮
// GRPC Request —> │  GRPC Server  │ —> Interceptor Chain —> │  Service  │ —> DB/API Calls ~ Tools ~ Etc ╮
//                 ╰───────────────╯                         ╰───────────╯                               │
// 		     			   ↑                                                                             │
//		     			   │													 	           Service GRPC Response
//                         │                                                                             │
//                  Middleware Chain                            ╭───────────────╮                        │
//                     ↑        │╰ —— ╮         ╭ [GRPC Out] <— │  GRPC Server  │ <— Interceptor Chain <—╯
//                     │        │     ╰ if http ╯               ╰───────────────╯
//                     │        │
//                     │        ↓
//                 ╭───────────────╮
// HTTP Request —> │  HTTP Server  │ —> [HTTP Out]
// —               ╰───────────────╯
