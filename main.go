package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

var (
	file string
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
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
		if fileExists(file) {
			b, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(v, "%s", b)
		}

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
		fmt.Fprintln(v, "Ctrl-M: Save As")
		fmt.Fprintln(v, "Ctrl-O: Toggle Overwrite")
	}
	return nil
}

func saveAs(g *gocui.Gui, v *gocui.View) error {
	f, err := os.Create("dockerfile")
	if err != nil {
		return err
	}
	defer f.Close()
	p := make([]byte, 100)
	v.Rewind()
	if _, err := io.CopyBuffer(f, v, p); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

func newFile(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
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

	if err := os.Rename(f.Name(), file); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func overwrite(g *gocui.Gui, v *gocui.View) error {
	v.Overwrite = true
	return nil
}
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlM, gocui.ModNone, saveAs); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlO, gocui.ModNone, overwrite); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlN, gocui.ModNone, newFile); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlA, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return saveView(g)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("savename", gocui.KeyEnter, gocui.ModNone, saveDeleteView); err != nil {
		return err
	}
	return nil
}
func saveDeleteView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("savename"); err != nil {
		return err
	}
	file = v.Buffer()
	if err := saveMain(g, v); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("body"); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("body", 0, 3, 2*maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dockerfile"
		v.Editable = true
		v.Wrap = true
	}
	if err := g.DeleteView("title"); err != nil {
		return err
	}
	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
	}
	return nil
}

func runGocui(cmd *cobra.Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file = dir + "/Dockerfile"
	if len(args) > 0 {
		file = dir + "/" + args[0]
	}
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

func saveView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("savename", 0, maxY/2-2, maxX, maxY/2+2)
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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "nanodock [file]",
		Short: "Terminal based Dockerfile Editor",
		Long:  "A Terminal based Dockerfile Editor",
		//Args:  cobra.MinimumNArgs(1),
		Run: runGocui,
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
