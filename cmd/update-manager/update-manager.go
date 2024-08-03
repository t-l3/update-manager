package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/0xAX/notificator"
	"github.com/t-l3/update-manager/internal/config"
	"github.com/t-l3/update-manager/internal/manager"

	"github.com/t-l3/update-manager/internal/notifications"

	"fyne.io/systray"
)

func main() {
	logger := log.New(os.Stdout, "app-manager-main  ", log.Ldate|log.Ltime|log.Lmsgprefix)
	appConfig := config.LoadConfig()

	err := os.MkdirAll(appConfig.TmpDownloadLocation, 0775)
	if err != nil {
		logger.Fatal("Error while creating download directory", err)
	}
	notif := notifications.New("update-manager", appConfig.SystrayIcon)

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

	removeTmp := true

	for _, app := range appConfig.Apps {
		if app.RetainDownload {
			removeTmp = false
		}
	}

	if removeTmp {
		logger.Println("Removing temporary files")
		os.RemoveAll(appConfig.TmpDownloadLocation)
	}

	notif.Terminate("")
	time.Sleep(20 * time.Millisecond) // Sleep to allow notification to terminate gracefully
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
		err := m.DownloadApp()
		if err != nil {
			logger.Printf("Download of %s failed", appConfig.Name)
			return
		}
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
