package version

import "fmt"

var Version = "unknown"

func Print() {
	fmt.Println("Version: " + Version)
}
