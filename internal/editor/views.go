package editor

import (
	"fmt"
	"nanodocker/internal/highlighting"
	"os"

	"github.com/awesome-gocui/gocui"
)

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
