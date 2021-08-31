package main

import (
	"bufio"
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
	cmd := "find . -name " //fileGo := `find . -name "*.go"`
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

func getOldName() (oldName string) { // 读取 go mod文件
	oldName = ""
	fs, err := os.Open("go.mod")
	if err != nil {
		return
	}

	defer fs.Close()
	buff := bufio.NewReader(fs)
	content, _, err := buff.ReadLine()
	if err != nil {
		return
	}

	subs := strings.Split(string(content), " ")
	oldName = subs[1]
	return oldName
}

func getNewName(args []string) string {
	return args[1]
}

func main() {
	if len(os.Args) != 2 { // 读取命令行参数
		log.Fatal("args incorrect")
	}

	newName := getNewName(os.Args)
	oldName := getOldName()

	// find files *.go, *.proto go.mod
	goFile := `"*.go"`
	pbFile := `"*.proto"`
	modFile := "go.mod"
	typeFiles := []string{goFile, pbFile, modFile}

	var allFiles []string
	for _, str := range typeFiles {
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

	// 替换文件
	for _, file := range allFiles {
		err := replaceFile(file, oldName, newName)
		if err != nil {
			log.Fatal("rewrite file failed: ", err)
		}
	}
}
