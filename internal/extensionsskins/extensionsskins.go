package extensionsskins

import (
	"fmt"
	"log"
	"strings"

	"github.com/CanastaWiki/Canasta-CLI-Go/internal/config"
	"github.com/CanastaWiki/Canasta-CLI-Go/internal/orchestrators"
)

type Item struct {
	Name                     string
	RelativeInstallationPath string
	PhpCommand               string
}

func Contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func List(instance config.Installation, constants Item) {
	log.Printf("Available %s:\n", constants.Name)
	log.Print(orchestrators.Exec(instance.Path, instance.Orchestrator, "web", "cd $MW_HOME/"+constants.RelativeInstallationPath+" && find * -maxdepth 0 -type d"))
}

func CheckInstalled(name string, instance config.Installation, constants Item) (string, error) {
	output := orchestrators.Exec(instance.Path, instance.Orchestrator, "web", "cd $MW_HOME/"+constants.RelativeInstallationPath+" && find * -maxdepth 0 -type d")
	if !Contains(strings.Split(output, "\n"), name) {
		return "", fmt.Errorf("%s %s doesn't exist", constants.Name, name)
	}
	return name, nil
}

func Enable(name, wiki string, instance config.Installation, constants Item) {
	if wiki == "" {
		log.Println("You didn't specify a wiki. The extension or skin will affect all wikis in the corresponding Canasta instance.")
	}

	var filePath string

	if wiki != "" {
		filePath = fmt.Sprintf("/mediawiki/config/%s/%s.php", wiki, name)
	} else {
		filePath = fmt.Sprintf("/mediawiki/config/settings/%s.php", name)
	}

	phpScript := fmt.Sprintf("<?php\n// This file was generated by Canasta\n%s( '%s' );", constants.PhpCommand, name)

	output, err := orchestrators.ExecWithError(instance.Path, instance.Orchestrator, "web", "ls "+filePath)

	if err == nil {
		log.Printf("%s %s is already enabled!\n", constants.Name, name)
	} else if Contains(strings.Split(output, ":"), " No such file or directory\n") {
		command := fmt.Sprintf(`echo -e "%s" > %s`, phpScript, filePath)
		orchestrators.Exec(instance.Path, instance.Orchestrator, "web", command)
		log.Printf("%s %s enabled\n", constants.Name, name)
	}
}

func CheckEnabled(name, wiki string, instance config.Installation, constants Item) (string, error) {
	var settingsPath string
	var filePath string

	if wiki != "" {
		settingsPath = fmt.Sprintf("/mediawiki/config/%s/", wiki)
		filePath = fmt.Sprintf("/mediawiki/config/%s/%s.php", wiki, name)
	} else {
		settingsPath = "/mediawiki/config/settings/"
		filePath = fmt.Sprintf("/mediawiki/config/settings/%s.php", name)
	}

	output := orchestrators.Exec(instance.Path, instance.Orchestrator, "web", "ls "+settingsPath)
	if !Contains(strings.Split(output, "\n"), name+".php") {
		return "", fmt.Errorf("%s %s is not enabled", constants.Name, name)
	}

	output = orchestrators.Exec(instance.Path, instance.Orchestrator, "web", fmt.Sprintf(`cat %s`, filePath))
	if !Contains(strings.Split(output, "\n"), "// This file was generated by Canasta") {
		return "", fmt.Errorf("%s %s was not generated by Canasta cli", constants.Name, name)
	}

	return name, nil
}

func Disable(name, wiki string, instance config.Installation, constants Item) {
	if wiki == "" {
		log.Println("You didn't specify a wiki. The common settings will disable the extension or skin in the corresponding Canasta instance.")
	}

	var filePath string

	if wiki != "" {
		filePath = fmt.Sprintf("/mediawiki/config/%s/%s.php", wiki, name)
	} else {
		filePath = fmt.Sprintf("/mediawiki/config/settings/%s.php", name)
	}

	command := fmt.Sprintf(`rm %s`, filePath)
	orchestrators.Exec(instance.Path, instance.Orchestrator, "web", command)
	log.Printf("%s %s disabled\n", constants.Name, name)
}
