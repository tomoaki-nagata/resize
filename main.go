package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/h2non/filetype"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	const quality = 95
	fmt.Println("")
	appPath, err := os.Executable()
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}
	appName := filepath.Base(appPath)
	appDir := filepath.Dir(appPath)
	srcDirName, width, height, err := getParams(appName)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}
	srcDirPath := filepath.Join(appDir, srcDirName)
	if !isDir(srcDirPath) {
		fmt.Printf("[ERROR] directory '%s' is not exist.\n", srcDirPath)
		os.Exit(1)
	}
	dstDirName := fmt.Sprintf("%s_%s_%dx%d", getTimestamp(), srcDirName, width, height)
	dstDirPath := filepath.Join(appDir, dstDirName)
	if err := os.MkdirAll(dstDirPath, 0755); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Input: %s\n", srcDirPath)
	fmt.Printf("Output: %s\n", dstDirPath)
	fmt.Printf("Width: %d\n", width)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Quality: %d\n", quality)
	fmt.Println("")
	result := map[string]int{"ok": 0, "skip": 0, "error": 0}
	if err := filepath.WalkDir(srcDirPath, func(srcFilePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if srcFilePath == srcDirPath {
			return nil
		}
		srcFileName := filepath.Base(srcFilePath)
		if srcFileName[0:1] == "." {
			return nil
		}
		if info.IsDir() {
			fmt.Printf("[SKIP]  %s is directory\n", srcFileName)
			result["skip"]++
			return filepath.SkipDir
		}
		if !isImage(srcFilePath) {
			fmt.Printf("[SKIP]  %s is not image\n", srcFileName)
			result["skip"]++
			return nil
		}
		dstFileName := changeExtension(srcFileName, ".jpg")
		dstFilePath := filepath.Join(dstDirPath, dstFileName)
		if err := resize(srcFilePath, dstFilePath, width, height, quality); err != nil {
			fmt.Printf("[ERROR] %v", err)
			result["error"]++
			return nil
		}
		fmt.Printf("[OK]    %s\n", srcFileName)
		result["ok"]++
		return nil
	}); err != nil {
		fmt.Printf("[FATAL] %s", err)
		os.Exit(1)
	}
	fmt.Println("")
	fmt.Printf("OK: %d, SKIP: %d, ERROR: %d\n", result["ok"], result["skip"], result["error"])
}

func getParams(name string) (srcDirName string, width int, height int, err error) {
	params := strings.Split(name, "+")
	if len(params) != 3 {
		err = fmt.Errorf("invalid name (must be 'resize+[directory]+[width]x[height]' (example: 'resize+input+1000x1000'))")
		return
	}
	srcDirName = params[1]
	size := strings.Split(params[2], "x")
	if len(size) != 2 {
		err = fmt.Errorf("invalid name (must be 'resize+[directory]+[width]x[height]' (example: 'resize+input+1000x1000'))")
		return
	}
	width, err = strconv.Atoi(size[0])
	if err != nil || width < 1 || width > 9999 {
		err = fmt.Errorf("invalid width")
		return
	}
	height, err = strconv.Atoi(size[1])
	if err != nil || height < 1 || height > 9999 {
		err = fmt.Errorf("invalid height")
		return
	}
	return
}

func isDir(path string) bool {
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return false
	}
	return true
}

func isImage(path string) bool {
	buf, _ := ioutil.ReadFile(path)
	return filetype.IsImage(buf)
}

func getTimestamp() string {
	now := time.Now()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return now.In(jst).Format("20060102_150405")
}

func resize(srcPath string, dstPath string, width int, height int, quality int) error {
	srcImage, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}
	dstImage := imaging.Thumbnail(srcImage, width, height, imaging.Lanczos)
	if err := imaging.Save(dstImage, dstPath, imaging.JPEGQuality(quality)); err != nil {
		return err
	}
	return nil
}

func changeExtension(path string, extention string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + extention
}
