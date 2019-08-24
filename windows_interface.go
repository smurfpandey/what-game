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
	procGetGUIThreadInfo         = mod.NewProc("GetGUIThreadInfo")
)

const (
	UWP_HOST_APP = "ApplicationFrameHost.exe"
)

type (
	ATOM            uint16
	BOOL            int32
	COLORREF        uint32
	DWM_FRAME_COUNT uint64
	DWORD           uint32
	HACCEL          HANDLE
	HANDLE          uintptr
	HBITMAP         HANDLE
	HBRUSH          HANDLE
	HCURSOR         HANDLE
	HDC             HANDLE
	HDROP           HANDLE
	HDWP            HANDLE
	HENHMETAFILE    HANDLE
	HFONT           HANDLE
	HGDIOBJ         HANDLE
	HGLOBAL         HANDLE
	HGLRC           HANDLE
	HHOOK           HANDLE
	HICON           HANDLE
	HIMAGELIST      HANDLE
	HINSTANCE       HANDLE
	HKEY            HANDLE
	HKL             HANDLE
	HMENU           HANDLE
	HMODULE         HANDLE
	HMONITOR        HANDLE
	HPEN            HANDLE
	HRESULT         int32
	HRGN            HANDLE
	HRSRC           HANDLE
	HTHUMBNAIL      HANDLE
	HWND            HANDLE
	LPARAM          uintptr
	LPCVOID         unsafe.Pointer
	LRESULT         uintptr
	PVOID           unsafe.Pointer
	QPC_TIME        uint64
	ULONG_PTR       uintptr
	WPARAM          uintptr
	UINT            uint32
	LONG            int32
	SHORT           int16
	SIZE_T          ULONG_PTR
)

type RECT struct {
	Left   LONG
	Top    LONG
	Right  LONG
	Bottom LONG
}

type LPGUITHREADINFO struct {
	CbSize        DWORD
	Flags         DWORD
	HwndActive    HWND
	HwndFocus     HWND
	HwndCapture   HWND
	HwndMenuOwner HWND
	HwndMoveSize  HWND
	HwndCaret     HWND
	RcCaret       RECT
}

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

func GetGUIThreadInfo(idThread DWORD) *LPGUITHREADINFO {
	lpgui := new(LPGUITHREADINFO)
	lpgui.CbSize = DWORD(unsafe.Sizeof(*lpgui))
	procGetGUIThreadInfo.Call(uintptr(idThread), uintptr(unsafe.Pointer(lpgui)))

	return lpgui
}

func GetUWPAppProcess() int {
	gti := GetGUIThreadInfo(0)
	return GetWindowProcess(gti.HwndFocus)
}
