package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
)

func RunAll(service *service.Service) {
	//go RunDallEWorker(service)
}

// ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— ——— —> DALL-E Worker

var colorPalettes = []string{
	"Tropical Sunset", "Mystic River", "Desert Bloom", "Retro Neon", "Ocean Mist", "Candy Pastels", "Burnt Sienna",
	"Lavender Fields", "Mango Sorbet", "Nightfall", "Rose Gold", "Nordic Frost", "Sunlit Meadow", "Indigo Twilight",
	"Pale Dawn", "Vintage Pastel", "Twilight Haze", "Mint Chocolate", "Autumn Harvest", "Dusty Rose",
}

var modernPeriods = []string{"1950s", "1960s", "1970s", "1980s", "1990s", "2000s", "2010s", "2020s", "2030s"}

var periods = []string{
	"3000BCE", "500BCE", "100AD", "900AD", "1400s", "1600s", "1750s", "1800s", "1850s",
	"1900s", "1950s", "1960s", "1970s", "1980s", "1990s", "2000s", "2010s", "2020s", "2030s", "2100s",
}

var people = []string{
	"young boy", "young girl", "teenage boy", "teenage girl", "adult man", "adult woman", "old man", "old woman", "baby",
}

var jobs = []string{
	"carpenter", "doctor", "nurse", "teacher", "police officer", "firefighter", "farmer", "fisherman", "chef", "baker", "mechanic", "architect",
	"engineer", "artist", "writer", "musician", "actor", "actress", "dancer", "athlete", "scientist", "explorer", "sailor", "pirate", "king", "queen",
}

var personalityTraits = []string{
	"easy-going", "serious", "funny", "smart", "silly", "optimistic", "pessimistic", "dumb", "kind", "mean", "selfish",
	"generous", "caring", "greedy", "lazy", "hardworking", "ambitious", "creative", "bold", "shy",
}

var creatures = []string{
	"dragon", "minotaur", "unicorn", "centaur", "mermaid", "merfolk", "hydra", "cyclops", "golem",
	"alien", "elf", "dwarf", "gnome", "halfling", "troll", "ogre", "goblin", "titan", "demon", "angel", "demon lord", "arcangel",
}

var poses = []string{
	"standing up", "sitting down", "jumping", "raising a fist", "laughing", "crying", "sleeping", "waking up", "reading a book", "writing a letter", "watching tv",
}

var moods = []string{"furious", "frustrated", "preoccupied", "worried", "frightened", "asian"}

var who = []string{"police officer", "boxer", "mother", "rapper", "musician", "clown", "salesman",
	"businessman", "athlete", "firefighter", "explorer", "pirate", "baron", "brat", "monarch", "minotaur", "greek gorgon"}

func RunDallEWorker(service *service.Service) {
	runID := new6DigitID()
	logs.InitModuleOK("DALL-E Worker", "🧑‍🏭")

	for _, mood := range moods {
		for _, who := range who {
			promptFmt := "an illustration of a %s %s, front-facing head shot, in ball-point-pen art style"
			prompt := fmt.Sprintf(promptFmt, mood, who)

			log.Printf("Generating %s...", prompt)

			ctx := context.Background()
			resp, err := service.NewGPTImage(ctx, &pbs.NewGPTImageRequest{Message: prompt, Size: pbs.GPTImageSize_DEFAULT})
			if err != nil {
				log.Printf("Error generating %s: %v", prompt, err)
				time.Sleep(25 * time.Second)
				continue
			}
			if resp.ImageUrl != "" {
				img, err := downloadImage(resp.ImageUrl)
				if err != nil {
					log.Printf("Error downloading %s: %v", resp.ImageUrl, err)
				}
				fileName := fmt.Sprintf("%s - %s %s.png", runID, mood, who)
				saveImage(img, fileName)
				log.Printf("Created %s!", fileName)
			} else {
				log.Printf("No image found in %+v", resp)
			}

			time.Sleep(12 * time.Second)
		}
	}
	/*for _, style := range artStyles {
		for i := 0; i < 8; i++ {
			job := jobs[rand.Intn(len(jobs))]
			personality := personalityTraits[rand.Intn(len(personalityTraits))]

			promptFmt := "A %s %s, illustrated in a %s style. This %s individual is the main focus of the composition."
			prompt := fmt.Sprintf(promptFmt, personality, job, style, job)

			log.Printf("Generating %s...", prompt)

			ctx := context.Background()
			resp, err := service.NewGPTImage(ctx, &pbs.NewGPTImageRequest{Message: prompt, Size: pbs.GPTImageSize_DEFAULT})
			if err != nil {
				log.Printf("Error generating [%s]: %v", prompt, err)
				time.Sleep(25 * time.Second)
				continue
			}
			if resp.ImageUrl != "" {
				img, err := downloadImage(resp.ImageUrl)
				if err != nil {
					log.Printf("Error downloading [%s]: %v", resp.ImageUrl, err)
				}
				fileName := fmt.Sprintf("%s - %s %s - %s.png", runID, personality, job, style)
				saveImage(img, fileName)
				log.Printf("Created %s!", fileName)
			} else {
				log.Printf("No image found in %+v", resp)
			}

			time.Sleep(12 * time.Second)
		}
	}*/
}

var artStyles = []string{
	"vibrant cell shaded concept art",
	"realistic cell shaded concept art",
	"detailed cell shaded concept art",

	"contemporary slice-of-life manga",
	"contemporary slice-of-life shonen manga",
	"intrincate slice-of-life manga",
	"contemporary sports manga",

	"cute detailed chibi anime",
	"simple cute chibi anime",

	"sharp-angled bold JRPG illustration",
	"vaporwave glitch aesthetic over a white background",
	"1700AD ink-over-paper intrincate",
	"vintage 1950s poster",
	"next-gen videogame concept art",

	//"monochrome modern shonen manga", "monochrome 1995 romance manga", "vaporwave aesthetic",
	//"surrealist painting", "pop art collage", "impressionist watercolor", "abstract geometric",
	//"digital glitch art", "ancient chinese ink engravings", "minimalist line art", "fantasy concept art",
	//"bold cell shaded", "futuristic sci-fi", "cyberpunk comic panel", "vintage retro poster",
	//"expressionist strokes", "stained glass mosaic", "cutting-edge cartooney", "psychedelic tie-dye pattern",
	//"noir detective scene", "pixel art sprite",
}

var uiElements = []string{
	"OK button", "Cancel button", "Submit button", "Login button", "Logout button", "Search button", "Refresh button",
	"Add button", "Remove button", "Edit button", "Delete button", "Save button",
}

var oldPrompts = []string{
	"A friendly small %s, illustrated in a %s style. The composition must be centered on the individual, and the individual must be the main focus of the composition.",
	"A teenage-aged goth/alternative minotaur on an oversize shirt and some sneakers, in a %s style.",
	"A teenage-aged goth/alternative centaur on an oversize shirt and some sneakers, in a %s style.",
}
