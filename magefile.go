// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	exeFileName = "tkihou.exe"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Run

// Clean clean up after yourself
func Clean() {
	fmt.Println("Clean...")
	os.Remove(exeFileName)
}

// Build build
func Build() error {
	mg.Deps(Clean)
	fmt.Println("Build...")
	return sh.RunV("go", "build", "-o", exeFileName, ".")
}

// Run execute app
func Run() error {
	mg.Deps(Build)
	fmt.Println("Run...")
	return sh.RunV("./" + exeFileName)
}

// Test execute test
func Test() error {
	fmt.Println("Test...")
	return sh.RunV("go", "test")
}

// Generate generate code
func Generate() error {
	fmt.Println("Generate...")
	return nil
}
