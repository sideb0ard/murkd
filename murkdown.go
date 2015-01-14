package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"
)

const (
	CODECOLOR = "white"
)

func cleanuprrr(line string) string {
	line = RemoveMarkUp(line)
	line = RemoveListMarker(line)
	line = RemoveMarkdownImages(line)
	line = HeaderToUpperBold(line)
	line = MarkdownHyperlinkToBold(line)
	line = CodeSpanReplace(line)
	line = Emphasis(line)
	return line
}

func RemoveListMarker(str string) string {
	relist := regexp.MustCompile(`^\s?\*\s`)
	return relist.ReplaceAllString(str, "  ")
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

func Emphasis(str string) string {
	// this is messy but only way i could think of without
	// regex lookbehinds - otherwise this was replacing _underscore_
	// parts of variable names, rather than on word boundary
	// ( asterix don't count as a word character )
	reemph := regexp.MustCompile(`([*_]+)([a-zA-Z\s'\."]+)([*_]+)`)
	return reemph.ReplaceAllStringFunc(str, func(m string) string {
		re := regexp.MustCompile(`([*_]+)([a-zA-Z\s'\."]+)([*_]+)(?:[^a-zA-Z])`)
		parts := re.FindStringSubmatch(m)
		if len(parts) > 1 {
			return Colorize(CODECOLOR, parts[2])
		} else {
			return m
		}
	})
}

func CodeSpanReplace(str string) string {
	respan := regexp.MustCompile("(`{1})([^`].*?)(`{1})")
	return respan.ReplaceAllStringFunc(str, func(m string) string {
		parts := respan.FindStringSubmatch(m)
		return Colorize(CODECOLOR, parts[2])
	})
}

func HeaderToUpperBold(str string) string {
	reheaders := regexp.MustCompile(`^#+ (.*)`)
	return reheaders.ReplaceAllStringFunc(str, func(m string) string {
		parts := reheaders.FindStringSubmatch(m)
		return Colorize("white", strings.ToUpper(parts[1]))
	})
}

func MarkdownHyperlinkToBold(str string) string {
	relynx := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	str = relynx.ReplaceAllStringFunc(str, func(m string) string {
		parts := relynx.FindStringSubmatch(m)
		return Bold(parts[1])
	})
	return str
}
func Colorize(clr string, str string) string {
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

	recommentline := regexp.MustCompile("^//")
	reheaderline := regexp.MustCompile(`^[=-]+$`)

	// temp to hold previous line so we can catch headers
	bufferline := "ignorefirsttime"
	refirst := regexp.MustCompile("ignorefirsttime")

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
		if backticksOn == true && len(cleanline) != 0 && !recommentline.MatchString(cleanline) {
			cleanline = Colorize(CODECOLOR, cleanline)
		}

		// reduce multiple consecutive lines to a single one
		if reempty.MatchString(cleanline) && previousLineEmpty == true {
			continue
		} else if reempty.MatchString(cleanline) {
			previousLineEmpty = true
		} else {
			previousLineEmpty = false
		}

		// buffer current in previous so we can check for header
		if refirst.MatchString(bufferline) {
			bufferline = cleanline
			continue
		}
		if reheaderline.MatchString(cleanline) {
			bufferline = Colorize("white", strings.ToUpper(bufferline))
		}

		fmt.Println(bufferline)
		bufferline = cleanline

	}
}

//////////////////////////////////////////

func usage() {
	fmt.Println("Yo! murk me up, brah - feed me markdown files..", os.Args[0], "<markdownfile>")
}
