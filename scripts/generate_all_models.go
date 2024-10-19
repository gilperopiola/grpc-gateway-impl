package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/fs"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

// This script generates the AllModels slice in /app/core/shared/models/all_models.go by parsing the models package
// and extracting the names of all model structs. It then writes the updated AllModels slice
// back to all_models.go. This script should be run whenever a new model is added to the models
// package or an existing model is removed.

func main() {

	modelsDir := "."
	allModelsFile := modelsDir + "/all_models.go"

	// Step 1️⃣ - Parse the models package and get all model names
	allModelNames, err := getModelStructNames(modelsDir)
	if err != nil {
		log.Fatalf("Failed to get model names: %v", err)
	}

	// Step 2️⃣ - Generate the new code to append
	allModelsCode := generateAllModelsSliceCode(allModelNames)

	// Step 3️⃣ - Update the all_models.go file with the new code
	err = updateAllModelsFile(allModelsFile, allModelsCode)
	if err != nil {
		log.Fatalf("Failed to update all_models.go: %v", err)
	}
}

// 🔽 Implementation

// Crawls through the models pkg, returns all model struct names in a slice.
func getModelStructNames(modelsFolder string) ([]string, error) {
	var out []string

	pkg, err := getPackage(modelsFolder, "models")
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	// For each file
	for _, file := range pkg.Syntax {

		// get all things declared in the top-level scope,
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			// if it's a new type that is being declared
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// and that type is a struct, then we append :)
				if _, ok := typeSpec.Type.(*ast.StructType); ok {
					out = append(out, typeSpec.Name.Name)
				}
			}
		}
	}

	return out, nil
}

// Generates the code for the AllModels slice
func generateAllModelsSliceCode(modelNames []string) string {
	var sb strings.Builder

	sb.WriteString("// DO NOT EDIT this slice manually, just run go generate ./...\n")
	sb.WriteString("// and any model defined in this package should be added automatically.\n")
	sb.WriteString("var AllModels = []any{\n")

	for _, modelName := range modelNames {
		sb.WriteString("\t&" + modelName + "{},\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// Reads the existing all_models.go file, removes the old auto-generated AllModels code
// and writes the new code, the AllModels slice with the updated model names
func updateAllModelsFile(allModelsFile string, allModelsCode string) error {

	// Read the existing all_models.go file
	fileBytes, err := os.ReadFile(allModelsFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	fileString := string(fileBytes)

	// Find the index where the auto-generated code starts
	startOfAutoGenCode := "// DO NOT EDIT this slice manually, just"
	appendIndex := strings.Index(fileString, startOfAutoGenCode)

	// If the auto-generated code is found, remove everything from that point onwards
	if appendIndex != -1 {
		fileString = fileString[:appendIndex]
	} else {
		fileString += "\n"
	}

	// Combine the content and the generated code
	fileString += allModelsCode

	// Fmt it
	newFileBytes, err := format.Source([]byte(fileString))
	if err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	// Replace the content of the file with the new code
	if err = os.WriteFile(allModelsFile, newFileBytes, fs.FileMode(0644)); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Helper. Parses a go pkg and returns it as a *packages.Package
func getPackage(dir, pkgName string) (*packages.Package, error) {
	cfg := packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps,
		Dir:  dir,
	}

	pkgs, err := packages.Load(&cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in directory %s", dir)
	}

	if pkgs[0].Name != pkgName {
		return nil, fmt.Errorf("package '%s' not found", pkgName)
	}

	return pkgs[0], nil
}
