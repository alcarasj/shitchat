package main

import (
    "flag"
    "fmt"
    "os"
    "net/http"
    "bufio"
    "log"
    "github.com/gin-gonic/gin"
)

// Too lazy to set up a central domain and deal with firewalls, so this'll do for now.
const SERVER_IP = "127.0.0.1:8080"
// Too lazy to dump client settings file, so this'll do for now.
var username string

func displayMessage(c *gin.Context) {
    message := c.PostForm("message")
    sender := c.PostForm("from")
    fmt.Println("[%s] %s", sender, message)
}

func sendMessage(recipient string, message string) {
    _, err := http.Get(SERVER_IP + "/users/" + recipient)
    if err != nil {
        log.Fatal(err)
    }
}

func startUI() {
    var recipient string
    reader := bufio.NewReader(os.Stdin)
    for {
        var name string
        fmt.Println("Welcome to ShitChat.")
        fmt.Print(">")
        input, _ := reader.ReadString('\n')
        if input == "!exit" {
            fmt.Println("Exiting...")
            break
        } else if input == "!users" {

        }
    }
}

func main() {
    var port string
    flag.StringVar(&port, "port", "3000", "Port number to use. Defaults to 3000.")
    flag.StringVar(&username, "username", "", "Your username.")

    if port == "" {
        log.Fatal("Port must be set.")
    } else if username == "" {
        log.Fatal("Username must be set.")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.POST("/inbox", displayMessage)
    router.Run(":" + port)
    startUI()
}