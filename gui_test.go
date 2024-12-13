package main

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/jroimartin/gocui"
)

func TestGui(t *testing.T) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layoutT)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitT); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type rectangle struct {
	x0 int
	x1 int
	y0 int
	y1 int
}

func addNewLine(arr []string) []string {
	newArr := []string{}
	for i := range arr {
		newArr = append(newArr, "\n" + arr[i]) 
	}
	return newArr
}

func longestStrLength(arr []string) int {
	maxLength := len(arr[0])
	for _, v := range arr {
		if maxLength < len(v) {
			maxLength = len(v)
		}
	}
	return maxLength
}

func drawRecWrite(g *gocui.Gui, name string, coordinates rectangle, toWrite string) (*gocui.View, error) {
	v, err := g.SetView(name, coordinates.x0, coordinates.y0, coordinates.x1, coordinates.y1);

	if  err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		fmt.Fprintln(v, toWrite)
	}
	return v, nil
}

func layoutT(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	chatsNames := []string{"General", "Annoucement", "Memes", "Anime"}
	membersNames := []string{"Alice", "Bob", "Francesca"}
	chatsR := rectangle{0, longestStrLength(chatsNames)+1, 0, maxY-1}
	messagesR := rectangle{chatsR.x1, maxX-longestStrLength(membersNames)-1-1, 0, maxY-4}
	membersR := rectangle{messagesR.x1, maxX-1, 0, chatsR.y1}
	writeMessageR := rectangle{chatsR.x1, membersR.x0, messagesR.y1, chatsR.y1}

	_, errV := drawRecWrite(g, "chats", chatsR, strings.Join(chatsNames, "\n"))
	if errV != nil {
		return errV
	}

	if v, err := g.SetView("messages", messagesR.x0, messagesR.y0, messagesR.x1, messagesR.y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
		fmt.Fprintln(v, maxX, maxY)
	}
	if v, err := g.SetView("members", membersR.x0, membersR.y0, membersR.x1 , membersR.y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, strings.Join(membersNames, "\n"))
	}
	if v, err := g.SetView("write message", writeMessageR.x0, writeMessageR.y0, writeMessageR.x1, writeMessageR.y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quitT(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
