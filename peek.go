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

func getFileContents(path string) string {
	contents, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(contents)
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

	// Display selected file contents
	if len(filePreviewList) > 0 {
		fileContents := getFileContents(filePreviewList[0].path)
		fileContentsView.SetText(fileContents)

		filesListView.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			if filePreviewList[index].show {
				fileContents := getFileContents(filePreviewList[index].path)
				fileContentsView.SetText(fileContents)
			}
		})
	}

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

			return nil
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
