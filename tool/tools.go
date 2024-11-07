package main

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func addHtmlFile(desFile *os.File, path string, prefix string) {
	webDir, err := os.Open(path)
	if err == nil {
		files, err2 := webDir.Readdir(-1)
		if err2 == nil {
			var name string
			for _, file := range files {
				name = file.Name()
				if name[0] == '.' {
					continue
				}
				if file.IsDir() {
					addHtmlFile(desFile, path+"/"+name, prefix+"/"+name)
					continue
				}
				desFile.Write([]byte("\tfiles[" + "\"" + prefix + "/" + name + "\"" + "] = "))
				con, err3 := os.ReadFile(path + "/" + name)
				if err3 == nil {
					src := hex.EncodeToString(con)
					desFile.Write([]byte("\""))
					desFile.Write([]byte(src))
					desFile.Write([]byte("\""))
					desFile.Write([]byte("\r\n\r\n"))
				}
			}
		}
	}
}

func addAbiFile(desFile *os.File, path string) {
	webDir, err := os.Open(path)
	if err == nil {
		files, err2 := webDir.Readdir(-1)
		if err2 == nil {
			var name string
			for _, file := range files {
				name = file.Name()
				if name[0] == '.' {
					continue
				}
				if file.IsDir() {
					continue
				}
				ln := len(name) - 4
				if ln <= 0 {
					continue
				}
				varname := name[:ln]
				desFile.Write([]byte("\tfiles[" + "\"" + varname + "\"" + "] = "))
				con, err3 := os.ReadFile(path + "/" + name)
				if err3 == nil {
					src := hex.EncodeToString(con)
					desFile.Write([]byte("\""))
					desFile.Write([]byte(src))
					desFile.Write([]byte("\""))
					desFile.Write([]byte("\r\n\r\n"))
				}
			}
		}
	}
}

func htmlToMap(dir string) {
	f, err := os.OpenFile(dir+"/../webfile.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		fmt.Println("open err")
	}
	f.Write([]byte("package main\r\n\r\n"))
	f.Write([]byte("import \"encoding/hex\"\r\n\r\n"))
	f.Write([]byte("var G_Files = make(map[string][]byte)\r\n\r\n"))

	f.Write([]byte("func InitHtmlFile(){\r\n\r\n"))
	f.Write([]byte("var files = make(map[string]string)\r\n\r\n"))
	addHtmlFile(f, dir+"/../res", "")
	f.Write([]byte("\r\n\r\n"))

	f.Write([]byte("\tfor key, value := range files {\r\n\t\ttmp, _ := hex.DecodeString(value)\r\n\t\tG_Files[key] = tmp\r\n\t}"))
	f.Write([]byte("\r\n\r\n"))

	f.Write([]byte("\tfiles = nil"))
	f.Write([]byte("\r\n"))
	f.Write([]byte("}"))

	f.Write([]byte("\r\n\r\n"))
	f.Close()
}

func abiToMap(dir string) {
	f, err := os.OpenFile(dir+"/../abi.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		fmt.Println("open err")
	}
	f.Write([]byte("package main\r\n\r\n"))
	f.Write([]byte("import \"encoding/hex\"\r\n\r\n"))
	f.Write([]byte("var G_Abis = make(map[string]string)\r\n\r\n"))

	f.Write([]byte("func InitAbiFile(){\r\n\r\n"))
	f.Write([]byte("var files = make(map[string]string)\r\n\r\n"))
	addAbiFile(f, dir+"/../abi")
	f.Write([]byte("\r\n\r\n"))

	f.Write([]byte("\tfor key, value := range files {\r\n\t\ttmp, _ := hex.DecodeString(value)\r\n\t\tG_Abis[key] = string(tmp)\r\n\t}"))
	f.Write([]byte("\r\n\r\n"))

	f.Write([]byte("\tfiles = nil"))
	f.Write([]byte("\r\n"))

	f.Write([]byte("}"))

	f.Write([]byte("\r\n\r\n"))
	f.Close()
}

func main() {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := strings.Replace(path, "\\", "/", -1)
	htmlToMap(dir)
	abiToMap(dir)
}
