package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-isatty"
	"github.com/rivo/tview"
)

func getFilesInPath(path string) []string {
	var filePreviewList []string

	files, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !isDir(file.Name()) {
			// fmt.Printf("File: %v", file.Name())
			filePreviewList = append(filePreviewList, file.Name())
		}
	}

	return filePreviewList
}

func getFileContents(path string) string {
	contents, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error reading path: '%v\n'", path)
		panic(err)
	}

	return string(contents)
}

func isDir(path string) bool {
	fileInfo, err := os.Lstat(path)

	if err != nil {
		fmt.Printf("Error checking if path is directory: '%v'", err)

		return true
	}

	// Ignore symbolic links
	if fileInfo.Mode()&fs.ModeSymlink != 0 {
		return true
	}

	return fileInfo.IsDir()
}

func main() {
	var filePathsList []string

	app := tview.NewApplication()
	defaultBorderColor := tcell.ColorWhite
	focusBorderColor := tcell.ColorGreen

	// Overriding tview default styles
	tview.Borders.HorizontalFocus = tview.BoxDrawingsLightHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsLightVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsLightDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsLightDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsLightUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsLightUpAndLeft

	// Check if it has stdin, e.g. piped input from other command
	hasStdin := !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd())

	if hasStdin {
		// fmt.Println("Reading from stdin...")
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			path := scanner.Text()

			if !isDir(path) {
				filePathsList = append(filePathsList, path)
			}
		}
	} else {
		// fmt.Println("Reading from current directory...")
		filePathsList = getFilesInPath(".")
	}

	// Left view
	filesListView := tview.NewList()

	for _, filePath := range filePathsList {
		filesListView.AddItem(filePath, "", 0, nil)
	}

	leftFrame := tview.NewFrame(filesListView)
	leftFrame.SetBorder(true).SetTitle("Files")
	leftFrame.SetBorderColor(focusBorderColor)
	leftFrame.AddText("Press (q) to exit.", false, tview.AlignTop, tcell.ColorGreen)
	leftFrame.AddText("Use (j) and (k) or the arrows to select file.", false, tview.AlignTop, tcell.ColorGreen)

	// Rigth view
	filePreviewContents := "Loading..."
	fileContentsView := tview.NewTextView().SetText(filePreviewContents)
	rightFrame := tview.NewFrame(fileContentsView)
	rightFrame.SetBorder(true).SetTitle("Preview")

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFrame, 0, 1, true).
		AddItem(rightFrame, 0, 2, true)

	// Display selected file contents
	if len(filePathsList) > 0 {
		filePreviewContents = getFileContents(filePathsList[0])
		fileContentsView.SetText(filePreviewContents)

		filesListView.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			filePreviewContents = "Loading..."

			filePreviewContents = getFileContents(filePathsList[index])
			fileContentsView.SetText(filePreviewContents)
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
				leftFrame.SetBorderColor(defaultBorderColor)
				rightFrame.SetBorderColor(focusBorderColor)
			} else {
				app.SetFocus(filesListView)
				rightFrame.SetBorderColor(defaultBorderColor)
				leftFrame.SetBorderColor(focusBorderColor)
			}

			return nil
		}

		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
