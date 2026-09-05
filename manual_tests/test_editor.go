package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const grokAPIURL = "https://api.groq.com/openai/v1/chat/completions"
const grokAPIKey = "your_groq_api_key" // Замени на свой API-ключ

// Editor структура для хранения состояния редактора
type Editor struct {
	screen    tcell.Screen
	lines     []string // Строки текста
	cursorX   int      // Позиция курсора X
	cursorY   int      // Позиция курсора Y
	filename  string   // Имя файла
	dirty     bool     // Флаг изменений
}

// NewEditor создаёт новый редактор
func NewEditor(filename string) (*Editor, error) {
	lines := []string{""} // Начальная пустая строка
	if filename != "" {
		data, err := ioutil.ReadFile(filename)
		if err == nil && len(data) > 0 {
			lines = strings.Split(string(data), "\n")
			if len(lines) == 0 || lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
		}
	}
	return &Editor{
		lines:    lines,
		filename: filename,
	}, nil
}

// InitScreen инициализирует экран
func (e *Editor) InitScreen() error {
	var err error
	e.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := e.screen.Init(); err != nil {
		return err
	}
	e.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	e.screen.Clear()
	return nil
}

// Draw отображает текст на экране
func (e *Editor) Draw() {
	e.screen.Clear()
	for y, line := range e.lines {
		for x, ch := range line {
			e.screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}
	e.screen.ShowCursor(e.cursorX, e.cursorY)
	e.screen.Show()
}

// Save сохраняет файл
func (e *Editor) Save() error {
	data := strings.Join(e.lines, "\n")
	if err := ioutil.WriteFile(e.filename, []byte(data), 0644); err != nil {
		return err
	}
	e.dirty = false
	return nil
}

// CallGrokAI вызывает Grok API для автодополнения
func (e *Editor) CallGrokAI() (string, error) {
	currentLine := e.lines[e.cursorY]
	payload := map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{"role": "user", "content": fmt.Sprintf("Complete the following text: %s", currentLine)},
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", grokAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+grokAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	return result.Choices[0].Message.Content, nil
}

// HandleInput обрабатывает ввод
func (e *Editor) HandleInput() bool {
	for {
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlQ: // Выход
				return false
			case tcell.KeyCtrlS: // Сохранение
				if err := e.Save(); err != nil {
					fmt.Printf("Error saving file: %v\n", err)
				}
			case tcell.KeyCtrlI: // Вызов ИИ
				completion, err := e.CallGrokAI()
				if err != nil {
					fmt.Printf("Error calling AI: %v\n", err)
					continue
				}
				e.lines[e.cursorY] += completion
				e.cursorX = len(e.lines[e.cursorY])
				e.dirty = true
			case tcell.KeyEnter:
				// Новая строка
				currentLine := e.lines[e.cursorY]
				afterCursor := currentLine[e.cursorX:]
				e.lines[e.cursorY] = currentLine[:e.cursorX]
				e.lines = append(e.lines[:e.cursorY+1], append([]string{afterCursor}, e.lines[e.cursorY+1:]...)...)
				e.cursorY++
				e.cursorX = 0
				e.dirty = true
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				// Удаление символа
				if e.cursorX > 0 {
					currentLine := e.lines[e.cursorY]
					e.lines[e.cursorY] = currentLine[:e.cursorX-1] + currentLine[e.cursorX:]
					e.cursorX--
					e.dirty = true
				} else if e.cursorY > 0 {
					// Объединить с предыдущей строкой
					currentLine := e.lines[e.cursorY]
					e.cursorY--
					e.cursorX = len(e.lines[e.cursorY])
					e.lines[e.cursorY] += currentLine
					e.lines = append(e.lines[:e.cursorY+1], e.lines[e.cursorY+2:]...)
					e.dirty = true
				}
			case tcell.KeyRune:
				// Ввод символа
				ch := ev.Rune()
				currentLine := e.lines[e.cursorY]
				e.lines[e.cursorY] = currentLine[:e.cursorX] + string(ch) + currentLine[e.cursorX:]
				e.cursorX++
				e.dirty = true
			case tcell.KeyLeft:
				if e.cursorX > 0 {
					e.cursorX--
				}
			case tcell.KeyRight:
				if e.cursorX < len(e.lines[e.cursorY]) {
					e.cursorX++
				}
			case tcell.KeyUp:
				if e.cursorY > 0 {
					e.cursorY--
					if e.cursorX > len(e.lines[e.cursorY]) {
						e.cursorX = len(e.lines[e.cursorY])
					}
				}
			case tcell.KeyDown:
				if e.cursorY < len(e.lines)-1 {
					e.cursorY++
					if e.cursorX > len(e.lines[e.cursorY]) {
						e.cursorX = len(e.lines[e.cursorY])
					}
				}
			}
		case *tcell.EventResize:
			e.screen.Sync()
		}
		e.Draw()
	}
}

func main() {
	filename := "example.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	editor, err := NewEditor(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := editor.InitScreen(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer editor.screen.Fini()

	editor.Draw()
	if !editor.HandleInput() {
		if editor.dirty {
			fmt.Println("Unsaved changes. Save before exiting? (y/n)")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" {
				if err := editor.Save(); err != nil {
					fmt.Printf("Error saving file: %v\n", err)
				}
			}
		}
	}
}