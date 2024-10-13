package hangman

import (
	"syscall"
	"time"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procKeybdEvent       = user32.NewProc("keybd_event")
	VK_F11          byte = 0x7A
	KEYEVENTF_KEYUP      = 0x0002
)

func PressF11() {
	procKeybdEvent.Call(uintptr(VK_F11), 0, 0, 0)
	time.Sleep(100 * time.Millisecond)
	procKeybdEvent.Call(uintptr(VK_F11), 0, uintptr(KEYEVENTF_KEYUP), 0)
}
