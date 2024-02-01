package webapi

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"worksp/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// IP whitelist to give access to specific IPs
// add IP to whitelist to access data
var IPWhitelist = map[string]bool{
	"127.0.0.1":     true,
	"192.168.0.100": true,
}

// runs webapi; capital letter to be public
func InitializeWebApi() {
	router := gin.Default()

	// apply CORS
	config := cors.DefaultConfig()
	// allows origins from clients specofoed in slice -> add to list is new services want to communicate
	config.AllowOrigins = []string{"http://192.168.0.100:8201", "http://localhost:8201"}
	router.Use(cors.New(config))

	router.GET("/api/v1/detection", getDios)

	// security; will not show API values if the IP isn't whitelisted
	restrictedPage := router.Group("/")
	restrictedPage.Use(middleware.IPWhiteList(IPWhitelist))

	// 0.0.0.0:6001
	err := router.Run(":6001")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

// dios json struct
type Dios struct {
	ID       int64  `json:"id"`
	Duration string `json:"duration"`
	WhatTime string `json:"whattime"`
}

// fetches dios in Dios format
func fetchDios() ([]Dios, error) {
	var db *sql.DB
	var err error
	var dios []Dios

	db, err = sql.Open("mysql", "root:TDC_arch2023@tcp(192.168.0.100:3306)/diobase")
	if err != nil {
		return nil, fmt.Errorf("dios %q", err)
	}
	defer db.Close()

	// fetches all dios
	rows, err := db.Query("SELECT * FROM dios;")
	if err != nil {
		return nil, fmt.Errorf("dios %q", err)
	}
	defer rows.Close()

	// looping through rows; using scan to assign column data to struct fields
	for rows.Next() {
		var di Dios
		if err := rows.Scan(&di.ID, &di.Duration, &di.WhatTime); err != nil {
			return nil, fmt.Errorf("dios %q", err)
		}
		dios = append(dios, di)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dios %q", err)
	}

	return dios, nil
}

// get all dios from database
func getDios(c *gin.Context) {
	dios, err := fetchDios()
	if err != nil {
		log.Fatal(err)
	}

	// show in JSON format
	c.JSON(http.StatusOK, dios)
}
