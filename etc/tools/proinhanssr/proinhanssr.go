package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// T0D0 make all enum variables a diff color

func main() {
	workingDir, _ := os.Getwd()
	fmt.Println("Working directory: " + workingDir)
	folderToCrawl := "./pkg"
	beginCrawling(folderToCrawl)
}

func beginCrawling(folderToCrawl string) {
	err := filepath.Walk(folderToCrawl, visitAll())
	if err != nil {
		fmt.Printf("error crawling %q: %v\n", folderToCrawl, err)
	}
}

func visitAll() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		if strings.Contains(path, "pb.") {
			return nil
		}

		if strings.Contains(path, "mock") {
			return nil
		}

		if strings.Contains(path, "db_adapter.go") {
			return nil
		}

		fmt.Printf("\nFile %s\n\n", path)
		fmt.Println("Functions:")
		structSlice := []string{}
		interfaceSlice := []string{}

		parseGoFile(path, &structSlice, &interfaceSlice)

		if len(structSlice) > 0 {
			fmt.Println("\nStructs:")
			for _, strct := range structSlice {
				fmt.Println("    *", strct)
			}
		}

		if len(interfaceSlice) > 0 {
			fmt.Println("\nInterfaces:")
			for _, iface := range interfaceSlice {
				fmt.Println("    *", iface)
			}
		}

		return nil
	}

}

// parseGoFile parses a single Go file and prints out structs, interfaces, and functions
func parseGoFile(filename string, structSlice, interfaceSlice *[]string) {
	node, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		if _, err := parser.ParseDir(token.NewFileSet(), filename, nil, parser.ParseComments); err != nil {
			fmt.Printf("error parsing %s: %v\n", filename, err)
			return
		}
		fmt.Printf("entering directory: %s", filename)
		return
	}

	// Inspect the AST and print all structs, interfaces, and functions
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			switch x.Type.(type) {

			case *ast.StructType:
				*structSlice = append(*structSlice, fmt.Sprintf(".. [%s] %s\n:", filepath.Base(filename), x.Name.Name))

			case *ast.InterfaceType:
				*interfaceSlice = append(*interfaceSlice, fmt.Sprintf(".. [%s] %s\n:", filepath.Base(filename), x.Name.Name))
			}
		case *ast.FuncDecl:
			for _, functionToOmit := range protoFunctions {
				if x.Name.Name == functionToOmit {
					return true
				}
			}

			if strings.HasSuffix(filename, "db_adapter.go") ||
				strings.HasSuffix(filename, "_mock.go") ||
				strings.Contains(filename, ".pb.") ||
				strings.HasSuffix(filename, "_test.go") {
				return true
			}

			fmt.Printf("..%s[%s] %s\n:", getSpacesRemaining(len(filepath.Base(filename))), filepath.Base(filename), x.Name.Name)
		}
		return true
	})
}

func getSpacesRemaining(length int) string {
	var spaces string
	for i := 0; i < 30-length; i++ {
		spaces += " "
	}
	return spaces
}

var protoFunctions = []string{
	"Reset", "String", "ProtoMessage", "ProtoReflect", "Descriptor",
}

func isInStringsSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func modifySettingsJSONFile() {
	filePath := `C:\Users\LUCFR\AppData\Roaming\Code\User`
	searchString := "// proinhanssr begin"
	appendString := `"ASDASD", {"text": "ASDASD", "color": "#FF0000"}`

	// Read the content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}

	// Convert content to string and append the new string after the specific string
	contentStr := string(content)
	index := strings.Index(contentStr, searchString)
	if index == -1 {
		log.Fatalf("string '%s' not found in file", searchString)
	}

	// Add the appendString immediately after the searchString
	modifiedContent := contentStr[:index+len(searchString)] + appendString + contentStr[index+len(searchString):]

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}
