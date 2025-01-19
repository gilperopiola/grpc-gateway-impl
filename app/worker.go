package app

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"golang.org/x/exp/rand"
)

var prompts = []string{
	"A friendly small %s, illustrated in a %s style. The composition must be centered on the individual, and the individual must be the main focus of the composition.",
}

var colorPalettes = []string{
	"Tropical Sunset", "Mystic River", "Desert Bloom", "Retro Neon", "Ocean Mist", "Candy Pastels", "Burnt Sienna",
	"Lavender Fields", "Mango Sorbet", "Nightfall", "Rose Gold", "Nordic Frost", "Sunlit Meadow", "Indigo Twilight",
	"Pale Dawn", "Vintage Pastel", "Twilight Haze", "Mint Chocolate", "Autumn Harvest", "Dusty Rose",
}

var modernPeriods = []string{
	"1950s", "1960s", "1970s", "1980s", "1990s", "2000s", "2010s", "2020s", "2030s",
}

var periods = []string{
	"3000BCE", "500BCE", "100AD", "900AD", "1400s", "1600s", "1750s", "1800s", "1850s",
	"1900s", "1950s", "1960s", "1970s", "1980s", "1990s", "2000s", "2010s", "2020s", // "2030s", "2100s",
}

var people = []string{
	"young boy", "young girl", "teenage boy", "teenage girl", "adult man", "adult woman",
	"old man", "old woman", "baby",
}

var professionals = []string{
	"carpenter", "doctor", "nurse", "teacher", "police officer", "firefighter",
	"farmer", "fisherman", "chef", "baker", "mechanic", "architect", "engineer", "artist", "writer", "musician",
	"actor", "actress", "dancer", "athlete", "scientist", "explorer", "sailor", "pirate", "king", "queen",
}

var personalities = []string{
	"easy-going", "serious", "funny", "smart", "silly", "optimistic", "pessimistic", "dumb", "kind", "mean", "selfish",
	"generous", "caring", "greedy", "lazy", "hardworking", "ambitious", "creative", "bold", "shy",
}

var creatures = []string{
	"dragon", "minotaur", "unicorn", "centaur", "mermaid", "merfolk", "hydra", "giant", "cyclops", "golem",
	"alien", "elf", "dwarf", "gnome", "halfling", "troll", "ogre", "goblin", "titan", "demon", "angel", "demon lord", "arcangel",
}

var poses = []string{
	"standing up", "sitting down", "jumping", "raising a fist",
	"laughing", "crying", "sleeping", "waking up",
	"reading a book", "writing a letter", "watching tv",
}

func runGPTWorker(app *App) {
	log.Println("DALL·E 3 worker...")
	runID := fmt.Sprintf("%06d", rand.Intn(999999))

	for _, period := range modernPeriods {
		promptFormat := "A %s %s %s, illustrated in the monochrome black and white palette of a modern shonen manga. This %s individual is the main focus of the composition."

		for i := 1; i <= 4; i++ {
			personality := personalities[rand.Intn(len(personalities))]
			professional := professionals[rand.Intn(len(professionals))]

			prompt := fmt.Sprintf(promptFormat, period, personality, professional, professional)
			fileName := fmt.Sprintf("%s - %s %s - %s.png", runID, period, professional, personality)
			log.Printf("Generating [%s]", fileName)

			ctx := context.Background()
			svcResp, err := app.Service.NewGPTImage(ctx, &pbs.NewGPTImageRequest{Message: prompt, Size: pbs.GPTImageSize_SMALL})
			if err != nil {
				time.Sleep(30 * time.Second)
				log.Printf("Error generating [%s]: %v", prompt, err)
				continue
			}

			if svcResp.ImageUrl != "" {
				img, err := downloadImage(svcResp.ImageUrl)
				if err != nil {
					log.Printf("Error downloading [%s]: %v", svcResp.ImageUrl, err)
				}
				saveImage(img, fileName)
				log.Printf("Generated [%s]", fileName)
			}

			time.Sleep(8 * time.Second)
		}
	}
}

var artStyles = []string{
	"monochrome modern shonen manga", "monochrome 1995 romance manga", "vaporwave aesthetic",
	"surrealist painting", "pop art collage", "impressionist watercolor", "abstract geometric",
	"digital glitch art", "ancient chinese ink engravings", "minimalist line art", "fantasy concept art",
	"bold cell shaded", "futuristic sci-fi", "cyberpunk comic panel", "vintage retro poster",
	"expressionist strokes", "stained glass mosaic", "cutting-edge cartooney", "psychedelic tie-dye pattern",
	"noir detective scene", "pixel art sprite",
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

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return png.Decode(resp.Body)
}

func saveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}
