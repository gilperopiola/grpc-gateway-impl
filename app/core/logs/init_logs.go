package logs

import (
	"log"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

func Step(n int, name string) {
	log.Printf("%s  ────────── %s\n", numbersToEmojis[n], name)
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
