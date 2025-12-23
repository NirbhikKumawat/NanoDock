package editor

import (
	"fmt"
	"nanodocker/internal/highlighting"
	"os"
	"regexp"

	"github.com/awesome-gocui/gocui"
)

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
