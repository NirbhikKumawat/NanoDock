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

__Ctrl+H__ - Toggle main and help  
__Ctrl+I__ - Close current help  
__Ctrl+S__ - Save file  
__Ctrl+N__ - New File  
__Ctrl+O__ - Toggle overwrite  
__Ctrl+C__ - Exit  

## Acknowledgements
- Inspired by [Nano](https://www.nano-editor.org/)
- Go Libraries used [Cobra](https://cobra.dev) and [Awesome GoCui](https://github.com/awesome-gocui)

## Roadmap
- [x] Basic editing functionality
- [x] Syntax Highlighting
- [x] Dockerfile help bar
- [ ] Undo-Redo
- [ ] Linting
- [ ] More features
- [ ] Updating README

---
Made by [NirbhikTheNice](https://github.com/NirbhikKumawat)