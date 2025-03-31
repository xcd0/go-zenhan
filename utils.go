package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"
)

// グローバル変数。
var (
	version  string = "debug build"   // makefileからビルドされると上書きされる。
	revision string = func() string { // {{{
		revision := ""
		modified := false
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					//return setting.Value
					revision = setting.Value
					if len(setting.Value) > 7 {
						revision = setting.Value[:7] // 最初の7文字にする
					}
				}
				if setting.Key == "vcs.modified" {
					modified = setting.Value == "true"
				}
			}
		}
		if modified {
			revision = "develop+" + revision
		}
		return revision
	}() // }}}
)

// init関数: ログの初期化。
func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func ShowHelp(post string) {
	buf := new(bytes.Buffer)
	parser.WriteHelp(buf)
	help := buf.String()
	help = strings.ReplaceAll(help, "display this help and exit", "ヘルプを出力する。")
	help = strings.ReplaceAll(help, "display version and exit", "バージョンを出力する。")

	lines := ""
	if !true {
		helps := strings.Split(help, "\n")
		for i, h := range helps {
			if strings.Contains(h, "Usage:") {
				usg := ""
				// `Usage: create_sw_iso [--root ROOT] [--install-media-dir-name NAME] [--year YEAR] [--sp SP] [--log LOGFILE] [--skip-get] [--skip-apply] [--skip-iso] [--skip-all] [--debug] [--verbose] [--code EXPORT-DIR]`
				ss := strings.Split(h, " [")
				p := GetFileNameWithoutExt(os.Args[0])
				for _, sss := range ss {
					if !strings.Contains(sss, "]") {
						usg += fmt.Sprintf("%v", sss)
					} else {
						sss = strings.ReplaceAll(sss, "[", "")
						sss = strings.ReplaceAll(sss, "]", "")
						k := strings.Split(sss, " ")
						if len(k) == 1 {
							usg += fmt.Sprintf("\n%v%v", strings.Repeat(" ", len(p)), sss)
						} else {
							usg += fmt.Sprintf("\n%v%-15s %s", strings.Repeat(" ", len(p)), k[0], k[1])
						}
					}
				}
				lines += "\n" + usg
			} else {
				if i != 0 {
					lines += "\n"
				}
				lines += h
			}
		}
	} else {
		lines = help
	}

	output := lines
	if len(post) != 0 {
		output += fmt.Sprintln(post)
	}
	fmt.Printf("%v\n", output)
	os.Exit(1)
}

func GetVersion() string {
	if len(revision) == 0 {
		// go installでビルドされた場合、gitの情報がなくなる。その場合v0.0.0.のように末尾に.がついてしまうのを避ける。
		return fmt.Sprintf("%v version %v\n", GetFileNameWithoutExt(os.Args[0]), version)
	} else {
		return fmt.Sprintf("%v version %v.%v\n", GetFileNameWithoutExt(os.Args[0]), version, revision)
	}
}

func ShowVersion() {
	fmt.Printf("%s", GetVersion())
	os.Exit(0)
}

func GetCurrentDir() string {
	ret, err := os.Getwd()
	if err != nil {
		panic(errors.Errorf("%v", err))
	}
	return filepath.ToSlash(ret)
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return !os.IsNotExist(err)
}

func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func WriteFile(filePath string, data []byte) {
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		panic(errors.Errorf("WriteFile: %v", err))
	}
	if _, err := f.Write(data); err != nil {
		panic(errors.Errorf("WriteFile: %v", err))
	}
}

func Abs(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	return p
}
