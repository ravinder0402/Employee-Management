package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Data struct {
	EID          uint      `json:"id" gorm:"primary_key"`
	ETitle       string    `json:"title" binding:"required"`
	EDescription string    `json:"description"`
	EDueDate     time.Time `json:"due_date" binding:"required"`
	EStatus      string    `json:"status" gorm:"default:'pending'"`
}

var d *gorm.DB

func main() {

	initDB()

	r := gin.Default()

	r.POST("/data", createTask)
	r.GET("/data/:id", getTask)
	r.PUT("/data/:id", updateTask)
	r.DELETE("/data/:id", deleteTask)
	r.GET("/data", listTasks)

	r.Run(":8888")
}

func initDB() {
	var err error
	d, err = gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("Not Connectd To Database: " + err.Error())
	}

	d.AutoMigrate(&Data{})
}

func createTask(c *gin.Context) {
	var data Data
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.EStatus == "" {
		data.EStatus = "pending"
	}

	d.Create(&data)

	c.JSON(http.StatusCreated, data)
}

func getTask(c *gin.Context) {
	var data Data
	id := c.Param("id")

	if err := d.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")

	var cdata Data
	if err := d.First(&cdata, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var ndata Data
	if err := c.BindJSON(&ndata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d.Model(&cdata).Updates(ndata)

	c.JSON(http.StatusOK, cdata)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	var data Data
	if err := d.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	d.Delete(&data)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func listTasks(c *gin.Context) {
	var tasks []Data
	d.Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}
