package main

import "fmt"

func main() {
	NewGame()
}

// ══════════════════════════
//  UTILS
// ══════════════════════════

func debuglog(str string) {
	fmt.Println(str)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SendMessageToGui(header string, message string) {
	return
}
