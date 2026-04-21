package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	var srcDir string

	modCache := os.Getenv("GOMODCACHE")
	if modCache != "" {
		srcDir = filepath.Join(modCache, "google.golang.org", "protobuf@v1.36.11", "src", "google", "protobuf")
	}

	if srcDir == "" || !protoFilesExist(srcDir) {
		srcDir = getProtocIncludeDir()
	}

	dstDir := "assets/google-protobuf"

	if !protoFilesExist(srcDir) {
		fmt.Println("Proto files not found, downloading...")
		cmd := exec.Command("go", "mod", "download", "google.golang.org/protobuf@v1.36.11")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Download failed: %v\n", err)
			os.Exit(1)
		}
		if modCache == "" {
			modCache = os.Getenv("GOMODCACHE")
		}
		if modCache != "" {
			srcDir = filepath.Join(modCache, "google.golang.org", "protobuf@v1.36.11", "src", "google", "protobuf")
		}
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dstDir, err)
		os.Exit(1)
	}

	files, err := filepath.Glob(filepath.Join(srcDir, "*.proto"))
	if err != nil {
		fmt.Printf("Error finding proto files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Printf("No proto files found in %s\n", srcDir)
		os.Exit(1)
	}

	for _, srcFile := range files {
		name := filepath.Base(srcFile)
		dstFile := filepath.Join(dstDir, name)

		srcData, err := os.ReadFile(srcFile)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", srcFile, err)
			continue
		}

		if err := os.WriteFile(dstFile, srcData, 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", dstFile, err)
			os.Exit(1)
		}

		fmt.Printf("Copied: %s\n", name)
	}

	fmt.Printf("\nDone! %d proto files copied to %s\n", len(files), dstDir)
}

func protoFilesExist(dir string) bool {
	if dir == "" {
		return false
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.proto"))
	return err == nil && len(files) > 0
}

func getProtocIncludeDir() string {
	var protocPath string

	switch runtime.GOOS {
	case "windows":
		protocPath = filepath.Join(os.Getenv("USERPROFILE"), "scoop", "apps", "protobuf", "current", "include")
	case "linux":
		protocPath = "/usr/include"
	case "darwin":
		protocPath = "/usr/local/include"
	}

	protoDir := filepath.Join(protocPath, "google", "protobuf")
	if protoFilesExist(protoDir) {
		return protoDir
	}

	out, err := exec.Command("go", "env", "GOTOOLSWORLD").Output()
	if err == nil {
		goBin := filepath.Dir(string(out))
		protocPath = filepath.Join(goBin, "..", "..", "pkg", "mod", "cache", "download", "github.com", "protocolbuffers", "protobuf", "raw")
	}

	return filepath.Join(protocPath, "google", "protobuf")
}
