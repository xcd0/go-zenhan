# go-zenhan

https://github.com/iuchim/zenhan

のgolang実装。

## モチベ

単に出来心で...

## 利点

- 簡単にビルドできる。
- インストールが簡単。
	```sh
	go install github.com/xcd0/go-zenhan@latest
	```
	- wslの場合。
		```sh
		unset GOBIN; GOOS=windows go install github.com/xcd0/go-zenhan@latest
		```
		下記に配置される。
		```
		$GOPATH/bin/windows_amd64/go-zenhan.exe
		```
- わかりやすいヘルプ付き。
	```sh
	$ ./go-zenhan.exe -h
	go-zenhan version 0.0.*

	Usage: go-zenhan [--code EXPORT-DIR] [--debug] [FLAG]

	Positional arguments:
	  FLAG                   設定したいIMEのON/OFF状態を 1 or 0 で指定する。指定がないとき現在の状態を出力する。

	Options:
	  --debug, -d            デバッグログ出力を有効にする。
	  --help, -h             ヘルプを出力する。
	  --version              バージョンを出力する。
	```

## 欠点

- おっきい (upxで圧縮しても800KB超)

## 備考

元の実装は **40KB** しかないので何ら持ち運びに困らない。


