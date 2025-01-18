package app

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

var artStyles = []string{
	"monochrome modern shonen manga", "monochrome 1995 romance manga", "vaporwave aesthetic",
	"surrealist painting", "pop art collage", "impressionist watercolor", "abstract geometric",
	"digital glitch art", "ancient chinese ink engravings", "minimalist line art", "fantasy concept art",
	"bold cell shaded", "futuristic sci-fi", "cyberpunk comic panel", "vintage retro poster",
	"expressionist strokes", "stained glass mosaic", "cutting-edge cartooney", "psychedelic tie-dye pattern",
	"noir detective scene", "pixel art sprite",
}

var bestArtStyles = []string{
	"monochrome modern shonen manga", "monochrome 1995 romance manga", "vaporwave aesthetic",
}

var colorPalettes = []string{
	"Tropical Sunset", "Mystic River", "Desert Bloom", "Retro Neon", "Ocean Mist", "Candy Pastels", "Burnt Sienna",
	"Lavender Fields", "Mango Sorbet", "Nightfall", "Rose Gold", "Nordic Frost", "Sunlit Meadow", "Indigo Twilight",
	"Pale Dawn", "Vintage Pastel", "Twilight Haze", "Mint Chocolate", "Autumn Harvest", "Dusty Rose",
}

var people = []string{
	"1990s teenage girl", "2000s teenage girl", "2010s teenage girl", "2020s teenage girl",
	"1990s teenage boy", "2000s teenage boy", "2010s teenage boy", "2020s teenage boy",
	"1990s adult woman", "2000s adult woman", "2010s adult woman", "2020s adult woman",
	"1990s adult man", "2000s adult man", "2010s adult man", "2020s adult man",
}

var creatures = []string{
	"dragon", "minotaur", "unicorn", "centaur", "mermaid", "merfolk", "hydra", "giant", "cyclops", "golem",
	"alien", "elf", "dwarf", "gnome", "halfling", "troll", "ogre", "goblin", "titan", "demon", "angel", "demon lord", "arcangel",
}

var uiElements = []string{
	"OK button", "Cancel button", "Submit button", "Login button", "Logout button", "Search button", "Refresh button",
	"Add button", "Remove button", "Edit button", "Delete button", "Save button",
}

var poses = []string{
	"standing up", "sitting down", "jumping", "raising a fist",
	"laughing", "crying", "sleeping", "waking up",
	"reading a book", "writing a letter", "watching tv",
}

var prompts = []string{}

func runGPTWorker(app *App) {

	log.Println("Running GPT worker...")
	runID := rand.Intn(999999)
	runIDStr := fmt.Sprintf("%06d", runID)

	for _, creature := range creatures {
		for _, style := range artStyles {
			ctx := context.Background()
			log.Printf("About to generate [Content: %s] [Style: %s]", creature, style)

			prompt := fmt.Sprintf(prompts[0], creature, style)
			svcResp, err := app.Service.NewGPTImage(ctx, &pbs.NewGPTImageRequest{Message: prompt})
			if err != nil {
				log.Printf("Error generating [Prompt: %s]: %v", prompt, err)
				continue
			}

			if svcResp.ImageUrl != "" {
				img, err := downloadImage(svcResp.ImageUrl)
				if err != nil {
					log.Printf("Error downloading image [URL: %s]: %v", svcResp.ImageUrl, err)
				}

				fileName := fmt.Sprintf("%s - %s - %s.png", runIDStr, creature, style)
				saveImage(img, fileName)
				log.Printf("Saved image [Filename: %s]", fileName)
			}

			time.Sleep(10 * time.Second)
		}
	}
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
