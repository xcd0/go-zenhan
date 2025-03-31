package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
)

// コマンドライン引数構造体。
type Args struct {
	// Status はIMEのオープン状態を設定するための値である。
	// 指定されなかった場合は現在の状態を取得する。
	Status     int    `arg:"positional" help:"設定したいIMEのON/OFF状態を 1 or 0 で指定する。指定がないとき現在の状態を出力する。" placeholder:"FLAG"`
	ExportCode string `arg:"--code"     help:"ソースコード出力先ディレクトリパス。バイナリに埋め込まれたソースコードを指定ディレクトリに出力する。" placeholder:"EXPORT-DIR"`
	Debug      bool   `arg:"-d,--debug" help:"デバッグログ出力を有効にする。"`
}

func (Args) Version() string {
	return GetVersion()
}

// グローバル変数。
var (
	args   Args
	parser *arg.Parser // ShowHelp() で使う
)

// コマンドライン引数の解析。
func ParseArgs() {

	// 初期値設定。
	args = Args{Status: -1}

	var err error
	parser, err = arg.NewParser(arg.Config{Program: GetFileNameWithoutExt(os.Args[0]), IgnoreEnv: false}, &args)
	if parser == nil {
		log.Printf("parser: %v", parser)
		log.Printf("err   : %v", err)
	}
	if err != nil {
		ShowHelp(fmt.Sprintf("%v", errors.Errorf("%v", err)))
		os.Exit(1)
	}
	err = parser.Parse(os.Args[1:])

	if err != nil {
		if err.Error() == "help requested by user" {
			ShowHelp("")
			os.Exit(1)
		} else if err.Error() == "version requested by user" {
			ShowVersion()
			os.Exit(0)
		} else {
			ShowHelp("")
			panic(errors.Errorf("%v", err))
		}
	}

	if !args.Debug {
		log.SetOutput(ioutil.Discard)
	}

	// 即時終了する処理
	if len(args.ExportCode) > 0 {
		if p, err := filepath.Abs(args.ExportCode); err != nil {
			panic(fmt.Errorf("ソースコードの出力に失敗しました: %v\n", err))
			os.Exit(1)
		} else {
			args.ExportCode = filepath.ToSlash(p)
		}
		if err := exportSourceCode(args.ExportCode); err != nil {
			panic(fmt.Errorf("ソースコードの出力に失敗しました: %v\n", err))
			os.Exit(1)
		}
		os.Exit(0)
	}
}
