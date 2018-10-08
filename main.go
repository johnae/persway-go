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

func currentWorkspace() *i3.Node {
	t, _ := i3.GetTree()
	focusedWs := t.Root.FindFocused(func(node *i3.Node) bool {
		if node.Type == i3.WorkspaceNode {
			n := node.FindFocused(func(node *i3.Node) bool {
				return node.Focused
			})
			return n != nil
		}
		return false
	})
	return focusedWs
}

func currentLayout(ws *i3.Node) string {
	if ws != nil {
		return string(ws.Layout)
	}
	return ""
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
	ws := currentWorkspace()
	resetOpacity()
	updateOpacity(currentLayout(ws))
	for recv.Next() {
		event := recv.Event()
		if ev, ok := event.(*i3.WindowEvent); ok {
			if ev.Change == "focus" {
				ws := currentWorkspace()
				updateOpacity(currentLayout(ws))
			}
		} else if ev, ok := event.(*i3.BindingEvent); ok {
			cmd := strings.Split(ev.Binding.Command, " ")[0]
			if cmd == "layout" {
				ws := currentWorkspace()
				updateOpacity(currentLayout(ws))
			}
		}
	}
	log.Fatal(recv.Close())
}
