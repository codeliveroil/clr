package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/codeliveroil/go-util/errors"
	"github.com/codeliveroil/niceflags"
)

const defaultClr = "\033[m"

var colorMap = map[string]string{
	"black":       "0",
	"red":         "1",
	"green":       "2",
	"yellow":      "3",
	"blue":        "4",
	"pink":        "5",
	"cyan":        "6",
	"lightgray":   "7",
	"darkgray":    "8",
	"lightred":    "9",
	"lightgreen":  "10",
	"lightyellow": "11",
	"lightblue":   "12",
	"lightpink":   "13",
	"lightcyan":   "14",
	"white":       "15",
}

type filter struct {
	exp           *regexp.Regexp
	highlightLine bool
	code          colorCode
}

func parse(args []string) []*filter {
	if len(args) == 0 {
		fmt.Println("No color-rules specified.")
		os.Exit(1)
	}

	var filters []*filter

	for _, arg := range args {
		if arg == "" { // can happen if there is a new line in the rules file
			continue
		}
		tokens := strings.Split(arg, "~")
		if len(tokens) < 2 {
			fmt.Println("Invalid syntax, see '-help':", arg)
			os.Exit(1)
		}

		// Parse regex/string
		exp, err := regexp.Compile(tokens[0])
		errors.Exit(err, "Cannot parse regex: %q", tokens[0])

		f := &filter{
			exp: exp,
		}

		// Parse foreground and background colours
		makeClr := func(clr string, bgFgOp string) string {
			if clr == "default" {
				return "\033[" + bgFgOp + "9m"
			}

			c, ok := colorMap[clr]
			if !ok {
				cint, err := strconv.Atoi(clr)
				if err != nil || (cint < 0 || cint > 255) {
					fmt.Println("Invalid color, see '-help':", clr)
					os.Exit(1)
				}
				c = clr
			}
			return "\033[" + bgFgOp + "8;5;" + c + "m"
		}

		clrTokens := strings.Split(tokens[1], ",")
		fgClr := makeClr(clrTokens[0], "3")

		var bgClr string
		if len(clrTokens) == 2 {
			bgClr = makeClr(clrTokens[1], "4")
		}

		f.code = colorCode(fgClr + bgClr)
		filters = append(filters, f)

		// Parse additional options
		if len(tokens) < 3 {
			continue
		}

		opts := strings.Split(tokens[2], ",")

		for _, o := range opts {
			switch o {
			case "line":
				f.highlightLine = true
			}
		}
	}

	return filters
}

