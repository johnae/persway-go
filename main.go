package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	i3 "github.com/johnae/go-i3"
)

func focusedLayout() string {
	t, _ := i3.GetTree()
	nodeWithLayout := t.Root.FindFocused(func(node *i3.Node) bool {
		return node.Layout != "none"
	})
	return string(nodeWithLayout.Layout)
}

func updateOpacity(layout string) {
	if layout == "tabbed" || layout == "stacked" {
		exec.Command("sway", "[tiling] opacity 1;").CombinedOutput()
	} else {
		exec.Command("sway", "[tiling] opacity 0.78; opacity 1").CombinedOutput()
	}
}

func resetOpacity() {
	exec.Command("sway", "[tiling] opacity 1;").CombinedOutput()
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		t, _ := i3.GetTree()
		if len(t.Root.Nodes) > 1 {
			resetOpacity()
		}
		os.Exit(0)
	}()
	recv := i3.Subscribe(i3.WindowEventType, i3.BindingEventType)
	resetOpacity()
	updateOpacity(focusedLayout())
	for recv.Next() {
		event := recv.Event()
		if ev, ok := event.(*i3.WindowEvent); ok {
			if ev.Change == "focus" {
				updateOpacity(focusedLayout())
			}
		} else if ev, ok := event.(*i3.BindingEvent); ok {
			cmd := strings.Split(ev.Binding.Command, " ")[0]
			if cmd == "layout" {
				updateOpacity(focusedLayout())
			}
		}
	}
	log.Fatal(recv.Close())
}
