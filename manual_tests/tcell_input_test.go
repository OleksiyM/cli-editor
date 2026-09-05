package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Set UTF-8 environment
	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("LC_CTYPE", "en_US.UTF-8")

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

	defer screen.Fini()

	// Set style
	screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))
	screen.Clear()

	// Display instructions
	instructions := []string{
		"Type some Cyrillic text (привет) and see what happens",
		"Press ESC to exit",
	}

	for y, line := range instructions {
		for x, r := range line {
			screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}

	screen.Show()

	// Event loop
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return
			}

			// Display the key event info
			screen.Clear()
			
			// Show event details
			info := fmt.Sprintf("Key: %v, Rune: %c (U+%04X), Mod: %v", 
				ev.Key(), ev.Rune(), ev.Rune(), ev.Modifiers())
			
			for x, r := range info {
				if x < 100 { // Limit display width
					screen.SetContent(x, 0, r, nil, tcell.StyleDefault)
				}
			}
			
			// Show instructions again
			for y, line := range instructions {
				for x, r := range line {
					screen.SetContent(x, y+2, r, nil, tcell.StyleDefault)
				}
			}
			
			screen.Show()
		}
	}
}