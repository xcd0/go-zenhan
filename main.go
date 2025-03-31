package main

import (
	"fmt"
	"log"

	// コマンドライン引数解析用ライブラリ。
	"github.com/pkg/errors"
	"golang.org/x/sys/windows" // Windows API呼び出し用ライブラリ。
)

func main() {

	// 引数の解析。
	ParseArgs()

	// 前面ウィンドウのハンドルを取得する。
	hwnd := getForegroundWindow()
	if hwnd == 0 {
		log.Printf("Error: 前面ウィンドウが取得できませんでした。")
		panic(errors.Errorf("Error: 前面ウィンドウが取得できませんでした。"))
	}

	// 前面ウィンドウに対するデフォルトIMEウィンドウのハンドルを取得する。
	ime := getDefaultIMEWnd(hwnd)
	if ime == 0 {
		log.Printf("Error: IMEウィンドウが取得できませんでした。")
		panic(errors.Errorf("Error: IMEウィンドウが取得できませんでした。"))
	}

	var stat uintptr
	if args.Status < 0 {
		// 引数が指定されなかった場合、IMEのオープン状態を取得する。
		stat = sendMessage(ime, WM_IME_CONTROL, IMC_GETOPENSTATUS, 0)
	} else {
		// 引数が指定された場合、IMEのオープン状態を設定する。
		if args.Status != 0 {
			args.Status = 1
		}
		stat = uintptr(args.Status)
		sendMessage(ime, WM_IME_CONTROL, IMC_SETOPENSTATUS, stat)
	}

	// 結果の状態を標準出力に出力する。
	log.Printf("IME status: %d", stat)
	fmt.Printf("%d\n", stat)
}

const (
	WM_IME_CONTROL    = 0x283 // WM_IME_CONTROL メッセージの定義である。
	IMC_GETOPENSTATUS = 5     // IMC_GETOPENSTATUS の定義である。
	IMC_SETOPENSTATUS = 6     // IMC_SETOPENSTATUS の定義である。
)

// getForegroundWindow は現在の前面ウィンドウのハンドルを取得する関数である。
func getForegroundWindow() windows.Handle {
	// user32.dll から GetForegroundWindow 関数を取得する。
	user32 := windows.NewLazySystemDLL("user32.dll")
	procGetForegroundWindow := user32.NewProc("GetForegroundWindow")
	ret, _, _ := procGetForegroundWindow.Call()
	return windows.Handle(ret)
}

// getDefaultIMEWnd は指定したウィンドウのデフォルトIMEウィンドウのハンドルを取得する関数である。
func getDefaultIMEWnd(hwnd windows.Handle) windows.Handle {
	// imm32.dll から ImmGetDefaultIMEWnd 関数を取得する。
	imm32 := windows.NewLazySystemDLL("imm32.dll")
	procImmGetDefaultIMEWnd := imm32.NewProc("ImmGetDefaultIMEWnd")
	ret, _, _ := procImmGetDefaultIMEWnd.Call(uintptr(hwnd))
	return windows.Handle(ret)
}

// sendMessage は指定したウィンドウにメッセージを送信する関数である。
func sendMessage(hwnd windows.Handle, msg, wParam, lParam uintptr) uintptr {
	// user32.dll から SendMessageW 関数を取得する。
	user32 := windows.NewLazySystemDLL("user32.dll")
	procSendMessage := user32.NewProc("SendMessageW")
	ret, _, _ := procSendMessage.Call(uintptr(hwnd), msg, wParam, lParam)
	return ret
}
