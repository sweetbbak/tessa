package main

import (
	"os/exec"
)

func notifyLinux(appName string, title string, text string, iconPath string) error {
	noti, err := exec.LookPath("notify-send")
	if err != nil {
		return err
	}

	cmd := exec.Command(noti, "-i", iconPath, title, text)
	cmd.Run()

	return nil
}
