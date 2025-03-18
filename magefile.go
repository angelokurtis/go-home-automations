//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Aliases maps command names to their corresponding actions.
// It allows using shorthand names for tasks, e.g., "gen" for Generate and "lint" for Lint.All.
var Aliases = map[string]any{
	"build": Build.Prod,
	"gen":   Generate,
	"lint":  Lint.All,
}

// Build defines the set of tasks related to the application's build process.
type Build mg.Namespace

// Prod builds the production version of the application
func (Build) Prod(ctx context.Context) error {
	mg.CtxDeps(ctx, Clean, Generate)

	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        runtime.GOOS,
		"GOARCH":      runtime.GOARCH,
	}
	ldFlags := []string{"-s", "-w"}
	args := []string{"build", "-ldflags", strings.Join(ldFlags, " "), "-o", "bin/go-merge", "./cmd/app/"}

	return sh.RunWithV(env, "go", args...)
}

// Dev builds the development version of the application
func (Build) Dev(ctx context.Context) error {
	mg.CtxDeps(ctx, Clean, Generate)

	args := []string{"build", "-o", "bin/go-merge", "./cmd/app/"}

	return sh.RunV("go", args...)
}

// Clean removes the "bin" directory to ensure a fresh build
func Clean() error {
	directory := "./bin/"

	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil
	}

	// Try to remove the directory
	if err := os.RemoveAll(directory); err != nil {
		return fmt.Errorf("failed to delete %s: %w", directory, err)
	}

	return nil
}

// Generate runs Wire to generate dependency injection code.
func Generate() error {
	dirs, err := scanWireDirs("./")
	if err != nil {
		return err
	}

	// Run Wire codegen tool on each directory found.
	for _, dir := range dirs {
		if err = sh.RunV("go", "run", "-mod=mod", "github.com/google/wire/cmd/wire@latest", dir); err != nil {
			return err
		}
	}

	return nil
}

// Lint defines the set of tasks related to linting the codebase.
type Lint mg.Namespace

// All runs linting on the entire codebase.
func (Lint) All() error {
	if err := sh.RunV("go", "run", "-mod=mod", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest", "run"); err != nil {
		return err
	}

	return nil
}

// Changes runs linting only on modified files.
func (Lint) Changes() error {
	if err := sh.RunV("go", "run", "-mod=mod", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest", "run", "--new"); err != nil {
		return err
	}

	return nil
}

// scanWireDirs finds directories with Go files using wire.Build.
func scanWireDirs(rootDir string) ([]string, error) {
	var dirs []string
	visited := make(map[string]bool)

	// Walk through the directory tree to find Wire Build calls.
	err := filepath.Walk(rootDir, func(filePath string, fileInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if fileInfo.IsDir() || filepath.Ext(filePath) != ".go" {
			return nil
		}

		fs := token.NewFileSet()

		node, parseErr := parser.ParseFile(fs, filePath, nil, parser.AllErrors)
		if parseErr != nil {
			log.Printf("Error parsing file: %s, %v", filePath, parseErr)
			return nil
		}

		ast.Inspect(node, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				if selector, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if selector.Sel.Name == "Build" {
						if ident, ok := selector.X.(*ast.Ident); ok && ident.Name == "wire" {
							dir := fmt.Sprintf("./%s/", filepath.Dir(filePath))
							if !visited[dir] {
								visited[dir] = true

								dirs = append(dirs, dir)
							}

							return false
						}
					}
				}
			}

			return true
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}
