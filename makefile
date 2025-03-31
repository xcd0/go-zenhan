BIN           := ./go-zenhan
REVISION      = $(shell git rev-parse --short HEAD)
VERSION       := 0.0.1
FLAGS_VERSION := -X main.version=$(VERSION)
FLAG          := -a -tags netgo -trimpath -ldflags='-s -w -extldflags="-static" $(FLAGS_VERSION) -buildid='
RESOURCE_DIR  := resources
BINDATA_FILE  := bindata.go
SOURCE_FILES  := go.mod go.sum *.go makefile *.md .gitignore
SET_GOOS      := GOOS=windows

.PHONY: all build clean release update-binary

all:
	cat ./makefile | grep '^[^ ]*:$$'

build:
	make update-binary
	$(SET_GOOS) go build -o $(BIN).exe

release:
	make clean-bindata
	make build
	GOOS=windows go build $(FLAG) -o $(BIN).exe
	#$(RESOURCE_DIR)/upx --lzma $(BIN).exe
	upx --lzma $(BIN).exe
	cp -rf $(BIN).exe ..
	echo Success!

# 埋め込むデータの更新
update-binary:
	@if ! which go >/dev/null 2>&1 ; then echo "goが見つかりません。wingetでGoのコンパイラをインストールします。"; winget install GoLang.Go; fi
	@if ! [ -e "$(RESOURCE_DIR)" ]; then mkdir -p "$(RESOURCE_DIR)"; fi
#	@if ! [ -e "$(RESOURCE_DIR)/src" ]; then mkdir -p "$(RESOURCE_DIR)/src"; fi
#	@if ! [ -e "$(RESOURCE_DIR)/busybox64u.exe" ]; then make get-busybox; fi
#	@if ! [ -e "$(RESOURCE_DIR)/upx.exe" ]; then make get-upx; fi
#	@if ! [ -e "$(RESOURCE_DIR)/jq.exe" ]; then make get-jq; fi
#	@if ! [ -e "$(RESOURCE_DIR)/7z.exe" ]; then make get-7z; fi
#	@if ! [ -e "$(RESOURCE_DIR)/nkf.exe" ]; then make get-nkf; fi
#	@if ! [ -e "$(RESOURCE_DIR)/mkisofs.exe" ]; then make get-mkisofs; fi
#	@if ! [ -e "$(RESOURCE_DIR)/cygwin1.dll" ]; then make get-mkisofs; fi
	cp -rfp $(SOURCE_FILES) $(RESOURCE_DIR)/src/
	rm -rf $(RESOURCE_DIR)/src/bindata.go
	make gen-bindata

gen-bindata:
	if which go-bindata >/dev/null; then :; else go install github.com/go-bindata/go-bindata/...@latest ; fi
	go-bindata -o $(BINDATA_FILE) $(RESOURCE_DIR)/ $(RESOURCE_DIR)/src/

clean-bindata:
	rm -rf "$(BINDATA_FILE)"
clean-resource:
	rm -rf "$(RESOURCE_DIR)"
clean:
	make clean-bindata
	make clean-resource
	rm -rf "$(BIN).exe"

get-golang:
	# これは流石にwingetしていいよね...
	# windows以外は色々なので頑張って_(:3 」∠ )_
	winget install GoLang.Go
