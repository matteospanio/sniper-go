package main

import (
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", home)
	router.GET("/snipe", snipe)

	router.Run("0.0.0.0:8080")
}

func home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

func snipe(c *gin.Context) {

	cCp := c.Copy()

	go func() {

		var command = fmt.Sprintf("sniper -t %s", cCp.Query("target"))
		cmd, err := exec.Command("bash", "-c", command).Output()

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"output": string(cmd)})
		}
	}()
}
