package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"xpack/src/base"
	"xpack/src/packager"

	"golang.org/x/term"
)

// var VERSION is updated by Scripts/versionController.sh
var VERSION = "1.2.5-stable"

// readLineRaw reads interactive input from terminal in raw mode, supporting Tab completion
func readLineRaw(promptText, defaultVal string) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// fallback to simple buffered read
		var line string
		if defaultVal != "" {
			fmt.Printf("%s [%s]: ", promptText, defaultVal)
		} else {
			fmt.Printf("%s: ", promptText)
		}
		_, err := fmt.Scanln(&line)
		if err != nil {
			// user might have pressed enter without input
			return defaultVal, nil
		}
		return strings.TrimSpace(line), nil
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	fmt.Print(promptText)
	if defaultVal != "" {
		fmt.Printf(" [%s]", defaultVal)
	}
	fmt.Print(": ")

	var buf []rune
	for {
		var b = make([]byte, 1)
		n, err := os.Stdin.Read(b)
		if err != nil || n == 0 {
			return "", err
		}
		c := b[0]
		if c == 3 { // Ctrl-C
			return "", fmt.Errorf("interrupted")
		}
		if c == '\r' || c == '\n' {
			// enter -- move to a new line so subsequent output doesn't continue the prompt
			fmt.Print("\n")
			if len(buf) == 0 {
				return defaultVal, nil
			}
			return string(buf), nil
		}
		if c == 9 { // Tab
			prefix := string(buf)
			matches, _ := filepath.Glob(prefix + "*")
			if len(matches) == 1 {
				// complete
				buf = []rune(matches[0])
				// re-render current buffer
				fmt.Print("\r")
				fmt.Print(strings.Repeat(" ", 80))
				fmt.Print("\r")
				if defaultVal != "" {
					fmt.Printf("%s [%s]: %s", promptText, defaultVal, string(buf))
				} else {
					fmt.Printf("%s: %s", promptText, string(buf))
				}
				continue
			}
			if len(matches) > 1 {
				// list options
				fmt.Print("\nMatches:\n")
				for i, m := range matches {
					fmt.Printf("  %d) %s\n", i+1, m)
				}
				fmt.Print("Select number or press enter to continue: ")
				var sel int
				_, err := fmt.Scanf("%d", &sel)
				if err == nil && sel > 0 && sel <= len(matches) {
					buf = []rune(matches[sel-1])
					fmt.Print("\r")
					fmt.Print(strings.Repeat(" ", 80))
					fmt.Print("\r")
					if defaultVal != "" {
						fmt.Printf("%s [%s]: %s", promptText, defaultVal, string(buf))
					} else {
						fmt.Printf("%s: %s", promptText, string(buf))
					}
				} else {
					// reprint prompt and buffer
					fmt.Print("\r")
					fmt.Print(strings.Repeat(" ", 80))
					fmt.Print("\r")
					if defaultVal != "" {
						fmt.Printf("%s [%s]: %s", promptText, defaultVal, string(buf))
					} else {
						fmt.Printf("%s: %s", promptText, string(buf))
					}
				}
				continue
			}
			// no matches, ignore
			continue
		}
		if c == 127 || c == 8 { // backspace
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\r")
				fmt.Print(strings.Repeat(" ", 80))
				fmt.Print("\r")
				if defaultVal != "" {
					fmt.Printf("%s [%s]: %s", promptText, defaultVal, string(buf))
				} else {
					fmt.Printf("%s: %s", promptText, string(buf))
				}
			}
			continue
		}
		// printable
		buf = append(buf, rune(c))
		fmt.Printf("%c", c)
	}
}

// selectOption displays an interactive list and allows arrow-key selection
func selectOption(options []string, defaultIndex int) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// fallback
		if defaultIndex >= 0 && defaultIndex < len(options) {
			return options[defaultIndex], nil
		}
		return options[0], nil
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	// create local menu with an Exit option appended
	menu := make([]string, 0, len(options)+1)
	menu = append(menu, options...)
	menu = append(menu, "Exit")

	selected := defaultIndex
	if selected < 0 || selected >= len(options) {
		selected = 0
	}

	clearMenu := func(n int) {
		// move cursor up to the start of menu
		fmt.Printf("\033[%dA", n)
		// clear each line in-place without emitting newlines
		for i := 0; i < n; i++ {
			// clear the entire line
			fmt.Print("\033[2K")
			// move cursor down one line to the next line
			fmt.Print("\033[1B")
		}
		// move cursor back up to the original position (just after the cleared block)
		fmt.Printf("\033[%dA", n)
		// ensure we're at the start of the line
		fmt.Print("\r")
	}

	// initial render
	for {
		// print the menu (including Exit)
		for i, opt := range menu {
			if i == selected {
				// clear line then print highlighted
				fmt.Printf("\033[2K\r\033[7m> %s\033[0m\n", opt)
			} else {
				fmt.Printf("\033[2K\r  %s\n", opt)
			}
		}

		// read a key
		var b = make([]byte, 1)
		if _, err := os.Stdin.Read(b); err != nil {
			return "", err
		}
		if b[0] == 0x1b { // ESC
			// likely an arrow key sequence: read two more
			var seq = make([]byte, 2)
			if _, err := os.Stdin.Read(seq); err != nil {
				return "", err
			}
			if seq[0] == '[' {
				switch seq[1] {
				case 'A': // up
					if selected > 0 {
						selected--
					} else {
						selected = len(menu) - 1
					}
				case 'B': // down
					selected = (selected + 1) % len(menu)
				}
			}
			// move cursor up to redraw in place
			fmt.Printf("\033[%dA", len(menu))
			continue
		}
		if b[0] == '\r' || b[0] == '\n' {
			// if Exit selected, clear the menu and return a sentinel error
			if selected == len(menu)-1 {
				clearMenu(len(menu))
				return "", fmt.Errorf("exit")
			}
			clearMenu(len(menu))
			return menu[selected], nil
		}
		// numeric quick-select (1..9)
		if b[0] >= '1' && b[0] <= '9' {
			idx := int(b[0] - '1')
			if idx >= 0 && idx < len(menu) {
				if idx == len(menu)-1 {
					clearMenu(len(menu))
					return "", fmt.Errorf("exit")
				}
				clearMenu(len(menu))
				return menu[idx], nil
			}
		}
		// otherwise ignore and continue
		// move cursor up to redraw
		fmt.Printf("\033[%dA", len(menu))
	}
}

