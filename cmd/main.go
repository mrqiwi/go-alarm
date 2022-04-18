package main

import (
	"go-alarm/pkg/alarm"
	"go-alarm/pkg/audio"
	"go-alarm/pkg/gui"
)

func main() {
	player := audio.NewPlayer()
	alarmer := alarm.NewAlarm()
	mainWindow := gui.InitGUI(player, alarmer)

	mainWindow.ShowAndRun()
}
