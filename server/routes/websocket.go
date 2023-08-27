package routes

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/matteospanio/sniper-go/sniper"
)

var (
	Wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	connections = make(map[*websocket.Conn]bool)
)

var mutex sync.Mutex

func printer(status *Status, start string, target string) {
	for status.Running {
		read := exec.Command("bash", "-c", "tail -n 200 "+sniper.OutPath+"/sniper-"+target+"-"+start+".txt")

		out, err := read.Output()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		for conn := range connections {
			mutex.Lock()
			conn.WriteMessage(websocket.TextMessage, out)
			mutex.Unlock()
		}

		time.Sleep(1 * time.Second)
	}
}

func HandleWebSocket(c *gin.Context) {
	conn, err := Wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	defer func() {
		conn.Close()
		delete(connections, conn)
	}()

	connections[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		target := string(msg)

		cmd := exec.Command("bash", "-c", "sniper -t "+target)
		fmt.Println("Running: ", cmd.String())

		status := Status{Running: true}

		date := time.Now().Format("2006-01-02-15-04")
		date = strings.Replace(date, "-", "", -1)

		cmd.Start()
		go printer(&status, date, target)

		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
			continue
		}

		status.Running = false

		for conn := range connections {
			mutex.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, []byte("[DONE]")); err != nil {
				fmt.Println("Failed to send message to client:", err)
			}
			mutex.Unlock()
		}

		fmt.Println("[DONE]")
	}
}
