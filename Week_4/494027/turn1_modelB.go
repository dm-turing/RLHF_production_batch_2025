package main

import (
	"fmt"
	// "log"
	"time"

	"github.com/andlabs/ui"
)

type mouseEvent struct {
	desc     string
	time     time.Duration
	err      error
	next     *mouseEvent
	previous *mouseEvent
}

var (
	mouseEventChain *mouseEvent
)

func (e *mouseEvent) String() string {
	return fmt.Sprintf("Event: %s, Time: %v, Error: %v", e.desc, e.time, e.err)
}

func newMouseEvent(desc string) *mouseEvent {
	return &mouseEvent{desc: desc}
}

func addMouseEventToChain(e *mouseEvent) {
	if mouseEventChain == nil {
		mouseEventChain = e
	} else {
		e.previous = mouseEventChain
		mouseEventChain.next = e
		mouseEventChain = e
	}
}

func handleMouseEvent(desc string) {
	start := time.Now()
	// Simulate some work for this event
	time.Sleep(time.Duration(100) * time.Millisecond)

	end := time.Now()
	duration := end.Sub(start)

	// Capture possible errors
	err := fmt.Errorf("Simulated error for %s", desc)

	e := newMouseEvent(desc)
	e.time = duration
	e.err = err

	addMouseEventToChain(e)

	// Update the GUI with the latest event
	ui.QueueMain(func() {
		text := fmt.Sprintf("%s\n%s", logText.Text(), e)
		logText.SetText(text)
	})
}

func onMouseDown(ev *ui.MouseEvent) {
	handleMouseEvent("Mouse Down")
}

func onMouseMove(ev *ui.MouseEvent) {
	handleMouseEvent("Mouse Move")
}

func onMouseUp(ev *ui.MouseEvent) {
	handleMouseEvent("Mouse Up")
}

func main() {
	ui.OnShouldQuit(func() bool {
		return true
	})

	logText := ui.NewMultilineEntry()
	logText.SetReadOnly(true)

	mainWindow := ui.NewWindow("Mouse Event Capture", 300, 200, false)
	mainWindow.SetChild(logText)

	mainWindow.OnMouseDown(onMouseDown)
	mainWindow.OnMouseMove(onMouseMove)
	mainWindow.OnMouseUp(onMouseUp)

	mainWindow.Show()

	ui.Main()
}
