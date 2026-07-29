package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	launcher := app.NewWithID("studio.corvus.launcher")
	window := launcher.NewWindow("Corvus Studio Launcher")
	window.SetContent(container.NewCenter(widget.NewLabel("Corvus Studio repository bootstrap")))
	window.Resize(fyne.NewSize(420, 180))
	window.ShowAndRun()
}
