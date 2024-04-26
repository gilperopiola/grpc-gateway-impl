package main

import (
	"github.com/gilperopiola/grpc-gateway-impl/app"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Welcome~! -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* This is the entrypoint of our app */
/* It runs a GRPC Server and points an HTTP Gateway towards it */
/* Even though it's not a complex app, its architecture and overall code are extremely polished */

func main() {

	// Init app
	app := app.NewApp()

	// Run app
	app.Run()

	// Exit app
	app.WaitForShutdown()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - T0D0 -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* Update README / Put tests back in */
/* Buf file / Dockerfile / Docker-compose / Kubernetes /
/* CI-CD / Metrics / Tracing / Caching / Tests / Always obfuscate requests on Logs */

/*

	intToStrByStrconv := func(i int) string { return strconv.Itoa(i + 1) }
	intToStrByFmt := func(i int) string { return fmt.Sprintf("%d", i+1) }
	intToStrByCasting := func(i int) string { return string(i + 1) }

	funcsToBenchmork := []func(int) string{
		intToStrByStrconv, intToStrByFmt, intToStrByCasting,
	}

	benchmorkar.Kumpare(funcsToBenchmork, 1000)

*/
