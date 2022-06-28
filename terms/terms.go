package terms

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed terms/*
var content embed.FS

func PrintLicenseTerms() {
	fmt.Print("Usage, modification, and distribution of this software and its components are subject to the following respective licensing terms:\n\n")

	_ = fs.WalkDir(content, "terms", MyWalkFunc)
}

func MyWalkFunc(path string, d fs.DirEntry, err error) error {
	if !d.IsDir() && d.Name() != "make_embedFS_happy" {
		components := strings.Split(path, "/")
		components = components[1 : len(components)-1]
		cleanPath := filepath.Join(components...)

		fileContent, _ := fs.ReadFile(content, path)

		fmt.Printf("License for %s:\n", cleanPath)
		fmt.Println("================================================================================")
		fmt.Print(string(fileContent))
		fmt.Println("================================================================================")
		fmt.Print("\n\n\n")
	}
	return nil
}
