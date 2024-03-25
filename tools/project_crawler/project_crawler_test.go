package main

import (
	"fmt"
	"os"
	"testing"
)

func TestCrawl(t *testing.T) {
	workingDir, _ := os.Getwd()
	fmt.Println("Working directory: " + workingDir)
	folderToCrawl := "../../pkg"
	crawlProject(folderToCrawl)
}
