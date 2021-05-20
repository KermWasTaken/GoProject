package main

import (
	"net"
	"log"
	"os"
	"fmt"
	"strings"
	"bufio"
	"io"
)

//Function called if error occurs
func logFatal(err error){
	if err != nil{
		log.Fatal(err) //Calls os.exit(1) which forces the program to terminate
	}
}

func main() {
	connection, err := net.Dial("tcp", "localhost:8000")
	logFatal(err) //Checks for error

	fmt.Println("Enter your username: ")

	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n') //Makes user enter in username first

	logFatal(err) //Checks for error

	username = strings.Trim(username, " \r\n") //Stores entered username into the username variable

	enterMsg := fmt.Sprintf("%s has joined the chatroom. \n", username)

	fmt.Println(enterMsg)

	go read(connection)
	write(connection, username)

	defer connection.Close()
}

func write(connection net.Conn, username string){
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')

		if (err != nil){ //error check
			break
		}

		message = fmt.Sprintf("%s: %s\n", username, strings.Trim(message, " \r\n")) //Formats the string so that the sender's username appears before message

		connection.Write([]byte(message))
	}
}

func read(connection net.Conn){
	for {
		reader := bufio.NewReader(connection)
		message, err := reader.ReadString('\n')

		if (err == io.EOF){ //error check
			connection.Close()
			fmt.Println("Error! Connection has been closed.")
			os.Exit(1)
		}

		fmt.Println(message)
		fmt.Println("")
	}
}