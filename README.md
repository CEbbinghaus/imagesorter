# ImageSorter

A little Project i wrote in GO using the [fyne](https://github.com/fyne-io/fyne) Libary to create a Image Sorting Program.

# Installation

Simply Download the Executable and place into a folder in your PATH

# Usage

The Program uses Commandline Arguments to Create the Folders to sort images into. They are numbered from the first argument being 0 all the way to 9

```bash
$ imagesorter folder1 folder2 folder3
```

Images will be sorted into the Folders you provide to using descriptive Names is Recommended

# Shortcuts

- Del - Delete Current Image \(places it in a Deleted folder)
- B - Move the Image back into Root Dir
- Right Arrow - Skip to next Image
- Left Arrow - Go back to Previous Image
- 0-9 - Places Image in Respective Folder
- F - Toggles Fullscreen mode on and off
