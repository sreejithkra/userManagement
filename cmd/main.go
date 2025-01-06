package main

import (
	"userManagement/internal/config"
	"userManagement/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	env := config.EnvConfig()
	database.ConnectDatabase(env)

	router.Run(":8086")

}
