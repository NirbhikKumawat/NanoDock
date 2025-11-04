package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Filename")
	}
	if v, err := g.SetView("body", 0, 3, 2*maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dockerfile"
		v.Editable = true
		v.Wrap = true

		if _, err := g.SetCurrentView("body"); err != nil {
			return err
		}

		b, err := os.ReadFile("Dockerfile")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
		fmt.Fprintln(v, "Dockerfile Area")
	}
	if v, err := g.SetView("help", 2*maxX/3+1, 3, maxX-1, 3*maxY/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Help regarding dockerfile")
	}
	if v, err := g.SetView("command", 2*maxX/3+1, 3*maxY/4+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Commands for using the editor")
		fmt.Fprintln(v, "Ctrl-C: Exit")
		fmt.Fprintln(v, "Ctrl-S: Save")
	}
	return nil
}
func saveMain(g *gocui.Gui, v *gocui.View) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, "Dockerfile")
	if err != nil {
		return err
	}

	p := make([]byte, 512)
	v.Rewind()
	if _, err := io.CopyBuffer(f, v, p); err != nil {
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

	if err := os.Rename(f.Name(), filepath.Join(dir, "Dockerfile")); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
		return err
	}
	return nil
}
func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
