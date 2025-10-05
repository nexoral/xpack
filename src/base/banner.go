package base

import (
	"fmt"
)

// const Version is kept so external scripts can update it via sed
const Version = "1.2.4-stable"

// PrintBanner prints a simple welcome banner. Version may be empty.
func PrintBanner(version string) {
	if version == "" {
		version = Version
	}
	fmt.Println("====================================")
	fmt.Println("          Welcome to xpack          ")
	fmt.Printf("           version: %s\n", version)
	fmt.Println("  A minimal packaging helper (deb,tgz)")
	fmt.Println("====================================")
}
