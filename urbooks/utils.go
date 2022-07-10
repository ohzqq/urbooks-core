package urbooks

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func shellout(args ...string) (error, string, string) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmdArgs := []string{"-c"}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("bash", cmdArgs...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	fmt.Println(cmd.String())
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func FindCover() string {
	if img := findFile(".jpg"); len(img) > 0 {
		return img[0]
	}
	return ""
}

func findFile(ext string) []string {
	var files []string

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range entries {
		if filepath.Ext(f.Name()) == ext {
			file, err := filepath.Abs(f.Name())
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, file)
		}
	}
	return files
}
