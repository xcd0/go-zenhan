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
	//if args.Verbose {
	//	log.Printf("%v\n", output)
	//}
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

func ShowReadme() {
	data, err := Asset("resources/src/readme.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: readme.md の取得に失敗しました: %v\n", err)
		return
	}
	fmt.Printf("%v", string(data))
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
		//panic(errors.Errorf("%v", err))
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
	//n, err := f.Write(data)

	if _, err := f.Write(data); err != nil {
		panic(errors.Errorf("WriteFile: %v", err))
	}
	//if args.Verbose {
	//	log.Printf("WriteFile: write %v bytes.", n)
	//}
}

func Abs(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	return p
}

// exportSourceCode は埋め込まれたソースコードを指定されたディレクトリに出力します
func exportSourceCode(outputDir string) error {
	fmt.Fprintf(os.Stderr, "ソースコードを %s に出力します\n", outputDir)
	{
		// ディレクトリを作成
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("出力ディレクトリの作成に失敗: %v", err)
		}
		//if err := os.MkdirAll(filepath.Join(outputDir, "src"), 0755); err != nil {
		//	return fmt.Errorf("出力ディレクトリの作成に失敗: %v", err)
		//}
	}

	// 埋め込みリソースディレクトリからソースファイルを取得
	srcDir := "resources/src"
	srcFiles := []string{}
	{
		var err error
		srcFiles, err = AssetDir(srcDir)
		if err != nil {
			return fmt.Errorf("埋め込みソースコードの取得に失敗: %v", err)
		}
		log.Printf("srcFiles: %v", srcFiles)
	}

	// 各ファイルを出力
	for _, filename := range srcFiles {
		path := filepath.ToSlash(filepath.Join(srcDir, filename))
		data, err := Asset(path)
		if err != nil {
			log.Printf("path: %v", path)
			fmt.Fprintf(os.Stderr, "警告: %s の取得に失敗しました: %v\n", filename, err)
			continue
		}
		outPath := filepath.Join(outputDir, filename)
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "警告: %s の書き込みに失敗しました: %v\n", outPath, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "  %s を出力しました\n", filename)
	}

	fmt.Fprintf(os.Stderr, "ソースコードの出力が完了しました\n")

	// 埋め込んでいるバイナリを出力。
	//exportEmbededBinary(outputDir)

	return nil
}

