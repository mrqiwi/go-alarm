package gui

import (
	"fmt"
	"strconv"
	"time"

	"go-alarm/pkg/alarm"
	"go-alarm/pkg/audio"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func InitGUI(player *audio.Player, alarmer *alarm.Alarm) fyne.Window {
	var (
		timeHours, timeMinutes, deferMinutes int
		setTime                              time.Time
		fileName                             string
	)

	myApp := app.New()
	mainWindow := myApp.NewWindow("Alarm")
	mainWindow.Resize(fyne.Size{
		Width:  320,
		Height: 240,
	})
	mainWindow.SetFixedSize(true)

	hoursVariants := []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}
	hoursSelect := widget.NewSelect(hoursVariants, func(value string) {
		timeHours, _ = strconv.Atoi(value)
	})
	hoursSelect.SetSelected("00")

	minutesVariants := []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25",
		"26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
		"40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53",
		"54", "55", "56", "57", "58", "59"}
	minutesSelect := widget.NewSelect(minutesVariants, func(value string) {
		timeMinutes, _ = strconv.Atoi(value)
	})
	minutesSelect.SetSelected("00")

	notifyLabel := widget.NewLabel("")

	browseButton := widget.NewButton("Choose...", func() {
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, _ error) {
			fileName = uc.URI().Path()
			mainWindow.Resize(fyne.Size{Width: 320, Height: 240})
		}, mainWindow)

		mainWindow.Resize(fyne.Size{Width: 800, Height: 600})

		fd.Resize(fyne.Size{Width: 1024, Height: 768})
		fd.Show()
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))

		notifyLabel.Hide()
	})

	deferSelect := widget.NewSelect([]string{"1", "2", "5", "10"}, func(value string) {
		deferMinutes, _ = strconv.Atoi(value)
	})
	deferSelect.SetSelected("1")
	deferSelect.Disable()

	deferAlarmButton := widget.NewButton("Defer", func() {
		alarmer.Reset()
		player.Stop()

		setTime = setTime.Add(time.Minute * time.Duration(deferMinutes))
		alarmer.SetUp(setTime, func() {
			player.Play(fileName)
		})

		notifyLabel.SetText(fmt.Sprintf("The alarm will ring at %d:%d", setTime.Hour(), setTime.Minute()))
		notifyLabel.Show()
	})
	deferAlarmButton.Disable()

	setAlarmButton := widget.NewButton("Start", func() {
		if fileName == "" {
			notifyLabel.SetText("Ringtone not selected!")
			notifyLabel.Show()
			return
		}

		deferSelect.Enable()
		deferAlarmButton.Enable()

		hoursSelect.Disable()
		minutesSelect.Disable()
		browseButton.Disable()

		setTime = hoursMinutesToTime(timeHours, timeMinutes)

		alarmer.SetUp(setTime, func() {
			player.Play(fileName)
		})
	})

	resetAlarmButton := widget.NewButton("Reset", func() {
		alarmer.Reset()
		player.Stop()

		deferSelect.Disable()
		deferAlarmButton.Disable()

		hoursSelect.Enable()
		minutesSelect.Enable()
		browseButton.Enable()
		notifyLabel.Hide()
	})

	mainWindow.SetContent(container.NewGridWithRows(4,
		container.NewGridWithColumns(3,
			container.NewVBox(hoursSelect),
			container.NewVBox(minutesSelect),
			container.NewVBox(browseButton),
		),
		container.NewVBox(container.NewGridWithColumns(2, setAlarmButton, resetAlarmButton)),
		container.NewVBox(container.NewGridWithColumns(2, deferAlarmButton, deferSelect)),
		container.NewVBox(notifyLabel),
	),
	)

	return mainWindow
}

func hoursMinutesToTime(h, m int) time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), h, m, now.Second(), now.Nanosecond(), now.Location())
}
