package editor

import (
	"fmt"
	"log"
	"nanodocker/internal/highlighting"
	"os"

	"github.com/awesome-gocui/gocui"
	"github.com/spf13/cobra"
)

const (
	ColorReset   = "\033[0m"
	ColorKeyword = "\033[1;32m"
	ColorComment = "\033[0;36m"
	ColorString  = "\033[0;34m"
	ColorCode    = "\033[0;33m"
)

var (
	file string
)
var dockerfileKeywords = []string{
	"FROM", "RUN", "CMD", "LABEL", "MAINTAINER", "EXPOSE", "ENV",
	"ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", "WORKDIR",
	"ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL",
}

type DockerfileEditor struct {
}

// Syntax Highlighting

func (e DockerfileEditor) ApplySyntaxHighlighting(v *gocui.View) {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	content := highlighting.StripAnsiCodes(v.Buffer())
	highlighted := highlighting.HighlightDockerfile(content)

	v.Clear()
	fmt.Fprint(v, highlighted)

	v.SetCursor(cx, cy)
	v.SetOrigin(ox, oy)
}

// utility functions

func Overwrite(g *gocui.Gui, v *gocui.View) error {
	v.Overwrite = !v.Overwrite
	return nil
}
func Nothing(g *gocui.Gui, v *gocui.View) error {
	return nil
}
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func (e DockerfileEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		e.ApplySyntaxHighlighting(v)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		e.ApplySyntaxHighlighting(v)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		e.ApplySyntaxHighlighting(v)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		e.ApplySyntaxHighlighting(v)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		v.EditNewLine()
		e.ApplySyntaxHighlighting(v)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0)
	case key == gocui.KeyHome:
		_, cy := v.Cursor()
		v.SetCursor(0, cy)
	case key == gocui.KeyEnd:
		_, cy := v.Cursor()
		line, _ := v.Line(cy)
		v.SetCursor(len(line), cy)
	}
}
func RunGocui(cmd *cobra.Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file = dir + "/Dockerfile"
	if len(args) > 0 {
		file = dir + "/" + args[0]
	}
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = false
	g.SetManagerFunc(Layout)
	if err := Keybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
