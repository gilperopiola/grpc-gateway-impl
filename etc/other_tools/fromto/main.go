package main

import (
	"log"
	"os/exec"
)

func main() {
	basePath := `C:\Users\LUCFR\Desktop\`
	fileName := `Webhooks`

	converter := &ffmpegConverter{}
	converter.FromMP4ToMP3(basePath+fileName+`.mp4`, basePath+fileName+`.mp3`)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type FromtoConverter interface {
	FromMP4ToMP3(inputPath, outputPath string)
}

type ffmpegConverter struct{}

func (ffconv *ffmpegConverter) FromMP4ToMP3(inputPath, outputPath string) {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-q:a", "0", "-map", "a", outputPath, "-y")

	if err := cmd.Run(); err != nil {
		log.Fatalf("ffmpeg error: %v", err)
	}
}
