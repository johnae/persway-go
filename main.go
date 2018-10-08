package main

import (
	"log"
	"os/exec"

	i3 "github.com/johnae/go-i3"
)

func main() {
	recv := i3.Subscribe(i3.WindowEventType)
	exec.Command("sway", "[title=\".*\"] opacity 1.0").CombinedOutput()
	for recv.Next() {
		ev := recv.Event().(*i3.WindowEvent)
		if ev.Change == "focus" {
			log.Printf("change: %s\n", ev.Change)

			exec.Command("sway", "[title=\".*\"] opacity 0.78").CombinedOutput()

			exec.Command("sway", "opacity", "1").CombinedOutput()
		}
	}
	log.Fatal(recv.Close())
}
