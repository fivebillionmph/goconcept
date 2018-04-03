package goconcept

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func serverCommands(server *Server) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		serverCommands__process(server, text)
	}
}

func serverCommands__process(server *Server, cmd string) {
	cmd_split1 := strings.Split(cmd, " ")
	var cmd_split2 []string
	for _, s := range cmd_split1 {
		cmd_split2 = append(cmd_split2, strings.Trim(s, " \t\n"))
	}
	switch cmd_split2[0] {
		case "useradd":
			serverCommands_useradd(server, cmd_split2)
		default:
			serverCommands_invalid(server, cmd_split2)
	}
}

func serverCommands_useradd(server *Server, cmds []string) {
	if len(cmds) < 5 {
		fmt.Printf("how to use: useradd email username password admin{1 or 0}\n")
		return
	}
	email := cmds[1]
	username := cmds[2]
	password := cmds[3]
	admin := cmds[4]
	var level uint8 = 1
	if admin == "1" {
		level = 2
	}

	_, err := DBUser__create(server.connection, email, password, username, level)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
		return
	}
	fmt.Printf("success\n")
}

func serverCommands_invalid(server *Server, cmds []string) {
	fmt.Printf("invalid command: %s\n", cmds[0])
}
