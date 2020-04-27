package main

import (
    "flag"
    "fmt"
    "os"
    "net/http"
    "bufio"
    "log"
    "net"
    "encoding/json"
)

// Too lazy to set up a central domain and deal with firewalls, so this'll do for now.
const DIRECTORY_SERVICE_IP = "127.0.0.1:8080"

type Message struct {
    senderName       string
    recipientName    string
    text             string
}

func sendMessage(sender string, recipient string, text string) {
    // message := Message{ senderName: sender, recipientName: recipient, text: text }
    response, error := http.Get(DIRECTORY_SERVICE_IP + "/users/" + recipient)
    if error != nil {
        log.Fatal(error)
    }
    var result map[string]interface{}
    json.NewDecoder(response.Body).Decode(&result)
    // TO-DO
}

func main() {
    var username string
    var port string
    flag.StringVar(&port, "port", "3000", "Port number to use. Defaults to 3000.")
    flag.StringVar(&username, "username", "", "Your username.")
    flag.Parse()

    if port == "" {
        log.Fatal("Port must be set.")
    } else if username == "" {
        log.Fatal("Username must be set.")
    }

    fmt.Println("Welcome to ShitChat!")
    fmt.Println("Your username is " + username + ".")
    listener, _ := net.Listen("tcp", ":" + port)
    connection, _ := listener.Accept()
    reader := bufio.NewReader(os.Stdin)

    for {
        message, _ := bufio.NewReader(connection).ReadString('\n')
        fmt.Println("Message received: %s", string(message))

        fmt.Println(">")
        input, _ := reader.ReadString('\n')
        if input == "!exit" {
            fmt.Println("Exiting...")
            break
        }
        // TO-DO
    }
}