package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - App -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

type App struct {
	*core.Config // 🗺️
	core.Servers // 🌐
	core.Service // 🌟
	core.Toolbox // 🛠️
}

var _ core.Servers = (*servers.Servers)(nil)
var _ core.Service = (*service.Service)(nil)
var _ core.Toolbox = (*toolbox.Toolbox)(nil)

func NewApp() (runAppFunc, cleanupFunc) {

	// -> ⭐️ This is our core App. This holds everything.
	app := &App{
		Config:  &core.Config{},     // 🗺️
		Servers: &servers.Servers{}, // 🌐
		Service: &service.Service{}, // 🌟
		Toolbox: &toolbox.Toolbox{}, // 🛠️
	}

	func() {
		app.Config = core.LoadConfig()
		core.SetupLogger(&app.LoggerCfg)
	}()

	func() {
		app.Toolbox = toolbox.Setup(app.Config)
		app.Service = service.Setup(app.Toolbox)
		app.Servers = servers.Setup(app.Service, app.Toolbox)
	}()

	func() {
		app.Toolbox.AddCleanupFunc(app.CloseDB)
		app.Toolbox.AddCleanupFunc(app.Servers.Shutdown)
		app.Toolbox.AddCleanupFuncWithErr(zap.L().Sync)
	}()

	return app.Servers.Run, app.Toolbox.Cleanup
}

type runAppFunc func()
type cleanupFunc func()

// -> Our App is like a little house 🏠
//
// It has 2 main entrances:
// -> One for GRPC 🪟 and another for HTTP 🪟. These doors are the 2 PORTS in which the app is run.
// And each entrance leads to a different room:
// -> The GRPC Room 💻 and the HTTP Room 💻. Our SERVERS struct represents these 2 rooms.
//
// -> So, when someone arrives at the GRPC entrance 🤓 and decides to come in, he only sees a corridor with a door
// at the end that reads SERVICE. That is our GRPC Server, a corridor, a small room leading to our SERVICE.
//
// -> Not so simple though, the GRPC corridor is divided into sections, each with a different purpose.
// They are all laid out one after the other, so our guest has to go through all of them to reach the end.
// He slowly starts reading each section's label: 'RATE LIMITER', 'LOGGER', 'TOKEN_VALIDATION', 'PWD_HASHER', etc...
// -> Each section is a GRPC INTERCEPTOR. Each section performs an action based on whoever is trying to get through,
// sometimes blocking his path and making him return back to the entrance with an error.
//
// -> And when he reaches the end of the INTERCEPTORS... Well, there's the door to our beloved SERVICE.
// And as he enters the Service, he is redirected to one of many small but different rooms, each one being a GRPC Service Method.
// There's the Login Room, Signup Room, etc.
//
// -> And then on the Service, he just uses the TOOLBOX stored in there to complete his mission.
// For example, he uses the DB Tool to retrieve and update data, or the Token Generator to get a new JWT on Login.
// In case he cannot continue due to an error, or if he just completed his task, then it returns with the obtained results.
//
// -> He goes back to the GRPC Room, traces back his steps through the Interceptors and goes out of the door.
//
// -> HTTP is another story. The HTTP Room actually has a quite similar structure to the GRPC one. It's a corridor, with a door at
// the end that reads GRPC Room. It's also divided into sections, but these are called MIDDLEWARE instead. And the logic is the same:
// Our guy Ronald 🤓 enters our App through Port :8081, accessing the HTTP Room. There he crosses sections like the CORS Handler,
// and reaches the door to the GRPC Room at the end. He enters, and has to go through all Interceptors, to reach the Service and
// fulfill his mission. When it's over, he heads back to the entrance of the GRPC Room, returns to the HTTP Room (converting the response)
// gotten on the Service from a GRPC format to an HTTP one.
