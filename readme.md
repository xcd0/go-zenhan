# go-zenhan

https://github.com/iuchim/zenhan

のgolang実装。

## モチベ

単に出来心で...

## 利点

- 簡単にビルドできる。
- わかりやすいヘルプ付き。
	```sh
	$ ./go-zenhan.exe -h
	go-zenhan version 0.0.1.develop+0539993

	Usage: go-zenhan [--code EXPORT-DIR] [--debug] [FLAG]

	Positional arguments:
	  FLAG                   設定したいIMEのON/OFF状態を 1 or 0 で指定する。指定がないとき現在の状態を出力する。

	Options:
	  --code EXPORT-DIR      ソースコード出力先ディレクトリパス。バイナリに埋め込まれたソースコードを指定ディレクトリに出力する。
	  --debug, -d            デバッグログ出力を有効にする。
	  --help, -h             ヘルプを出力する。
	  --version              バージョンを出力する。
	```

## 欠点

- おっきい (upxで圧縮しても800KB超)

## 備考

元の実装は **40KB** しかないので何ら持ち運びに困らない。


