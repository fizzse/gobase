package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
  重命名项目
*/

func replaceFile(file, oldStr, newStr string) (err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	oldBytes := []byte(oldStr)
	newBytes := []byte(newStr)

	if !bytes.Contains(content, oldBytes) {
		return
	}

	content = bytes.ReplaceAll(content, oldBytes, newBytes)
	err = os.Remove(file)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, content, 0644)
	return
}

func findFiles(name string) ([]string, error) {
	//fileGo := `find . -name "*.go"`
	cmd := "find . -name "
	cmd += name

	cmder := exec.Command("bash", "-c", cmd)
	out, err := cmder.CombinedOutput()
	if err != nil {
		return nil, err
	}

	strOut := string(out)
	files := strings.Split(strOut, "\n")
	return files, nil
}

func main() {
	// 读取命令行参数
	if len(os.Args) != 2 {
		log.Fatal("args incorrectt")
	}

	// find files *.go, *.proto go.mod
	goFile := `"*.go"`
	pbFile := `"*.proto"`
	gomod := "go.mod"
	strs := []string{goFile, pbFile, gomod}

	var allFiles []string
	for _, str := range strs {
		files, err := findFiles(str)
		if err != nil {
			log.Fatal("find file failed: ", err)
		}

		for _, file := range files {
			if file != "" {
				allFiles = append(allFiles, file)
			}
		}
	}

	oldName := "github.com/fizzse/gobase"
	newName := os.Args[1]

	// 替换文件
	for _, file := range allFiles {
		err := replaceFile(file, oldName, newName)
		if err != nil {
			log.Fatal("rewrite file failed: ", err)
		}
	}
}
