package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	iprompt    = flag.String("iprompt", "> ", "Input prompt line prefix")
	oprompt    = flag.String("oprompt", "> ", "Output prompt line prefix")
	delay      = flag.Int("delay", 80, "Typing delay between chars in ms")
	delayvar   = flag.Int("delay-var", 10, "Typing delay variation")
	spdelay    = flag.Int("spdelay", 10, "Extra typing delay before a space, to stager manual imput a little more naturally")
	spdelayvar = flag.Int("spdelay-var", 80, "Extra typing delay before a space variation")

	cmddelay    = flag.Int("cmddelay", 200, "Command output delay between lines in ms")
	cmddelayvar = flag.Int("cmddelay-var", 100, "Command output delay variation")

	title  = flag.String("title", "Screencast", "Title of cast")
	cmd    = flag.String("cmd", "asciinemagen", "Command used for cast")
	height = flag.Int("height", 62, "Height of cast")
	width  = flag.Int("width", 135, "Width of cast")
)

// {
//   "title": "Screencast",
//   "version": 1,
//   "command": null,
//   "stdout": [
//     [
//       0.021967,
//       "Hola señor"
//     ],
//   ],
//   "height": 62,
//   "env": {
//     "SHELL": "/bin/bash",
//     "TERM": "xterm-256color"
//   },
//   "duration": 31.891245,
//   "width": 135
// }
type asciinemaV1 struct {
	Title    string  `json:"title,omitempty"`
	Version  int     `json:"version,omitempty"`
	Command  string  `json:"command,omitempty"`
	Height   int     `json:"height,omitempty"`
	Width    int     `json:"width,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Env      env     `json:"env,omitempty"`
	Stdout   []input `json:"stdout"`
}

//   {
//     "SHELL": "/bin/bash",
//     "TERM": "xterm-256color"
//   },
type env struct {
	Shell string `json:"SHELL,omitempty"`
	Term  string `json:"TERM,omitempty"`
}

//     [
//       0.021967,
//       "Hola señor"
//     ]
type input struct {
	delay  time.Duration
	output string
}

func (i input) MarshalJSON() ([]byte, error) {
	out := strconv.AppendFloat([]byte{'['}, i.delay.Seconds(), 'f', 6, 64)
	out = append(out, ',')
	res, err := json.Marshal(i.output)
	if err != nil {
		return nil, err
	}
	out = append(out, res...)
	out = append(out, ']')
	return out, nil
}

func main() {
	flag.Parse()

	out := &asciinemaV1{
		Title:   *title,
		Version: 1,
		Command: *cmd,
		Height:  *height,
		Width:   *width,
		Env: env{
			Shell: "/bin/bash",
			Term:  "xterm-256color",
		},
	}

	reader := bufio.NewReader(os.Stdin)
	var line, in []byte
	var isPrefix bool
	var err error
OUTER_LOOP:
	for {
		for {
			in, isPrefix, err = reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break OUTER_LOOP
				}
				log.Fatalf("stdin reading failed: %s", err)
			}
			line = append(line, in...)
			if !isPrefix {
				break
			}
		}

		l := string(line)
		if strings.HasPrefix(l, *iprompt) {
			d := vardelay(*cmddelay, *cmddelayvar)
			out.Stdout = append(out.Stdout, input{
				delay:  time.Duration(d) * time.Millisecond,
				output: *oprompt,
			})
			d = vardelay(*cmddelay, *cmddelayvar)
			extraDelay := d * 2
			for _, char := range l[2:] {
				d := vardelay(*delay, *delayvar)
				// fmt.Fprintln(os.Stderr, d, extraDelay)
				d += extraDelay
				if char == ' ' {
					spd := vardelay(*spdelay, *spdelayvar)
					d += spd
				}
				out.Stdout = append(out.Stdout, input{
					delay:  time.Duration(d) * time.Millisecond,
					output: string(char),
				})
				extraDelay = 0
			}
			d = vardelay(*cmddelay, *cmddelayvar)
			out.Stdout = append(out.Stdout, input{
				delay:  time.Duration(d) * time.Millisecond,
				output: "\n\r",
			})
		} else {
			d := vardelay(*cmddelay, *cmddelayvar)
			out.Stdout = append(out.Stdout, input{
				delay:  time.Duration(d) * time.Millisecond,
				output: l + "\n\r",
			})
		}

		line = line[:0]
	}

	duration := time.Duration(0)
	for _, input := range out.Stdout {
		duration += input.delay
	}
	out.Duration = duration.Seconds()

	// fmt.Println(out)
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func vardelay(dd, ddv int) int {
	d := dd
	v := rand.Int31n(int32(ddv))
	if rand.Int31n(2) == 1 {
		d += int(v)
	} else {
		d -= int(v)
	}
	if d < 0 {
		return 0
	}
	return d
}
