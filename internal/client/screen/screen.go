package screen

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/seandisero/celaeno/internal/shared"
	"golang.org/x/term"
)

const spaceRune = 1
const widthBuffer = 2

type Screen struct {
	Width  int
	input  chan shared.Message
	buffer string
}

func NewScreen(width int) *Screen {
	scrn := &Screen{
		Width: width,
		input: make(chan shared.Message),
	}

	return scrn
}

func (s *Screen) ClearInput(message string) error {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}
	totalLength := len(message)
	lines := int(math.Ceil(float64(totalLength) / float64(width)))

	for range lines {
		s.ClearMessageBox()
	}
	return nil
}

func (s *Screen) HandleMessage(message shared.Message) {
	s.input <- message
}

func (s *Screen) ClearMessageBox() {
	fmt.Print("\033[A")
	fmt.Print("\033[K")
}

func (s *Screen) MessageLoop(ctx context.Context) error {
	for {
		select {
		case msg := <-s.input:
			err := s.printToScreen(msg, msg.Incoming)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			fmt.Println("leaving message loop")
			return nil
		}
	}
}

func (s *Screen) CelaenoResponse(message string) {
	celMessage := shared.Message{
		Username: "celaeno",
		Message:  message,
		Incoming: true,
	}
	s.HandleMessage(celMessage)
}

func (s *Screen) printToScreen(message shared.Message, incoming bool) error {
	lines, err := s.formatMessage(message.Message)
	if err != nil {
		return fmt.Errorf("could not format message: %w", err)
	}

	lines = append([]string{fmt.Sprintf(":%s:", message.Username)}, lines...)

	if !incoming {
		s.ClearMessageBox()
		for _, line := range lines {
			s.buffer += fmt.Sprintln(" > ", line)
		}
	} else {
		s.ClearMessageBox()
		for _, line := range lines {
			bufferRight(s, line)
		}
	}

	fmt.Println(s.buffer)
	s.buffer = ""

	s.CreateInputBox()
	return nil
}

func (s *Screen) CreateInputBox() {
	for range s.Width {
		fmt.Printf("-")
	}
	fmt.Println()
}

func bufferRight(s *Screen, line string) {
	dLine := fmt.Sprintln(line, " < ")
	numSpaces := s.Width - len(line)
	spaces := ""
	for range numSpaces {
		spaces += " "
	}
	s.buffer += spaces + dLine
}

func (s *Screen) formatMessage(message string) ([]string, error) {
	splitMsg := strings.Split(message, " ")

	lines, err := s.splitsToLines(splitMsg)
	if err != nil {
		return make([]string, 0), err
	}

	return lines, nil
}

func (s *Screen) splitsToLines(splits []string) ([]string, error) {
	if len(splits) < 1 {
		return make([]string, 0), fmt.Errorf("no message to split")
	}
	lines := make([]string, 0, 5)

	runeCount := len(splits[0]) + spaceRune
	lastPosition := 0

	for i, split := range splits[1:] {

		if runeCount+len(split)+spaceRune+widthBuffer >= s.Width {

			line := strings.Join(splits[lastPosition:i-1], " ")
			lines = append(lines, line)
			lastPosition = i - 1
			runeCount = len(splits[i-1]) + spaceRune
		}
		runeCount += len(split) + spaceRune
	}
	line := strings.Join(splits[lastPosition:], " ")
	lines = append(lines, line)
	return lines, nil
}

func (s *Screen) PrintWelcomeMessage() {
	fmt.Println("welcome to celaeno!")
	fmt.Println("the encrypted terminal messaging app")
	fmt.Println()
	fmt.Println("commands:")
	fmt.Println()
	fmt.Println("to register a new user:")
	fmt.Println("/register <username> <password>")
	fmt.Println()
	fmt.Println("to login:")
	fmt.Println("/login <username> <password>")
	fmt.Println()
}
