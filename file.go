package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FileEntry struct {
	NameLabel *widget.Label
	Button    *widget.Button
}

func main() {
	a := app.New()
	w := a.NewWindow("设置运行路径")
	w.Resize(fyne.NewSize(600, 400)) // 设置窗口宽度为600像素，高度为400像素
	createRow := func(index int) FileEntry {
		nameLabel := widget.NewLabel("xshell路径")
		button := widget.NewButton("Selected exe", func() {
			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				defer reader.Close()
				fmt.Println("Selected file:", reader.URI())
			}, w)
		})
		return FileEntry{
			NameLabel: nameLabel,
			Button:    button,
		}
	}
	// Create four rows with different names
	rows := make([]FileEntry, 4)
	for i := range rows {
		rows[i] = createRow(i)
	}
	table := container.NewGridWithColumns(1)
	for _, row := range rows {
		rowContainer := container.NewHBox(row.NameLabel, row.Button)
		table.Add(rowContainer)
	}

	w.SetContent(table)
	w.ShowAndRun()
}
