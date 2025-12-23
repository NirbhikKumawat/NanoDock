package editor

import (
	"os"

	"github.com/awesome-gocui/gocui"
)

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
