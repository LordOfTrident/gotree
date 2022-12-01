package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/LordOfTrident/gotree/pkg/gotree"
)

var (
	blocking = flag.Bool("blocking", true, "Wait for a key press")

	width  = flag.Int("width",  19, "Width of the tree")
	height = flag.Int("height", 10, "Height of the tree")

	offset = flag.Int("offset", 1, "Tree offset from the left")

	trunkStr = flag.String("trunkStr", "wWw", "Trunk rendering string")
	trunkLen = flag.Int("trunkLen",    2,     "Length of the trunk")

	leafChance = flag.Int("leafChance", 3, "The chance of a leaf appearing")

	lightsSpeed = flag.Int("lightsSpeed", 6,
	                       "Speed of the lights flickering in ticks, 1 tick = 60ms")
)

func init() {
	flag.Parse()

	gotree.Width  = *width
	gotree.Height = *height

	gotree.Offset = *offset

	gotree.TrunkStr = *trunkStr
	gotree.TrunkLen = *trunkLen

	gotree.LeafChance = *leafChance

	gotree.LightsSpeed = *lightsSpeed
}

func main() {
	err := gotree.StartTree(*blocking)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())

		os.Exit(1)
	}
}
