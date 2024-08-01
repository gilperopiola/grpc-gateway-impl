package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	basePath := `C:\Users\LUCFR\Desktop\`

	filenamer := &osFilenamer{}
	filenames := filenamer.GetFilenamesInFolder(basePath)

	splitter := &ffmpegMP3Splitter{}
	for _, filename := range filenames {
		if strings.HasSuffix(filename, ".mp3") {
			if err := splitter.SplitMP3File(basePath+filename, 4, basePath); err != nil {
				log.Fatalf("ffmpeg error: %v", err)
			}
		}
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type MP3Splitter interface {
	SplitMP3File(pathToFile string, parts int, outputDir string) error
}

type ffmpegMP3Splitter struct{}

func (ffspl *ffmpegMP3Splitter) SplitMP3File(pathToFile string, parts int, outputDir string) error {

	duration, err := getFileDuration(pathToFile)
	if err != nil {
		return err
	}

	partDuration := duration / float64(parts)

	for part := 0; part < parts; part++ {
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s - Part %d.mp3", strings.TrimRight(filepath.Base(pathToFile), ".mp3"), part+1))
		startTime := partDuration * float64(part)
		cmd := exec.Command("ffmpeg", "-ss", fmt.Sprintf("%f", startTime), "-t", fmt.Sprintf("%f", partDuration), "-i", pathToFile, "-acodec", "copy", outputFile)

		// Capturing standard error
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("ffmpeg error: %v, stderr: %s", err, stderr.String())
		}
	}

	return nil
}

func getFileDuration(pathToFile string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", pathToFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	cleanOutput, ok := strings.CutSuffix(string(output), "\r\n")
	if !ok {
		cleanOutput = string(output)
	}

	duration, err := strconv.ParseFloat(cleanOutput, 64)
	if err != nil {
		if duration, err = strconv.ParseFloat(cleanOutput, 32); err != nil {
			return 0, err
		}
	}

	return duration, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type Filenamer interface {
	GetFilenamesInFolder(path string) []string
}

type osFilenamer struct{}

func (osf *osFilenamer) GetFilenamesInFolder(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var filenames []string
	for _, entry := range entries {
		if info, err := entry.Info(); err == nil && !info.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}

	return filenames
}
