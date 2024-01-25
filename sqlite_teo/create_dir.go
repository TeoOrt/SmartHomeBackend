package sqlite_teo

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDir(user_email *string) (string, error) {
	root := "/home/teoortega/AndroidStudioProjects/SmartHomeBackend/video_storage"
	path := filepath.Join(root, *user_email)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		fmt.Println("Could not create dir", err)
		return "", err
	}
	return path, nil
}
