package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func uninstall(js string) (string, error) {
	if strings.Contains(js, beginJs) || strings.Contains(js, endJs) {
		if !*flagAttemptReinstall {
			fmt.Println("Bootstrapper is installed. Either attempt a reinstall via --reinstall or make a clean Discord reinstall")
			return "", errors.New("Bootstrapper installed")
		}

		fmt.Println("Bootstrapper already installed, removing old install")

		start := strings.LastIndex(js, beginJs)
		end := strings.LastIndex(js, endJs) + len(endJs) - 1
		js = js[:start] + js[end+1:]
		fmt.Println("Uninstall was completed. If success or not you'll find out.")
	}
	return js, nil
}

func install(js, bootstrapper string) (string, error) {
	return strings.Replace(js,
			lookFor,
			lookFor+"\n"+bootstrapper+"\n", 1),
		nil
}
