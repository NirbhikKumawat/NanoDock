package editor

import (
	"errors"
	"fmt"
	"log"
	"nanodocker/dockerfile"
	"nanodocker/internal/highlighting"
	"os"
	"regexp"

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

// Cursor navigation

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
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

// help bar

func GetLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return err
	}
	maxX, maxY := g.Size()
	if v, err := g.SetView("information", 2*maxX/3+1, 3, maxX-1, 3*maxY/4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		fmt.Fprintln(v, ColorKeyword+l+ColorReset)
		v.Wrap = true
		g.Cursor = true
		v.Editable = true
		for _, keyword := range dockerfileKeywords {
			if l == keyword {
				fmt.Fprintln(v, dockerfile.Info[keyword])
			}
		}
		if _, err := g.SetCurrentView("information"); err != nil {
			return err
		}
	}
	return nil
}
func HelpToBodyView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("body"); err == nil {
		g.Mouse = false
		g.Cursor = true
		return nil
	}
	return nil
}
func BodyToHelpView(g *gocui.Gui, v *gocui.View) error {
	if ViewExist(g, "information") {
		if _, err := g.SetCurrentView("information"); err == nil {
			v.Editable = false
			g.Mouse = false
			g.Cursor = false
			return err
		}
	} else {
		if v, err := g.SetCurrentView("help"); err == nil {
			v.Editable = false
			g.Mouse = false
			g.Cursor = false
			return err
		}
	}
	return nil
}
func DeleteInformationView(g *gocui.Gui, v *gocui.View) error {
	if !ViewExist(g, "information") {
		return nil
	}
	if err := g.DeleteView("information"); err != nil {
		return err
	}
	if err := BodyToHelpView(g, v); err != nil {
		return err
	}
	return nil
}
func ViewExist(g *gocui.Gui, s string) bool {
	if _, err := g.View(s); err != nil {
		return false
	}
	return true
}

// new file

func NewFile(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	return nil
}
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// save file

func SaveView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("savename", 0, maxY/2-1, maxX, maxY/2+1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
		v.Editable = true
	}
	if _, err := g.SetCurrentView("savename"); err != nil {
		return err
	}
	return nil
}
func SaveDeleteView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("savename"); err != nil {
		return err
	}
	s := v.Buffer()
	re := regexp.MustCompile(`\s+`)
	file = re.ReplaceAllString(s, "")
	v, err := g.SetCurrentView("body")
	if err != nil {
		return err
	}
	if err := SaveMain(g, v); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("body", 0, 3, 2*maxX/3, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dockerfile"
		v.Editable = true
		v.Wrap = true
		v.Editor = DockerfileEditor{}
	}
	if err := g.DeleteView("title"); err != nil {
		return err
	}
	if v, err := g.SetView("title", 0, 0, maxX-1, 2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
	}
	return nil
}
func SaveMain(g *gocui.Gui, v *gocui.View) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, "Dockerfile")
	if err != nil {
		return err
	}

	content := highlighting.StripAnsiCodes(v.Buffer())

	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return err
	}

	if err := os.Rename(f.Name(), file); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlI, gocui.ModNone, DeleteInformationView); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyEnter, gocui.ModNone, GetLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlO, gocui.ModNone, Overwrite); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlN, gocui.ModNone, NewFile); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlS, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return SaveView(g)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("savename", gocui.KeyEnter, gocui.ModNone, SaveDeleteView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.MouseWheelDown, gocui.ModNone, Nothing); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.MouseWheelUp, gocui.ModNone, Nothing); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyCtrlH, gocui.ModNone, HelpToBodyView); err != nil {
		return err
	}
	if err := g.SetKeybinding("information", gocui.KeyCtrlH, gocui.ModNone, HelpToBodyView); err != nil {
		return err
	}
	/*if err := g.SetKeybinding("information", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("information", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}*/
	if err := g.SetKeybinding("body", gocui.KeyCtrlH, gocui.ModNone, BodyToHelpView); err != nil {
		return err
	}
	return nil
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
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("title", 0, 0, maxX-1, 2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
	}
	if v, err := g.SetView("body", 0, 3, 2*maxX/3, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dockerfile"
		v.Editable = true
		v.Wrap = true
		v.Editor = DockerfileEditor{}

		if _, err := g.SetCurrentView("body"); err != nil {
			return err
		}
		if FileExists(file) {
			b, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}

			content := string(b)
			highlighted := highlighting.HighlightDockerfile(content)
			fmt.Fprint(v, highlighted)
		}

	}
	if v, err := g.SetView("help", 2*maxX/3+1, 3, maxX-1, 3*maxY/4, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "\033[32mADD")
		fmt.Fprintln(v, "ARG")
		fmt.Fprintln(v, "CMD")
		fmt.Fprintln(v, "COPY")
		fmt.Fprintln(v, "ENTRYPOINT")
		fmt.Fprintln(v, "ENV")
		fmt.Fprintln(v, "EXPOSE")
		fmt.Fprintln(v, "FROM")
		fmt.Fprintln(v, "HEALTHCHECK")
		fmt.Fprintln(v, "LABEL")
		fmt.Fprintln(v, "MAINTAINER")
		fmt.Fprintln(v, "ONBUILD")
		fmt.Fprintln(v, "RUN")
		fmt.Fprintln(v, "SHELL")
		fmt.Fprintln(v, "STOPSIGNAL")
		fmt.Fprintln(v, "USER")
		fmt.Fprintln(v, "VOLUME")
		fmt.Fprintln(v, "WORKDIR\033[0m")
	}
	if v, err := g.SetView("command", 2*maxX/3+1, 3*maxY/4+1, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Commands for using the editor")
		fmt.Fprintln(v, "Ctrl-C: Exit")
		fmt.Fprintln(v, "Ctrl-S: Save")
		fmt.Fprintln(v, "Ctrl-N: New file")
		fmt.Fprintln(v, "Ctrl-O: Toggle Overwrite")
		fmt.Fprintln(v, "Ctrl-H: Toggle Help")
		fmt.Fprintln(v, "Ctrl-I: Close information tab")

	}
	return nil
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
