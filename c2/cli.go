package c2

import (
	"fmt"
	"github.com/chzyer/readline"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var helpMenu string = `

CSC C2 V0.0.1 (2024-03-01)
_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_
Help Menu :p							 	
Commands:			Description			 
_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_
help				Show this menu		 
clear				Clear the console	 
exit				Clean Exit of C2
agent <agent_name>	Interact w/ target agent
exec <command>		Execute command with current agent set by agent command
agents				List all the agents
history				See what you did fool
_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_
`

var CurrentAgent *Agent = nil

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
		readline.PcItem("agents"),
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
		if CurrentAgent != nil {
			rl.SetPrompt(fmt.Sprintf("C2 %s > ", CurrentAgent.ID))
		} else {
			rl.SetPrompt("C2 > ")
		}

		switch {
		case command == "exit":
			{
				// do a thing
				fmt.Println("C2 Shutting down C2, gooodbye world...")
				os.Exit(0)
			}
		case command == "help":
			{
				fmt.Print(helpMenu)
			}
		case command == "clear":
			{
				if err := clearConsole(); err != nil {
					fmt.Println(err)
				}
			}
		case command == "agents":
			{
				fmt.Printf("%10s %10s %10s\n", "ID", "IP", "Last Call")
				for _, agent := range AgentMap.Agents {
					fmt.Printf("%10s %10s %5.0f seconds ago\n", agent.ID, agent.IP, time.Since(agent.LastCall).Seconds())
				}
			}
		case command == "history":
			{
				// print everything from history file
				readHistoryFile()
			}
		case strings.HasPrefix(command, "agent"):
			{
				agentId := strings.TrimPrefix(command, "agent ")
				if agent := AgentMap.Get(agentId); agent != nil {
					CurrentAgent = agent
					// TODO: Log Agent call backs to file -- new file -- read since last login
					rl.SetPrompt(fmt.Sprintf("C2 %s > ", CurrentAgent.ID))
				} else {
					fmt.Printf("KYS -- Agent %s does not exist\n", agentId)
				}
			}
		case strings.HasPrefix(command, "exec"):
			{
				if CurrentAgent == nil {
					// TODO add log entry
					fmt.Println("No agent selected!")
					continue
				}

				cmd := strings.TrimPrefix(command, "exec ")
				fullCmd := strings.Split(cmd, " ")

				AgentMap.Enqueue(CurrentAgent.ID, fullCmd)
			}
		}

	}

}
