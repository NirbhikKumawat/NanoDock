# NanoDock
A lightweight,fast and extensible terminal-based text editor built with GoLang

## Features
- **Syntax Highlighting** - Dynamically highlights syntax for your Dockerfiles
- **Help Bar** - In hand Dockerfile documentation
- **Performance** - Very lightweight offering great performance

## Quick Start
```bash
# clone the git repository
git clone https://github.com/NirbhikKumawat/NanoDock.git
# change directory
cd NanoDock
# build the go binary
go build -o nanodock main.go
# move to binary so it is included int the PATH
sudo mv nanodock /usr/local/bin/
# make it executable
sudo chmod+x /usr/local/bin/nanodock
# verify installation
nanodock -h
```

## Usage

### Basic Commands

### Keyboard Shortcuts

## Acknowledgements
- Inspired by [Nano](https://www.nano-editor.org/)
- Go Libraries used [Cobra](https://cobra.dev) and [Awesome GoCui](https://github.com/awesome-gocui)

## Roadmap
- [x] Basic editing functionality
- [x] Syntax Highlighting
- [ ] Undo-Redo
- [ ] Dockerfile help bar
- [ ] Linting
- [ ] More features
- [ ] Updating README

---
Made by [NirbhikTheNice](https://github.com/NirbhikKumawat)