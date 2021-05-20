package main

import (
	"log"
	"net"
	"bufio"
)

//Global variables
var (
	openConnections = make(map[net.Conn]bool) //Creates a map of all active connections with their values set to true
	newConnection = make(chan net.Conn) //Whenever there is a new connection made it is past through the channel
	deadConnection = make(chan net.Conn) //Whenever there is an ended connection  it is past through the channel
)

//Function called if error occurs
func logFatal(err error){
	if err != nil{
		log.Fatal(err) //Calls os.exit(1) which forces the program to terminate
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8000") //Listen takes a network "tcp" and address "port 8080"
	logFatal(err) //Checks to see if there is error
	go func() {
		//For loop that runs forever in able to always process when there are new clients
		for {
			conn, err := ln.Accept() //Accepts new clients
			logFatal(err) //Checks for errors
			openConnections[conn] = true //Adds client to active connections
			newConnection <- conn //Passes interface into newConnection so conn can be accessed outside of goroutine
		}
	}()

	for {
		select {
			case conn := <-newConnection: //If this case, broadcastMessage sends message to all clients
				go broadcastMessage(conn)
			case conn := <-deadConnection:
				for item := range openConnections {
					if item == conn {
						break
					}
				}
				delete(openConnections, conn) //If connection is found in active maps, breaks out of case, otherwise connection is deleted
		}
	}
	defer ln.Close()
}

func broadcastMessage(conn net.Conn){
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err != nil { //Breaks if there is an error
			break
		}

		for item := range openConnections {
			if item != conn {
				item.Write([]byte(message))
			}
		}
	}

	deadConnection <- conn
}