func main() {
	// flags
	inputFlag := flag.String("i", "", "relative path to input binary")
	archFlag := flag.String("arch", "", "target architecture (amd64, arm64, i386)")
	verFlag := flag.String("v", "", "version string for the package (e.g. 1.1.1)")
	flag.Parse()

	// show banner
	bannerVersion := VERSION
	// prefer VERSION file if present (over var)
	if data, err := os.ReadFile(filepath.Join("..", "..", "VERSION")); err == nil {
		bannerVersion = strings.TrimSpace(string(data))
	}
	base.PrintBanner(bannerVersion)

	// interactive prompts if flags not provided
	inputPath := *inputFlag
	if inputPath == "" {
		if v, err := readLineRaw("Path to input binary (relative allowed)", "./bin/xpack"); err == nil {
			inputPath = v
		} else {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	}
	// normalize relative path (do not resolve to absolute; preserve relative semantics)
	inputPath = filepath.Clean(inputPath)

	defaultArch := runtime.GOARCH
	if defaultArch == "amd64" {
		defaultArch = "amd64"
	}
	arch := *archFlag
	if arch == "" {
		// present interactive selector
		opts := []string{"amd64", "arm64", "i386"}
		// find default index
		defIdx := 0
		for i, o := range opts {
			if o == defaultArch {
				defIdx = i
				break
			}
		}
		if sel, err := selectOption(opts, defIdx); err == nil {
			arch = sel
		} else {
			// if user selected Exit, quit
			if err.Error() == "exit" {
				fmt.Println("Exiting")
				os.Exit(0)
			}
			arch = defaultArch
		}
	}

	// Resolve binary file path relative to current working directory
	// Do not convert to absolute so relative paths remain meaningful when used in CI scripts
	if _, err := os.Stat(inputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: input binary not found at %s\n", inputPath)
		os.Exit(1)
	}

	// infer app name from binary filename
	var appName string
	defaultCandidate := filepath.Base(inputPath)
	 v, _ := readLineRaw("Application name for this package", filepath.Base(inputPath))
		if v != "" {
					appName = v
				} else {
					appName = defaultCandidate
				}
	// remove extension
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))

	// determine package version:
	// - if -v provided, use it
	// - otherwise, if running interactively prompt the user (default from VERSION file -> var VERSION -> bannerVersion)
	// - if not interactive (CI), fall back to VERSION file -> var VERSION -> bannerVersion
	versionStr := "0.0.0"
	if *verFlag != "" {
		versionStr = *verFlag
	} else {
		// pick a sensible default candidate from VERSION file, build var, or banner
		defaultCandidate := bannerVersion
		if data, err := os.ReadFile(filepath.Join("..", "..", "VERSION")); err == nil {
			defaultCandidate = strings.TrimSpace(string(data))
		} else if VERSION != "" {
			defaultCandidate = VERSION
		}

		// if stdin is a terminal, prompt the user so they can change or accept the default
		if term.IsTerminal(int(os.Stdin.Fd())) {
			if v, err := readLineRaw("Version for this build", defaultCandidate); err == nil {
				if v != "" {
					versionStr = v
				} else {
					versionStr = defaultCandidate
				}
			} else {
				// on prompt error, fallback to defaultCandidate
				versionStr = defaultCandidate
			}
		} else {
			// non-interactive: use the default candidate
			versionStr = defaultCandidate
		}
	}

	// build packages
	outDir := "dist"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// ensure there is a blank line separating user input from program output
	fmt.Print("\n")
	fmt.Printf("Packaging %s (version %s) for arch %s...\n", appName, versionStr, arch)
	// call packager to build .deb and tar.gz
	if err := packager.BuildAll(inputPath, appName, versionStr, arch, outDir); err != nil {
		fmt.Fprintf(os.Stderr, "Packaging failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Packaging complete. Artifacts are in %s\n", outDir)
}
