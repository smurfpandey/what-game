package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	ps "github.com/keybase/go-ps"
	"syscall"
	"time"
	"unsafe"
)

var (
	mod                          = windows.NewLazyDLL("user32.dll")
	procGetWindowText            = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength      = mod.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = mod.NewProc("GetWindowThreadProcessId")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

func GetWindowProcess(hwnd HWND) int {
	procId := 0;
	procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&procId)))

	return procId
}

func main() {
	for {
		if hwnd := getWindow("GetForegroundWindow"); hwnd != 0 {
			text := GetWindowText(HWND(hwnd))
			procId := GetWindowProcess(HWND(hwnd))
			
			yoProcess, _ := ps.FindProcess(procId)
			processExecName := yoProcess.Executable()
			processPath, _ := yoProcess.Path()

			fmt.Println("window :", text, "# procId:", procId, "# Process? ", processExecName, "# Path? ", processPath)
		}

		time.Sleep(2 * time.Second)
	}
}
