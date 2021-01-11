package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type FileHandle struct {
	Directory string
	Name      string
	movedTo   string
}

func (f FileHandle) GetPath() string {
	return path.Join(f.Directory, f.movedTo, f.Name)
}

var index int
var size fyne.Size
var files []FileHandle

func FilesLeft() bool {
	return index < len(files)
}

func remove(s []FileHandle, i int) (FileHandle, []FileHandle) {
	e := s[i]
	return e, append(s[:i], s[i+1:]...)
}

func main() {

	index = 0

	args := os.Args[1:]

	basePath := "./"

	_ = os.Mkdir(path.Join(basePath, "deleted"), os.ModeDir)

	fileinfo, _ := ioutil.ReadDir(basePath)

	for _, f := range fileinfo {
		if !(strings.HasSuffix(f.Name(), ".png") || strings.HasSuffix(f.Name(), ".jpg") || strings.HasSuffix(f.Name(), ".jpeg")) {
			continue
		}
		fmt.Println(path.Join(basePath, f.Name()))
		files = append(files, FileHandle{Directory: basePath, Name: f.Name()})
	}

	if !FilesLeft() {
		fmt.Println("No Images Could be found in this Directory")
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	a := app.New()
	//a.SendNotification(fyne.NewNotification("Test", "Testing"))
	w := a.NewWindow("Image Sorter")
	w.SetFullScreen(true)
	w.Show()

	buttons := []fyne.CanvasObject{}

	buttons = append(buttons, widget.NewButton("Back", func() {
		index -= 1
		if index < 0 {
			index = 0
		}
		BuildContent(w, buttons, files[index])
	}))

	buttons = append(buttons, widget.NewButton("Skip", func() {
		index += 1
		if index >= len(files) {
			index = len(files) - 1
		}
		BuildContent(w, buttons, files[index])
	}))

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {

		switch k.Name {
		case "Left":
			index -= 1
			if index < 0 {
				index = 0
			} else {

				BuildContent(w, buttons, files[index])
			}
			return
		case "Right":
			index += 1
			if index >= len(files) {
				index = len(files) - 1
			} else {

				BuildContent(w, buttons, files[index])
			}
			return
		case "Delete":
			var file FileHandle
			file, files = remove(files, index)

			oldpath := file.GetPath()
			file.movedTo = "deleted"
			newpath := file.GetPath()

			err := os.Rename(oldpath, newpath)

			if err != nil {
				fmt.Errorf(err.Error())
			}

			if !FilesLeft() {
				w.Close()
				return
			}

			BuildContent(w, buttons, files[index])
			return
		case "F":
			w.SetFullScreen(!w.FullScreen())
			return
		case "B":

			oldpath := files[index].GetPath()
			files[index].movedTo = ""
			newpath := files[index].GetPath()

			err := os.Rename(oldpath, newpath)

			if err != nil {
				fmt.Errorf(err.Error())
			}

			if !FilesLeft() {
				w.Close()
				return
			}

			BuildContent(w, buttons, files[index])
			return
		}

		char := fmt.Sprintf("%s", k.Name)
		if localIndex, err := strconv.Atoi(char); err == nil {
			if localIndex >= len(args) {
				return
			}
			tmparg := args[localIndex]

			oldpath := files[index].GetPath()
			files[index].movedTo = tmparg
			newpath := files[index].GetPath()

			err := os.Rename(oldpath, newpath)

			if err != nil {
				fmt.Errorf(err.Error())
			}

			index += 1
			if !FilesLeft() {
				w.Close()
				return
			}
			BuildContent(w, buttons, files[index])
		}
	})

	for i, arg := range args {
		tmparg := arg
		_ = os.Mkdir(path.Join(basePath, tmparg), os.ModeDir)
		buttons = append(buttons, widget.NewButton(fmt.Sprintf("(%d)%s", i, tmparg), func() {
			oldpath := files[index].GetPath()
			files[index].movedTo = tmparg
			newpath := files[index].GetPath()

			err := os.Rename(oldpath, newpath)

			if err != nil {
				fmt.Errorf(err.Error())
			}

			index += 1
			if !FilesLeft() {
				w.Close()
				return
			}
			BuildContent(w, buttons, files[index])
		}))
	}

	BuildContent(w, buttons, files[index])

	w.ShowAndRun()
}

func BuildContent(w fyne.Window, buttons []fyne.CanvasObject, file FileHandle) {
	image := &canvas.Image{
		File: file.GetPath(),
		// FillMode: canvas.ImageFillOriginal,
		FillMode: canvas.ImageFillContain,
	}

	size = w.Canvas().Size()
	width, height := GetSize(file.GetPath())

	image.SetMinSize(fyne.Size{Width: fyne.Min(width, size.Width), Height: fyne.Min(height, size.Height)})

	builtbuttons := widget.NewHBox(
		buttons...,
	)

	w.SetContent(
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewMaxLayout(), image),
			fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel(file.GetPath()+"  "+fmt.Sprintf("%d/%d", index, len(files)))),
			fyne.NewContainerWithLayout(layout.NewCenterLayout(), builtbuttons),
		),
	)
}

func GetSize(file string) (width int, height int) {
	handle, _ := os.Open(file)
	defer handle.Close()

	pixels, _, err := image.Decode(handle)

	if err != nil {
		fyne.LogError("image err", err)

		return 0, 0
	}
	origSize := pixels.Bounds().Size()
	return origSize.X, origSize.Y
}
