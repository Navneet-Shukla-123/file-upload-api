package main

import (
	"file-api/routes"

	"github.com/gin-gonic/gin"
)


func main(){
	router:=gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/",routes.Home)
	//router.POST("/upload",routes.Upload)
	router.POST("/upload",routes.UploadInChunk)
	router.Run()
}