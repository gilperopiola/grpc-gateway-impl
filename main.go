package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app"
)

// 🔻 - --------------------------------------- - 🔻
// 🔻 - --- - GRPC Gateway Implementation - --- - 🔻
// 🔻 - --------------------------------------- - 🔻

func main() {
	runApp, cleanUp := app.Setup()
	defer cleanUp()

	runApp()

	waitForSignal := make(chan os.Signal, 1)
	signal.Notify(waitForSignal, stopSignals...)
	<-waitForSignal
}

var stopSignals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM}

// ---             ╭───────────────╮                         ╭───────────╮
// GRPC Request -> │  GRPC Server  │ -> Interceptor Chain -> │  Service  │ -> DB/API Calls ~ Tools ~ Etc ╮
//                 ╰───────────────╯                         ╰───────────╯                               |
// 		     			   ↑                                                                             ↓
//		     			   |													 	           Service GRPC Response
//                         |                                                                             |
//                  Middleware Chain                            ╭───────────────╮                        |
//                     ↑        |╰ -- -- -- -- -- - GRPC OUT <- │  GRPC Server  │ <- Interceptor Chain <-╯
//                     |        |      if http                  ╰───────────────╯
//                     |        |
//                     |        ↓
//                 ╭───────────────╮
// HTTP Request -> │  HTTP Server  │ -> HTTP OUT
// ---             ╰───────────────╯