func color(line string, filters []*filter) string {
	type cell struct {
		startColor []colorCode
		endColor   []colorCode
	}

	cells := make(map[int]*cell)

	// Assign start and end colors to each character position in the line
	//
	setCell := func(start, end int, color colorCode) {
		get := func(index int) *cell {
			c, ok := cells[index]
			if !ok {
				c = &cell{}
				cells[index] = c
			}
			return c
		}

		c := get(start)
		c.startColor = append(c.startColor, color)

		c = get(end)
		c.endColor = append(c.endColor, color)
	}

	for _, f := range filters {
		if f.highlightLine {
			if found := f.exp.MatchString(line); found {
				setCell(0, len(line), f.code)
			}
			continue
		}

		locs := f.exp.FindAllStringIndex(line, -1)
		for _, loc := range locs {
			setCell(loc[0], loc[1], f.code)
		}
	}

	// Generate the colorized string
	//
	s := stack{}
	var buf bytes.Buffer
	for i, c := range line {
		cell, ok := cells[i]
		if !ok {
			buf.WriteString(string(c))
			continue
		}

		// First end any highlights that need to be ended.
		for _, ec := range cell.endColor {

			// Pop the ended color, this may not always be on the top of the stack.
			// For example:
			//   echo "Hello there, Kobe!" | clr "Hello  there, Kobe!~red" "there~green" "there, Kobe!~yellow"
			// will yield:
			//   [Red]Hello [Green,Yellow]there[\Green], Kobe![\Red,\Yellow]
			// When "there" is being printed, the stack is going to be [Yellow,Green,Red]. If we simply pop(), we'll end the yellow
			// highlight. This is wrong. We must end the green highlight. That's why we need to pop() until we see green, then push
			// back all those that were popped() and were not green.

			// pop all start colors till we meet the end color
			tmp := stack{}
			for e, ok := s.pop(); ok; e, ok = s.pop() {
				if e == ec {
					break
				}
				tmp.push(e)
			}
			// push the popped but unended start colors back
			for e, ok := tmp.pop(); ok; e, ok = tmp.pop() {
				s.push(e)
			}

			// reset background color if any. Otherwise, we'll color the remaining text in the next foreground color,
			// but maintain the previous background color.
			buf.WriteString(defaultClr)

			// If we have nested colors: e.g. [Red]Hello [Green]there[\Green], Kobe![\Red]
			// Then, when "there" is done being highlighted green, we need to switch to red again and not the
			// default color (\033[m).
			endClr, ok := s.peek()
			if !ok {
				endClr = defaultClr
			}
			buf.WriteString(string(endClr))
		}

		// Next, start coloring any highlights that need to be colored.
		for _, sc := range cell.startColor {
			buf.WriteString(string(sc))
			s.push(sc)
		}

		// Last, append the character
		buf.WriteString(string(c))

	}

	buf.WriteString(defaultClr)

	return buf.String()
}

func printSwatch() {
	flipColors := []int{0, 16, 17, 18, 19, 232, 233, 234, 235, 236}

	m := make(map[int]bool)
	for _, c := range flipColors {
		m[c] = true
	}

	for i := 0; i <= 255; i++ {
		fg := 16
		if _, ok := m[i]; ok {
			fg = 248
		}
		fmt.Printf("\033[38;5;%dm\033[48;5;%dm% 4d\033[m ", fg, i, i)
		if (i+1)%16 == 0 {
			fmt.Printf("\n")
		}
	}
}

func main() {
	flags := niceflags.NewFlags("clr",
		"clr - Output colorizer",
		"Colorizes piped output. For instance, you can pipe logs to clr and set a color rule for info and error logs.",
		`[color-rule] [color-rule] ...

color-rule: is a tilde separated rule of the following format:
  search-string~colors[~extra-options] 
 
search-string: can be any string or a regular expression of RE2 syntax (https://github.com/google/re2/wiki/Syntax).

colors: is of the format:
  foreground-color[,background-color]
The foreground-color and background-color can be one of black, red, green, yellow, blue, pink, cyan, lightgray, darkgray, lightred, lightgreen, lightyellow, lightblue, lightpink, lightcyan, white or default.
Alternatively, you can use a color from 0 to 255 from the color chart rendered with the -swatch option.

extra-options: is optional and is a comma-separated list of the following options:
  line: highlight the entire line.

Examples:
  echo 'Hello, World 2019' | clr "Hello~red" "[[:digit:]]~blue,yellow"
  echo 'Hello, World 2019' | clr "World~default,pink~line"
  echo 'Hello, World 2019' | clr "2019~11,52"
`,
		"help",
		false,
	)

	rules := flags.String("rules", "", "Read rules from file instead of the command line. Separate each color-rule by a new line.")
	swatch := flags.Bool("swatch", false, "Display the colors for the 0-255 color range.")

	args := os.Args[1:]
	err := flags.Parse(args)
	errors.Exit(err, "Cannot parse")

	flags.Help()

	if *swatch {
		printSwatch()
		return
	}

	if f := *rules; f != "" {
		b, err := ioutil.ReadFile(f)
		errors.Exit(err, "Cannot read file")
		args = strings.Split(string(b), "\n")
	}

	filters := parse(args)

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			errors.Exit(err, "Cannot read from pipe")
		}

		line = line[:len(line)-1] // remove the newline char

		fmt.Println(color(line, filters))
	}
}
