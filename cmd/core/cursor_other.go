//go:build !windows

package main

func getCursorPos() (int, int, bool) {
	return 0, 0, false
}
