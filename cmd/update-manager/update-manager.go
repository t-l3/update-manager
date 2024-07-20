package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/0xAX/notificator"
	"github.com/t-l3/update-manager/internal/config"
	"github.com/t-l3/update-manager/internal/manager"

	"github.com/t-l3/update-manager/internal/notifications"

	"fyne.io/systray"
)

func main() {
	notif := notifications.New("update-manager", "Started update-manager")

	logger := log.New(os.Stdout, "app-manager-main  ", log.Ldate|log.Ltime|log.Lmsgprefix)
	appConfig := config.LoadConfig()

	err := os.MkdirAll(appConfig.TmpDownloadLocation, 0775)
	if err != nil {
		logger.Fatal("Error while creating download directory", err)
	}

	logger.Printf("  === Starting app checks ===  ")

	var wg sync.WaitGroup

	for _, app := range appConfig.Apps {
		wg.Add(1)
		go updateApplication(&app, &appConfig.TmpDownloadLocation, &wg)
	}

	configureSystray := func() { systrayOnReady(appConfig.SystrayIcon) }
	startSystray, _ := systray.RunWithExternalLoop(configureSystray, func() {})
	startSystray()

	wg.Wait()
	logger.Println("App updates completed")
	logger.Println("Removing temporary files")
	os.RemoveAll(appConfig.TmpDownloadLocation)
	notif.Terminate("")
}

func updateApplication(appConfig *config.App, tmpDir *string, wg *sync.WaitGroup) {
	logger := log.New(os.Stdout, fmt.Sprintf("app-manager-%s  ", appConfig.Name), log.Ldate|log.Ltime|log.Lmsgprefix)
	notify := notificator.New(notificator.Options{
		DefaultIcon: appConfig.Icon,
		AppName:     "update-manager",
	})
	m := manager.New(appConfig, tmpDir, logger, notify)

	shouldInstall := m.CheckVersion()
	if shouldInstall {
		m.DownloadApp()
		m.InstallApp()
	}
	wg.Done()
}

func systrayOnReady(icon string) {
	updateIcon, _ := os.Open(icon)
	updateIconBytes, _ := io.ReadAll(updateIcon)
	systray.SetIcon(updateIconBytes)
	updateIcon.Close()

	systray.SetTitle("update-manager")
	systray.SetTooltip("update-manager")
	quitButton := systray.AddMenuItem("Quit", "Quit update-manager")
	go func() {
		<-quitButton.ClickedCh
		systray.Quit()
	}()
}
