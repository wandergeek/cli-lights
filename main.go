package main

import (
	// "bufio"
	"flag"
	"fmt"
	// "io/ioutil"
	"os"
)

var colors = map[string][]byte{
	"black":   {0, 0, 0},
	"red":     {1, 0, 0},
	"green":   {0, 1, 0},
	"blue":    {0, 0, 1},
	"cyan":    {0, 1, 1},
	"magenta": {1, 0, 1},
	"yellow":  {1, 1, 0},
	"white":   {1, 1, 1},
}

var led_types = map[string]int{
	"power":          0,
	"wired_internet": 1,
	"wireless":       2,
	"pairing":        3,
	"radio":          4,
}

var led_positions = [][]byte{
	{15, 13, 14},
	{12, 10, 11},
	{9, 1, 8},
	{2, 4, 3},
	{5, 7, 6},
}

// var leds = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var leds = []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

var p = fmt.Println

func main() {

	colorPtr := flag.String("color", "black", "give me a color mofo")
	flag.Parse()

	fmt.Print("You passed: ")
	fmt.Println(*colorPtr)
	setColor(1, *colorPtr, false)
	setLEDs()
}

func _setColor(position int, color []byte) {
	var indexes = led_positions[position]
	for i := 0; i < 3; i++ {
		leds[indexes[i]] = color[i]
	}
}

func setColor(position int, color string, flash bool) {
	_setColor(position, colors[color])

	p(flash)

}

func setLEDs() {

	writetofile("/sys/kernel/debug/omap_mux/lcd_data15", []byte("27"))
	writetofile("/sys/kernel/debug/omap_mux/lcd_data14", []byte("27"))
	writetofile("/sys/kernel/debug/omap_mux/uart0_ctsn", []byte("27"))
	writetofile("/sys/kernel/debug/omap_mux/mii1_col", []byte("27"))

	if _, err := os.Stat("/sys/class/gpio/gpio11/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", []byte("11"))
	}

	if _, err := os.Stat("/sys/class/gpio/gpio10/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", []byte("10"))
	}

	if _, err := os.Stat("/sys/class/gpio/gpio40/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", []byte("40"))
	}

	if _, err := os.Stat("/sys/class/gpio/gpio96/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", []byte("96"))
	}

	writetofile("/sys/class/gpio/gpio11/direction", []byte("low"))
	writetofile("/sys/class/gpio/gpio10/direction", []byte("low"))
	writetofile("/sys/class/gpio/gpio40/direction", []byte("low"))
	writetofile("/sys/class/gpio/gpio96/direction", []byte("low"))

	for i := 0; i < len(leds); i++ {
		writetofile("/sys/class/gpio/gpio40/value", leds[i:i+1])
		p(leds[i : i+1])
		writetofile("/sys/class/gpio/gpio96/value", []byte("1"))
		writetofile("/sys/class/gpio/gpio96/value", []byte("0"))
	}

	writetofile("/sys/class/gpio/gpio11/value", []byte("1"))
	writetofile("/sys/class/gpio/gpio11/value", []byte("0"))
}

func writetofile(fn string, val []byte) error {
	p("writing to" + fn)
	df, err := os.OpenFile(fn, os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(df, val)
	df.Close()
	return nil
}
