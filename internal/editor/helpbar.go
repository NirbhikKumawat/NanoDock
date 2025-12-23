package editor

import (
	"errors"
	"fmt"
	"nanodocker/docsreference"

	"github.com/awesome-gocui/gocui"
)

func GetLine(g *gocui.Gui, v *gocui.View) error {
	docsreference.InitializeMap()
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
				fmt.Fprintln(v, docsreference.Info[keyword])
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
