package sqlite_teo

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDir(last_name *string) (string, error) {
	root, _ := os.Getwd()
	path := filepath.Clean(filepath.Join(root, *last_name))
	err := os.MkdirAll(path, os.ModePerm)
	// fmt.Printf("Directory Path: %s \n", filepath.Join(root, *last_name))

	if err != nil {
		fmt.Println("Could not create dir", err)
		return "", err
	}
	return path, nil
}
