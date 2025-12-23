package editor

import "github.com/awesome-gocui/gocui"

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