/*
// ここから下は使ってない。

var (
	tempDir          string
	toolPaths        = make(map[string]string) // ツール名からパスへのマッピング
	initialized      = false
	binary_debug_log = false
)

func SetDebugLogFlag(debug bool) {
	binary_debug_log = debug
}

func init() {
	if err := SetupTool("busybox", "resources/busybox64u.exe"); err != nil {
		fmt.Printf("エラー: Busyboxの初期化に失敗: %v\n", err)
		os.Exit(1)
	}
	if err := SetupTool("7z", "resources/7z.exe"); err != nil {
		fmt.Printf("エラー: 7zの初期化に失敗: %v\n", err)
		os.Exit(1)
	}
	if err := SetupTool("nkf", "resources/nkf.exe"); err != nil {
		fmt.Printf("エラー: nkfの初期化に失敗: %v\n", err)
		os.Exit(1)
	}
	if err := SetupTool("mkisofs", "resources/mkisofs.exe"); err != nil {
		fmt.Printf("エラー: mkisofsの初期化に失敗: %v\n", err)
		os.Exit(1)
	}
	if err := SetupTool("cygwin1.dll", "resources/cygwin1.dll"); err != nil {
		fmt.Printf("エラー: cygwin1.dllの初期化に失敗: %v\n", err)
		os.Exit(1)
	}
	setupCleanupHandler()
}

func example_tools() {
	// main関数の例
	// 起動時に各種ツールを展開
	if err := SetupTool("busybox", "resources/busybox64u.exe"); err != nil {
		log.Printf("エラー: Busyboxの初期化に失敗: %v\n", err)
		os.Exit(1)
	}

	// 別のツールも追加できます
	// err = SetupTool("other-tool", "resources/other-tool.exe")

	// ここで通常のアプリケーションロジックを実行
	// defer CleanupTools() は不要 - シグナルハンドラが処理します

	// busyboxコマンドの例
	if IsToolAvailable("busybox") {
		output, err := ExecuteTool("busybox", "ls", "-la")
		// または以下のように互換性関数も使用可能
		// output, err := ExecuteBusybox("ls", "-la")
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
		} else {
			fmt.Println("コマンド出力:")
			fmt.Println(output)
		}
	}

	// アプリケーションのメインロジックをここに記述
	// ... その他の処理

	// 明示的にクリーンアップする場合（通常は必要ありません）
	// CleanupTools()
}

func exportEmbededBinary(outputDir string) error {
	// 埋め込んだbinaryを出力する。
	{
		srcFiles := []string{
			"resources/7z.exe",
			"resources/bash.exe",
			"resources/busybox64u.exe",
			"resources/cygwin1.dll",
			"resources/jq.exe",
			"resources/mkisofs.exe",
			"resources/nkf.exe",
			"resources/upx.exe",
		}

		outPath := filepath.Join(outputDir, "src", "resources")
		if err := os.MkdirAll(outPath, 0750); err != nil {
			fmt.Fprintf(os.Stderr, "警告: %s の作成に失敗しました: %v\n", outPath, err)
		}
		for _, filename := range srcFiles {
			data, err := Asset(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "警告: %s の取得に失敗しました: %v\n", filename, err)
				continue
			}
			outPath := filepath.Join(outputDir, "src", filename)
			if err := os.WriteFile(outPath, data, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "警告: %s の書き込みに失敗しました: %v\n", outPath, err)
				continue
			}
			fmt.Fprintf(os.Stderr, "  %s を出力しました\n", filename)
		}

		move := [][]string{
			[]string{"src/build.bat", "build.bat"},
			[]string{"src/readme.md", "readme.md"},
			[]string{"src/run.bat", "run.bat"},
			[]string{"src/.gitignore", ".gitignore"},
		}
		for _, m := range move {
			src := filepath.Join(outputDir, m[0])
			dst := filepath.Join(outputDir, m[1])
			err := os.Rename(src, dst)
			if err != nil {
				fmt.Fprintf(os.Stderr, "警告: %s の %s への移動に失敗しました: %v\n", src, dst, err)
			} else {
				fmt.Fprintf(os.Stderr, "  %s を出力しました\n", dst)
			}
		}
	}
}

// SetupTool は指定されたツールを一時ディレクトリに展開します
func SetupTool(toolName, assetPath string) error {

	// 一時ディレクトリの初期化
	if tempDir == "" {
		var err error
		tempDir, err = os.MkdirTemp("", "app_tools")
		if err != nil {
			return fmt.Errorf("一時ディレクトリの作成に失敗: %w", err)
		}

		// シグナルハンドラを設定してプログラム終了時にクリーンアップ
		setupCleanupHandler()
		initialized = true
	}

	// 埋め込まれたツールデータを取得
	toolData, err := Asset(assetPath)
	if err != nil {
		return fmt.Errorf("埋め込みファイル %s の取得に失敗: %w", assetPath, err)
	}

	// 一時ファイルとして書き出す
	toolPath := filepath.Join(tempDir, filepath.Base(assetPath))
	if err := os.WriteFile(toolPath, toolData, 0755); err != nil {
		return fmt.Errorf("%s の書き出しに失敗: %w", toolName, err)
	}

	// ツールパスを記録
	toolPaths[toolName] = toolPath

	if binary_debug_log {
		log.Printf("%s が展開されました: %s\n", toolName, toolPath)
	}
	return nil
}

// GetToolPath は指定されたツールのパスを取得します
func GetToolPath(toolName string) (string, bool) {
	path, exists := toolPaths[toolName]
	return path, exists
}

// IsToolAvailable は指定されたツールが利用可能かどうかを確認します
func IsToolAvailable(toolName string) bool {
	_, exists := toolPaths[toolName]
	return exists
}

// ExecuteTool は指定されたツールでコマンドを実行します
func ExecuteTool(toolName string, args ...string) (string, error) {
	toolPath, exists := toolPaths[toolName]
	if !exists {
		return "", fmt.Errorf("%s が初期化されていません。先にSetupToolを呼び出してください", toolName)
	}

	if binary_debug_log {
		log.Printf("exec.Command(%v, %#q)", toolPath, args)
	}
	log.Printf("exec.Command(%v, %#q)", toolPath, args)

	log_cmd := fmt.Sprintf("%#v", filepath.ToSlash(toolPath))
	for _, a := range args {
		log_cmd += fmt.Sprintf(" %v", a)
	}
	log.Printf("%v", log_cmd)

	cmd := exec.Command(toolPath, args...)
	output, err := cmd.CombinedOutput()
	output_str := string(output)
	log.Printf("stdout: %v", output_str)
	log.Printf("stderr: %v", err)

	if err != nil {
		return "", fmt.Errorf("%s コマンド実行エラー: %w", toolName, err)
	}
	return output_str, nil
}

func ExecuteToolOutputNkf(toolName string, cmd_args ...string) (string, error) {
	toolPath, exists := toolPaths[toolName]
	if !exists {
		return "", fmt.Errorf("%s が初期化されていません。先にSetupToolを呼び出してください", toolName)
	}

	if args.Verbose {
		log.Printf("exec.Command(%v, %#q)", toolPath, cmd_args)
	}

	log_cmd := fmt.Sprintf("%#v", filepath.ToSlash(toolPath))
	for _, a := range cmd_args {
		log_cmd += fmt.Sprintf(" %v", a)
	}
	if args.Verbose {
		log.Printf("%v", log_cmd)
	}

	cmd := exec.Command(toolPath, cmd_args...)
	output, err := cmd.CombinedOutput()
	output_str := string(output)

	if args.Verbose {
		log.Printf("out: %#v", output_str)
	}
	if err != nil {
		log.Printf("err: %#v", err.Error())
	}

	if out_u8 := nkf(output); len(out_u8) > 0 {
		filePath := filepath.ToSlash(filepath.Join(GetCurrentDir(), "cmd.stdout.txt"))
		WriteFile(filePath, []byte(out_u8))
		if args.Verbose {
			log.Printf("filePath.stdout: %v", filePath)
			if args.Verbose {
				s := string(out_u8)
				ss := strings.Split(s, "\n")
				for _, sss := range ss {
					log.Printf("stdout: %v", sss)
				}
			}
		}
		os.RemoveAll(filePath)
	}
	if err != nil {
		if err_u8 := nkf([]byte(err.Error())); len(err_u8) > 0 {
			filePathErr := filepath.ToSlash(filepath.Join(GetCurrentDir(), "cmd.stderr.txt"))
			WriteFile(filePathErr, []byte(err_u8))
			log.Printf("filePath.stderr: %v", filePathErr)
			log.Printf("stderr: %v", string(err_u8))
		}
		return "", fmt.Errorf("%s コマンド実行エラー: %w", toolName, err)
	}
	return output_str, nil
}

func mktemp() string {
	in, err := os.CreateTemp("", "temp.*")
	if err != nil {
		panic(errors.Errorf("%v", err))
	}
	return filepath.ToSlash(in.Name())
}

func nkf(b []byte) []byte {
	if args.Verbose {
		log.Printf("nkf: input length: %v", len(b))
	}

	in, err := os.CreateTemp("", "nkf.*.in.txt")
	if err != nil {
		panic(errors.Errorf("%v", err))
	}
	defer os.Remove(in.Name())
	out, err := os.CreateTemp("", "nkf.*.out.txt")
	if err != nil {
		panic(errors.Errorf("%v", err))
	}
	defer os.Remove(out.Name())

	WriteFile(filepath.ToSlash(in.Name()), b)

	callNkf := func(in, out string) ([]byte, error) {
		toolPath, exists := toolPaths["nkf"]
		if !exists {
			return nil, fmt.Errorf("%s が初期化されていません。先にSetupToolを呼び出してください", "nkf")
		}
		toolPath = filepath.ToSlash(toolPath)
		if args.Verbose {
			log.Printf("exec.Command(%#v, %#v, %#v, %#v)", toolPath, "-wO", in, out)
		}
		cmd := exec.Command(toolPath, "-wO", in, out)
		output, err := cmd.CombinedOutput()
		return output, err
	}

	// 文字コード変換
	if _, err := callNkf(filepath.ToSlash(in.Name()), filepath.ToSlash(out.Name())); err != nil {
		panic(errors.Errorf("nkf コマンド実行エラー: %v", err))
	}
	bout, err := os.ReadFile(out.Name())
	if err != nil {
		panic(errors.Errorf("nkf コマンド実行エラー: %v", err))
	}
	return bout
}

// CleanupTools は一時ディレクトリとツールを削除します
func CleanupTools() {
	if tempDir != "" {
		os.RemoveAll(tempDir)
		fmt.Println("一時ディレクトリを削除しました")
		tempDir = ""
		toolPaths = make(map[string]string)
		initialized = false
	}
}

// グローバル変数として関数ポインタのスライスを定義する。
// 各要素は引数なし、戻り値なしの関数を指す。
var funcSlice []func()

// ! addFunction はグローバルな関数スライスに新たな関数を追加する関数である。
func addFunction(f func()) {
	// スライスに関数を追加する。
	funcSlice = append(funcSlice, f)
}

// ! callAllFunctions はグローバルな関数スライスに保持されている全ての関数を呼び出す関数である。
func callAllFunctions() {
	// スライス内の各関数を順に呼び出す。
	for _, f := range funcSlice {
		f()
	}
}

// setupCleanupHandler は終了シグナルを捕捉してクリーンアップを実行するハンドラを設定します
func setupCleanupHandler() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n終了シグナルを受信しました")
		CleanupTools()
		callAllFunctions()
		os.Exit(0)
	}()
}
*/
