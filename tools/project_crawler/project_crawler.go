package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	workingDir, _ := os.Getwd()
	fmt.Println("Working directory: " + workingDir)
	folderToCrawl := "../../pkg"
	crawlProject(folderToCrawl)
}

func visitAll(folderToCrawl string) filepath.WalkFunc {
	fmt.Printf("\nFile %s\n\n", folderToCrawl)
	fmt.Println("Functions:")

	structSlice := []string{}
	interfaceSlice := []string{}

	parseGoFile(folderToCrawl, &structSlice, &interfaceSlice)

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

// parseGoFile parses a single Go file and prints out structs, interfaces, and functions
func parseGoFile(filename string, structSlice, interfaceSlice *[]string) {
	node, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("error parsing %s: %v\n", filename, err)
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

func crawlProject(folderToCrawl string) {
	err := filepath.Walk(folderToCrawl, visitAll(folderToCrawl))
	if err != nil {
		fmt.Printf("error crawling %q: %v\n", folderToCrawl, err)
	}
}
