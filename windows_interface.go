package main

import (
	"golang.org/x/sys/windows"
	"syscall"
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

func getWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := getWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func GetWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

func GetWindowProcess(hwnd HWND) int {
	procId := 0
	procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&procId)))

	return procId
}
