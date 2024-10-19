package logs

import (
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"go.uber.org/zap"
)

func Step(step int) {
	log.Printf("%s\n", numbersToEmojis[step])
}

func SubstepOK(name, emoji string) {
	log.Printf("\t %s %s OK\n", emoji, name)
}

func EnvVars() {
	log.Println(" \t🎈 APP = " + core.G.AppName + " " + core.G.Version)
	log.Println(" \t🌐 ENV = " + core.G.Env)
	if core.G.TLSEnabled {
		log.Println(" \t✅ TLS = on")
	} else {
		log.Println(" \t⚠️  TLS = off")
	}
}

var numbersToEmojis = map[int]string{
	0:  "0️⃣",
	1:  "1️⃣",
	2:  "2️⃣",
	3:  "3️⃣",
	4:  "4️⃣",
	5:  "5️⃣",
	6:  "6️⃣",
	7:  "7️⃣",
	8:  "8️⃣",
	9:  "9️⃣",
	10: "🔟",
}

// On Windows I get a *fs.PathError calling zap.L().Sync() to flush logger on shutdown.
// This just calls zap.L().Sync() and ignores that specific error. See https://github.com/uber-go/zap/issues/991
func SyncLogger() error {
	var pathErr *fs.PathError
	if err := zap.L().Sync(); err != nil && !errors.As(err, &pathErr) {
		return fmt.Errorf("error syncing logger: %w", err)
	}
	return nil
}