#get-upx:
#	# winget install upxでもよい。
#	# ここではビルド時にバイナリに埋め込むことを想定して配置する。
#	curl -L `curl -s https://api.github.com/repos/upx/upx/releases/latest | grep "browser_download_url" | grep "win64.zip" | cut -d"\"" -f4` -o upx.zip; unzip -jo upx.zip "upx*/upx.exe" -d .; mv upx.exe "$(RESOURCE_DIR)/" ; rm upx.zip
#get-busybox:
#	curl.exe -L "https://frippery.org/files/busybox/busybox64u.exe" -o "$(RESOURCE_DIR)/busybox64u.exe"
#	cp -rfp "$(RESOURCE_DIR)/busybox64u.exe" bash.exe || :
#	cp -rfp "$(RESOURCE_DIR)/busybox64u.exe" "$(RESOURCE_DIR)/bash.exe" || :
#get-jq:
#	# winget install jqlang.jqでもよい。
#	# ここではビルド時にバイナリに埋め込むことを想定して配置する。
#	# jqを使わず、jqの最新版をGitHubから取ってくるのは難しい。1.7.1を取ってきて取る。
#	curl.exe -s -L https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-windows-amd64.exe -o "./$(RESOURCE_DIR)/jq171.exe" \
#		&& curl -L "$$( \
#			curl.exe -s "https://api.github.com/repos/jqlang/jq/releases/latest" \
#			| ./$(RESOURCE_DIR)/jq171.exe -r '.assets[].browser_download_url | select(endswith("amd64.exe"))' \
#		)" -o $(RESOURCE_DIR)/jq.exe \
#		&& rm ./$(RESOURCE_DIR)/jq171.exe \
#		&& ./$(RESOURCE_DIR)/jq.exe --version
#get-7z:
#	curl -L "$$( \
#		curl.exe -s "https://api.github.com/repos/ip7z/7zip/releases/latest" \
#		| ./$(RESOURCE_DIR)/jq.exe -r '.assets[].browser_download_url | select(endswith("7zr.exe"))' \
#	)" -o $(RESOURCE_DIR)/7z.exe \
#	&& ./$(RESOURCE_DIR)/7z.exe | grep "7-Zip"
#get-nkf:
#	# winget にない。GitHubから最新版を取得する。jqを使用する。とはいえ最新といっても...
#	# ここではビルド時にバイナリに埋め込むことを想定して配置する。
#	JQ=$$(which jq); [ -n "$$JQ" ] && echo 1 || JQ="./$(RESOURCE_DIR)/jq.exe"; \
#		echo "$$JQ"; \
#		curl.exe -s "https://api.github.com/repos/kkato233/nkf/releases/latest" \
#		| "$$JQ" -r '.assets[].browser_download_url | select(endswith(".zip"))' \
#		| { \
#				read url; \
#				curl.exe -L "$$url" -o "nkf.zip" \
#				&& e=$$( unzip -l "nkf.zip" | grep "nkf.exe" | awk '{print $$4}' ); \
#				unzip -p "nkf.zip" "$$e" > "./$(RESOURCE_DIR)/nkf.exe" && rm "nkf.zip"; \
#			} \
#		&& [ -e "./$(RESOURCE_DIR)/nkf.exe" ] \
#			&& "./$(RESOURCE_DIR)/nkf.exe" "--version" | head -n1 | cut -c 0-35 \
#			|| { \
#				echo "Error: Download nkf failed !!!" >&2; \
#				exit 1; \
#			}
#get-mkisofs:
#	# 10MBくらいあるのでちょっと時間かかる。
#	
#	cd "./$(RESOURCE_DIR)" ; \
#	if [ -e "./mkisofs.exe" ] && "./mkisofs.exe" --version | head -n 1 | grep "mkisofs" -q; then \
#		"./mkisofs.exe" --version | head -n 1; \
#	else \
#		if [ ! -e "./cdrtfe.zip" ] || ! unzip -l "./cdrtfe.zip" >/dev/null 2>&1 ; then \
#			echo "cdrtfe.zip をダウンロードします。"; \
#			rm -f "./cdrtfe.zip"; \
#			curl -L http://sourceforge.net/projects/cdrtfe/files/cdrtfe/cdrtfe%201.5.9/cdrtfe-1.5.9.zip -o cdrtfe.zip; \
#		fi; \
#		echo "cdrtfe.zip からmkisofs.exeを展開します。"; \
#		echo "$$(unzip -l "./cdrtfe.zip" | grep -e "cygwin1.dll" -e "mkisofs.exe" | awk '{print $$4}')"; \
#		echo "$$(unzip -l "./cdrtfe.zip" | grep -e "cygwin1.dll" -e "mkisofs.exe" | awk '{print $$4}')" | xargs -I{} 7z e -aoa "./cdrtfe.zip" "{}" >/dev/null 2>&1; \
#		v=$$("./mkisofs.exe" --version | head -n 1); \
#		if echo "$$v" | grep "mkisofs" -q; then \
#			echo "$$v"; \
#			rm "./cdrtfe.zip"; \
#		else \
#			echo "Error: Download mkisofs failed !!!" >&2; \
#			exit 1; \
#		fi; \
#	fi

