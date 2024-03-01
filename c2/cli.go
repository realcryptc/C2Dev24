package c2

import (
	"fmt"
	"github.com/chzyer/readline"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var helpMenu string = `

CSC C2 V0.0.1 (2024-03-01)
_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-
Help Menu :p							 	
Commands:			Description			 
_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-
help				Show this menu		 
clear				Clear the console	 
exit				Clean Exit of C2	 
history				See what you did fool
_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-

`

func readHistoryFile() {
	data, err := os.ReadFile(".history")
	if err != nil {
		log.Fatal("Couldn't find the history file")
	}
	os.Stdout.Write(data)
}

func clearConsole() error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "pls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func installStuffOnLinux() error {
	var cmd *exec.Cmd

	//if runtime.GOOS == "windows" {
	//	cmd = exec.Command("cmd", "/c", "pls")
	//} else {
	cmd = exec.Command("apt-get")
	//}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func StartCLI() {
	autoCompleter := readline.NewPrefixCompleter(
		readline.PcItem("clear"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("history"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "C2 > ",
		AutoComplete:    autoCompleter,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	historyFile, err := os.OpenFile(".history", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer historyFile.Close()

	for {
		command, err := rl.Readline()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := historyFile.WriteString(command + "\n"); err != nil {
			fmt.Println(err)
		}

		command = strings.TrimSpace(command)

		switch command {
		case "exit":
			{
				// do a thing
				fmt.Println("[C2] Shutting down C2, gooodbye world...")
				os.Exit(0)
			}
		case "help":
			{
				fmt.Print(helpMenu)
			}
		case "clear":
			{
				if err := clearConsole(); err != nil {
					fmt.Println(err)
				}
			}
		case "history":
			// print everything from history file
			readHistoryFile()
		}

	}

}
