package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func init() {
	// Configure slog for human-readable output during development
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	rootDir := "../" // directory to scan recursively

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("Error accessing path",
				"path", path,
				"error", err)
			return nil // continue walking
		}

		if info.IsDir() {
			return nil // skip directories
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			processFile(path)
		}
		return nil
	})

	if err != nil {
		slog.Error("Error walking the path",
			"rootDir", rootDir,
			"error", err)
	}
}

func processFile(filename string) {
	// First, use AST to find functions and types that need comments
	fset := token.NewFileSet()
	src, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Failed to read file",
			"filename", filename,
			"error", err)
		return
	}

	node, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		slog.Error("Failed to parse file",
			"filename", filename,
			"error", err)
		return
	}

	// Find functions and types that need comments
	functionsNeedingComments := make(map[string]bool)
	typesNeedingComments := make(map[string]bool)

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Handle both regular functions and methods
			if !d.Name.IsExported() {
				continue // skip unexported functions
			}
			if d.Doc == nil || len(d.Doc.List) == 0 {
				functionsNeedingComments[d.Name.Name] = true
				if d.Recv != nil {
					slog.Info("🎯 Found method needing comment",
						"method", d.Name.Name)
				} else {
					slog.Info("🎯 Found function needing comment",
						"function", d.Name.Name)
				}
			} else {
				if d.Recv != nil {
					slog.Debug("ℹ️ Skipped method: comment already exists",
						"method", d.Name.Name)
				} else {
					slog.Debug("ℹ️ Skipped function: comment already exists",
						"function", d.Name.Name)
				}
			}

		case *ast.GenDecl:
			// Handle type declarations
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if !typeSpec.Name.IsExported() {
							continue // skip unexported types
						}
						if typeSpec.Doc == nil || len(typeSpec.Doc.List) == 0 {
							// Check if the general declaration has a comment
							if d.Doc == nil || len(d.Doc.List) == 0 {
								typesNeedingComments[typeSpec.Name.Name] = true
								slog.Info("🎯 Found type needing comment",
									"type", typeSpec.Name.Name)
							} else {
								slog.Debug("ℹ️ Skipped type: comment already exists on declaration",
									"type", typeSpec.Name.Name)
							}
						} else {
							slog.Debug("ℹ️ Skipped type: comment already exists",
								"type", typeSpec.Name.Name)
						}
					}
				}
			}
		}
	}

	if len(functionsNeedingComments) == 0 && len(typesNeedingComments) == 0 {
		slog.Info("ℹ️ No functions or types need comments",
			"filename", filename)
		return
	}

	// Now use string manipulation to add comments
	lines := strings.Split(string(src), "\n")
	modified := false

	// Regex to match function declarations (including methods)
	funcRegex := regexp.MustCompile(`^func\s+(?:\([^)]*\)\s+)?([A-Z][a-zA-Z0-9_]*)\s*\(`)

	// Regex to match type declarations
	typeRegex := regexp.MustCompile(`^type\s+([A-Z][a-zA-Z0-9_]*)\s+`)

	// Process lines in reverse order to avoid index shifting issues
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Process function declarations (including methods)
		if strings.HasPrefix(trimmed, "func ") {
			matches := funcRegex.FindStringSubmatch(trimmed)
			if len(matches) > 1 {
				funcName := matches[1]
				if functionsNeedingComments[funcName] {
					if !hasCommentAbove(lines, i) {
						comment := createComment(line, funcName)
						lines = insertComment(lines, i, comment)
						modified = true
						slog.Info("✅ Added comment to function/method",
							"function", funcName,
							"line", i+1)
						delete(functionsNeedingComments, funcName)
					}
				}
			}
		}

		// Process type declarations
		if strings.HasPrefix(trimmed, "type ") {
			matches := typeRegex.FindStringSubmatch(trimmed)
			if len(matches) > 1 {
				typeName := matches[1]
				if typesNeedingComments[typeName] {
					if !hasCommentAbove(lines, i) {
						comment := createComment(line, typeName)
						lines = insertComment(lines, i, comment)
						modified = true
						slog.Info("✅ Added comment to type",
							"type", typeName,
							"line", i+1)
						delete(typesNeedingComments, typeName)
					}
				}
			}
		}
	}

	if modified {
		// Write the modified content back
		newContent := strings.Join(lines, "\n")
		err = os.WriteFile(filename, []byte(newContent), 0644)
		if err != nil {
			slog.Error("Failed to write file",
				"filename", filename,
				"error", err)
			return
		}
		slog.Info("📝 Updated file",
			"filename", filename)
	}
}

// hasCommentAbove checks if there's already a comment above the given line
func hasCommentAbove(lines []string, lineIndex int) bool {
	for j := lineIndex - 1; j >= 0; j-- {
		prevLine := strings.TrimSpace(lines[j])
		if prevLine == "" {
			continue // skip empty lines
		}
		if strings.HasPrefix(prevLine, "//") {
			return true
		}
		break // stop at first non-empty line
	}
	return false
}

// createComment creates a properly indented comment for the given line and name
func createComment(line, name string) string {
	// Find the proper indentation by looking at the line
	indent := ""
	for _, char := range line {
		if char == ' ' || char == '\t' {
			indent += string(char)
		} else {
			break
		}
	}
	return indent + fmt.Sprintf("// %s TODO: Add description", name)
}

// insertComment inserts a comment line before the specified line index
func insertComment(lines []string, index int, comment string) []string {
	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:index]...)
	newLines = append(newLines, comment)
	newLines = append(newLines, lines[index:]...)
	return newLines
}
