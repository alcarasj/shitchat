package main

import (
    "flag"
    "time"
    "fmt"
    "regexp"
    "net/http"
    "strconv"
    "log"
    "github.com/gin-gonic/gin"
)

// Too lazy to set up auth, user expiry will do for now.
const USER_TTL_SECONDS = 86400

type User struct {
    ip          string
    port        uint64
    lastSeen    int64
}

// Too lazy to set up a real database, so this'll do for now.
var directory = make(map[string]User)

func keepUserAlive(username string) {
    if user, userExists := directory[username]; userExists {
        user.lastSeen = time.Now().Unix()
        directory[username] = user
    }
}

func createUser(c *gin.Context) {
    username, valueExists := c.GetPostForm("username")
    if !valueExists {
        c.JSON(http.StatusConflict, gin.H{ "message": "Username field is required." })
        return
    }

    re := regexp.MustCompile("^[a-zA-Z0-9_]*$z")

    if _, userExists := directory[username]; userExists {
        c.JSON(http.StatusConflict, gin.H{ 
            "message": "This username already exists. Please provide another username.",
        })
    } else if usernameIsValid := re.MatchString(username); !usernameIsValid && username != "" {
        c.JSON(http.StatusConflict, gin.H{ 
            "message": "Username must only contain alphanumeric characters.",
        })
    } else {
        clientPortStr, valueExists := c.GetPostForm("port")
        if clientPort, error := strconv.ParseUint(clientPortStr, 0, 16); valueExists && error == nil {
            clientIP := c.ClientIP()
            user := User{ ip: clientIP, lastSeen: time.Now().Unix(), port: clientPort }
            directory[username] = user
            fmt.Println("%s registered with username %s.", clientIP, username)
            c.JSON(http.StatusOK, gin.H{ "message": "You have successfully registered as " + username + "." })
        } else {
            c.JSON(http.StatusConflict, gin.H{ 
                "message": "Port field is required and must be a positive integer.",
            })
        }
    }
}

func getUsers(c *gin.Context) {
    usernames := make([]string, 0, len(directory))
    now := time.Now().Unix()
    for username, user := range directory {
        if (now - user.lastSeen) < USER_TTL_SECONDS {
            usernames = append(usernames, username)
        } else {
            delete(directory, username)
        }
    }
    c.JSON(http.StatusOK, gin.H{ "users": usernames })
}

func getUser(c *gin.Context) {
    username := c.Param("username")
    if user, userExists := directory[username]; userExists {
        c.JSON(http.StatusOK, gin.H{ "ip": user.ip, "port": user.port })
    } else {
        c.JSON(http.StatusNotFound, gin.H{ "message": username + " is not a registered user." })
    }
}

func updateUser(c *gin.Context) {
    username := c.Param("username")
    if _, userExists := directory[username]; userExists {
        keepUserAlive(username)
        c.Status(http.StatusOK)
    } else {
        c.JSON(http.StatusNotFound, gin.H{ "message": username + " is not a registered user." })
    }
}

func main() {
    var port string
    flag.StringVar(&port, "port", "8080", "Port number to use. Defaults to 8080.")
    flag.Parse()

    if port == "" {
        log.Fatal("Port must be set.")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.POST("/users", createUser)
    router.GET("/users", getUsers)
    router.PUT("/users/:username", updateUser)
    router.GET("/users/:username", getUser)
    router.Run(":" + port)
}