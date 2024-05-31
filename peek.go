package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FilePreview struct {
	path string
	show bool
}

func isDirectory(path string) bool {
	return false
}

func getFilesInPath(path string) []FilePreview {
	var filePreviewList []FilePreview

	files, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePreviewList = append(filePreviewList, FilePreview{path: file.Name(), show: true})
		}
	}

	return filePreviewList
}

func main() {
	app := tview.NewApplication()
	filePreviewList := getFilesInPath(".")

	// Left view
	filesListView := tview.NewList()

	for _, filePreview := range filePreviewList {
		if filePreview.show {
			filesListView.AddItem(filePreview.path, "", 0, nil)
		}
	}

	leftFrame := tview.NewFrame(filesListView)
	leftFrame.SetBorder(true).SetTitle("Files")
	leftFrame.AddText("Press (q) to exit.", false, tview.AlignTop, tcell.ColorGreen)
	leftFrame.AddText("Use (j) and (k) or the arrows to select file.", false, tview.AlignTop, tcell.ColorGreen)

	// Rigth view
	fileContentsView := tview.NewTextView().SetText("File Contents...")
	rightFrame := tview.NewFrame(fileContentsView)
	rightFrame.SetBorder(true).SetTitle("Preview")

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFrame, 0, 1, true).
		AddItem(rightFrame, 0, 2, true)

	// Keyboard events
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			// Exit with "q" key
			app.Stop()
		case 'j':
			// Move down with "j" key
			if app.GetFocus() == filesListView {
				currentItemIndex := filesListView.GetCurrentItem()

				if currentItemIndex+1 >= filesListView.GetItemCount() {
					filesListView.SetCurrentItem(0)
				} else {
					filesListView.SetCurrentItem(currentItemIndex + 1)
				}
			}
		case 'k':
			// Move up with "k" key
			if app.GetFocus() == filesListView {
				currentItemIndex := filesListView.GetCurrentItem()

				if currentItemIndex == 0 {
					filesListView.SetCurrentItem(filesListView.GetItemCount() - 1)
				} else {
					filesListView.SetCurrentItem(currentItemIndex - 1)
				}
			}
		case rune(tcell.KeyTAB):
			// Switch focus between left and right views with "TAB" key
			if app.GetFocus() == filesListView {
				app.SetFocus(fileContentsView)
			} else {
				app.SetFocus(filesListView)
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
