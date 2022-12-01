package gotree

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
	"math/rand"
)

const (
	AttrReset     = "\x1b[0m"
	AttrBold      = "\x1b[1m"
	AttrItalics   = "\x1b[3m"
	AttrUnderline = "\x1b[4m"
	AttrBlink     = "\x1b[5m"

	AttrBlack   = "\x1b[30m"
	AttrRed     = "\x1b[31m"
	AttrGreen   = "\x1b[32m"
	AttrYellow  = "\x1b[33m"
	AttrBlue    = "\x1b[34m"
	AttrMagenta = "\x1b[35m"
	AttrCyan    = "\x1b[36m"
	AttrWhite   = "\x1b[37m"

	AttrGrey          = "\x1b[90m"
	AttrBrightRed     = "\x1b[91m"
	AttrBrightGreen   = "\x1b[92m"
	AttrBrightYellow  = "\x1b[93m"
	AttrBrightBlue    = "\x1b[94m"
	AttrBrightMagenta = "\x1b[95m"
	AttrBrightCyan    = "\x1b[96m"
	AttrBrightWhite   = "\x1b[97m"
)

// Config
var (
	Width  = 19
	Height = 10

	Offset = 1

	LeafChar   = '*'
	LeafAttr   = AttrGreen
	LeafChance = 4

	LightBulbChar  = 'o'
	LightBulbAttr  = AttrBold
	LightBulbAttrs = []string{AttrBrightRed,
	                          AttrBrightGreen,
	                          AttrBrightYellow,
	                          AttrBrightBlue,
	                          AttrBrightMagenta,
	                          AttrBrightCyan,
	                          AttrBrightWhite}

	// Speed in ticks, a tick takes 60 miliseconds
	LightsSpeed = 6

	TrunkStr  = "wWw"
	TrunkAttr = AttrBrightRed
	TrunkLen  = 2

	StarChar = '&'
	StarAttr = AttrBold + AttrBrightYellow
)

func StartTree(blocking bool) error {
	if Width <= 0 {
		return fmt.Errorf("Width <= 0")
	} else if Height <= 0 {
		return fmt.Errorf("Height <= 0")
	} else if TrunkLen < 0 {
		return fmt.Errorf("Trunk length <= 0")
	} else if LeafChance < 2 {
		return fmt.Errorf("Leaf chance < 2")
	} else if Offset < 0 {
		return fmt.Errorf("Offset < 0")
	}

	rand.Seed(time.Now().UnixNano())

	if blocking {
		// Init terminal for non-blocking input
		if err := termInit(); err != nil {
			return err
		}
		defer termRestore()

		hideCursor()
		defer showCursor()

		tick := 0

		// Input
		in := make([]byte, 1)

		for {
			if tick % LightsSpeed == 0 {
				if tick > 0 {
					// Move cursor to the start of the tree to re-render it
					moveCursorUp(TrunkLen + Height)
				}

				renderTree()
			}

			os.Stdin.Read(in)
			if in[0] != 0 { // When no key is pressed, Read returns 0
				break
			}

			time.Sleep(60 * time.Millisecond)
			tick ++
		}

	} else {
		// Make the light bulbs blink so it still looks animated
		tmp := LightBulbAttr
		LightBulbAttr = AttrBlink + LightBulbAttr

		renderTree()

		LightBulbAttr = tmp
	}

	return nil
}

func renderTree() {
	// Render the leaves
	for i := 0; i < Height; i ++ {
		if i == 0 {
			fmt.Printf("%v%v\n", strings.Repeat(" ", Width / 2 + Offset),
			           StarAttr + string(StarChar) + AttrReset)
		} else {
			// Calculate the current leaf row width
			width := int(float64(i + 1) / float64(Height) * float64(Width))
			if width <= 0 {
				width = 1
			}

			renderLeafRow(width)
		}
	}

	// Render the trunk
	offset := Width / 2 - TrunkLen / 2 + Offset
	out    := TrunkAttr + TrunkStr + AttrReset

	for i := 0; i < TrunkLen; i ++ {
		fmt.Printf("%v%v\n", strings.Repeat(" ", offset), out)
	}
}

func renderLeafRow(width int) {
	offset := Width / 2 - width / 2 + Offset

	fmt.Print(strings.Repeat(" ", offset))

	for i := 0; i < width; i ++ {
		if rand.Intn(LeafChance) == 1 { // Render a light bulb
			attr := LightBulbAttr + LightBulbAttrs[rand.Intn(len(LightBulbAttrs))]

			fmt.Print(attr + string(LightBulbChar) + AttrReset)
		} else { // Render a leaf
			fmt.Print(LeafAttr + string(LeafChar) + AttrReset)
		}
	}

	fmt.Println()
}

func moveCursorUp(by int) {
	fmt.Printf("\x1b[%vA", by)
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
}

var mode string

func termInit() error {
	// Save the previous terminal attributes
	bytes, err := exec.Command("stty", "-F", "/dev/tty", "-g").Output()
	if err != nil {
		return err
	}

	mode = string(bytes[:len(bytes) - 1])

	// Ignore CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c

		showCursor()
		termRestore()
		os.Exit(0)
	}()

	// Non-blocking input, no input echo and no timeout
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "-echo", "-echok", "-icanon",
	             "min", "0", "time", "0").Run()

	return nil
}

func termRestore() {
	exec.Command("stty", "-F", "/dev/tty", mode).Run()
}
