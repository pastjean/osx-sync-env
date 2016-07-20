package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	currentUser, _   = user.Current()
	plistFilePath, _ = filepath.Abs(currentUser.HomeDir + "/Library/LaunchAgents/osx-sync-env.plist")

	errPlistFileExists = fmt.Errorf("PlistFile '%s' already exists", plistFilePath)
	errPlistFileIsDir  = fmt.Errorf("PlistFile '%s' is a directory, go home you are drunk", plistFilePath)
)

func parseEnvironment() map[string]string {
	envsMap := make(map[string]string)

	for _, env := range os.Environ() {
		keyValue := strings.SplitN(env, "=", 2)
		envsMap[keyValue[0]] = keyValue[1]
	}
	return envsMap
}

func setEnv(key, value string) {
	launchctlSetenv := exec.Command("launchctl", "setenv", key, value)
	if err := launchctlSetenv.Run(); err != nil {
		log.Printf("Error setting env variable using '%v' \n Error:%v\n", launchctlSetenv.Args, err)
	}
}

func launchctlLoad() error {
	loadCMD := exec.Command("launchctl", "load", plistFilePath)
	return loadCMD.Run()
}
func launchctlUnLoad() error {
	loadCMD := exec.Command("launchctl", "unload", plistFilePath)
	return loadCMD.Run()
}

func launchctlReload() error {
	if err := launchctlUnLoad(); err != nil {
		return err
	}
	return launchctlLoad()
}

func sync() {
	for env, val := range parseEnvironment() {
		setEnv(env, val)
	}
	fmt.Println("Environment variables reloaded. Now relaunch your GUI apps to make them aware.")
}

func install() error {
	_, err := os.Stat(plistFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return errPlistFileExists
	}

	if err := createPlistFile(); err != nil {
		return err
	}

	return launchctlReload()
}

func uninstall() error {
	if err := launchctlUnLoad(); err != nil {
		return err
	}
	return deletePlistFile()
}

func upgrade() error {
	filestat, err := os.Stat(plistFilePath)
	if err != nil {
		return err
	}
	if filestat.IsDir() {
		return errPlistFileIsDir
	}

	// TODO: extra verification step
	// we should verify if the file was created by us to not erase
	// other softwares that would conflict with us

	if err := deletePlistFile(); err != nil {
		return err
	}
	if err := createPlistFile(); err != nil {
		return err
	}

	return launchctlReload()
}

func deletePlistFile() error {
	return os.Remove(plistFilePath)
}

func createPlistFile() error {
	file, err := os.Create(plistFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	programPath, _ := filepath.Abs(os.Args[0])
	loginShell := os.Getenv("SHELL")
	plistFileContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
	<key>Label</key>
	<string>osx-sync-env</string>
	<key>ProgramArguments</key>
	<array><string>%s</string><string>-c</string><string>-l</string><string>%s</string><string>sync</string></array>
	<key>RunAtLoad</key>
	<true/>
	</dict>
</plist>`, loginShell, programPath)

	file.WriteString(plistFileContent)
	return nil
}

func main() {

	var RootCmd = &cobra.Command{
		Use:   "osx-sync-env", // os.Args[0]
		Short: "osx-sync-env is and easy to use environment variable manager",
		Long: `An easy to use environment variable manager. It loads the environment
variables exported in the user shell into the osx GUI app context
using launchctl. Built with love by pastjean.`,
	}

	var installCmd = &cobra.Command{
		Use:   "install",
		Short: "Installs the sync on user login",
		Long: `It creates a launch command that is run on user's log. That command
launches osx-sync-env with the "sync" command (sets up the GUI app's context
environment variables).`,
		Run: func(cmd *cobra.Command, args []string) {

			if err := install(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Successfully loaded and installed %s into launchctl\n", plistFilePath)
		},
	}

	var uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Removes the sync from the user's login",
		Long: `It removes the login command from the user's login. Environment variables are
still exported and a logout+login is necessary to remove them.`,
		Run: func(cmd *cobra.Command, args []string) {

			if err := uninstall(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Successfully unloaded and uninstalled %s from launchctl\n", plistFilePath)
		},
	}

	var upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades the plist file to the most recent version.",
		Long: `Replace the plist file with a newer one. 
Used only in case the software changes the plist format`,
		Run: func(cmd *cobra.Command, args []string) {
			upgrade()
		},
	}

	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Exports the env. vars. of the current shell into the GUI context",
		Long: `This command uses "launchctl setenv" to synchronize environment variables 
set in the current shell into the OSX GUI application context.`,
		Run: func(cmd *cobra.Command, args []string) {
			sync()
		},
	}

	RootCmd.AddCommand(installCmd)
	RootCmd.AddCommand(uninstallCmd)
	RootCmd.AddCommand(upgradeCmd)
	RootCmd.AddCommand(syncCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
