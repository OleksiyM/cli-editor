package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Set UTF-8 environment
	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LANG", "en_US.UTF-8")

	// Create screen
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	// Enable UTF-8 fallback
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)

	if err := screen.Init(); err != nil {
		panic(err)
	}

	// Set style
	screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))

	// Clear screen
	screen.Clear()

	// Draw test text with Cyrillic characters
	text := "hell привет"
	x, y := 0, 0
	for _, r := range text {
		screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		x++
	}

	// Show instructions
	instructions := "Press any key to exit"
	x, y = 0, 2
	for _, r := range instructions {
		screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		x++
	}

	// Refresh screen
	screen.Show()

	// Wait for event
	screen.PollEvent()

	// Clean up
	screen.Fini()
}