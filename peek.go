package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Left view
	leftView := tview.NewTextView().SetText("Left Side")
	leftFrame := tview.NewFrame(leftView)
	leftFrame.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Files")

	// Rigth view
	rightView := tview.NewTextView().SetText("Right Side")
	rightFrame := tview.NewFrame(rightView)
	rightFrame.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Preview")

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
		case rune(tcell.KeyTAB):
			// Switch focus between left and right views with "TAB" key
			if app.GetFocus() == leftView {
				app.SetFocus(rightView)
			} else {
				app.SetFocus(leftView)
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
