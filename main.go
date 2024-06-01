package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ */
/* -~-~-~-~-~- GRPC Gateway Implementation -~-~-~-~-~ */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ */

// -> Welcome 🌈

func main() {

	// -> Init app
	runAppFn, cleanupFn := app.NewApp()
	defer cleanupFn()

	// -> Run app
	runAppFn()

	// -> Exit app
	func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-ch
	}()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* Update README / Put tests back in */
/* Buf file / Dockerfile / Docker-compose / Kubernetes /
/* CI-CD / Metrics / Tracing / Caching / Tests / Always obfuscate requests on Logs */
