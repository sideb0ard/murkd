package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"
)

func cleanuprrr(line string) string {
	line = RemoveMarkUp(line)
	line = RemoveMarkdownImages(line)
	line = HeaderToUpperBold(line)
	line = MarkdownHyperlinkToBold(line)
	return line
}

func RemoveMarkUp(str string) string {
	rehtml := regexp.MustCompile(`<.*>`)
	return rehtml.ReplaceAllString(str, "")
}

func RemoveMarkdownImages(str string) string {
	rehtml := regexp.MustCompile(`^\[!\[.*`)
	return rehtml.ReplaceAllString(str, "")
}

func RemoveRegex(rex *regexp.Regexp, str string) string {
	return rex.ReplaceAllString(str, "")
}

func HeaderToUpperBold(str string) string {
	// change headers into upper case bold
	reheaders := regexp.MustCompile(`^## (.*)`)
	str = reheaders.ReplaceAllStringFunc(str, func(m string) string {
		parts := reheaders.FindStringSubmatch(m)
		//return Bold(strings.ToUpper(parts[1]))
		return Colorize("white", strings.ToUpper(parts[1]))
	})

	return str
}

func MarkdownHyperlinkToBold(str string) string {
	// change MD hyperlinks into bold text
	relynx := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	str = relynx.ReplaceAllStringFunc(str, func(m string) string {
		parts := relynx.FindStringSubmatch(m)
		return Bold(parts[1])
	})
	return str
}
func Colorize(clr string, str string) string {
	//lime := ansi.ColorCode("green+h:black")
	clr = ansi.ColorCode(clr)
	reset := ansi.ColorCode("reset")
	return clr + str + reset
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

//////////////////////////////////////////

func main() {

	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	markdownfile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Err, mate:", err.Error)
		os.Exit(1)
	}
	defer markdownfile.Close()

	reempty := regexp.MustCompile(`^\s*$`)
	previousLineEmpty := false

	rebackticks := regexp.MustCompile("^```")
	backticksOn := false

	scanner := bufio.NewScanner(markdownfile)
	for scanner.Scan() {

		cleanline := cleanuprrr(scanner.Text())

		// multiline code block
		if rebackticks.MatchString(cleanline) {
			if backticksOn != true {
				backticksOn = true
				cleanline = RemoveRegex(rebackticks, cleanline)
			} else {
				backticksOn = false
				cleanline = RemoveRegex(rebackticks, cleanline)
			}
		}
		if backticksOn == true && len(cleanline) != 0 {
			cleanline = Colorize("green", cleanline)
		}

		// reduce multiple consecutive lines to a single one
		if reempty.MatchString(cleanline) && previousLineEmpty == true {
			continue
		} else if reempty.MatchString(cleanline) {
			previousLineEmpty = true
		} else {
			previousLineEmpty = false
		}

		fmt.Println(cleanline)
	}
}

//////////////////////////////////////////

func usage() {
	fmt.Println("Yo! murk me up, brah - feed me markdown files..", os.Args[0], "<markdownfile>")
}
