package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Item struct {
	Caption string  `json:"caption"`
	Weight  float32 `json:"weight"`
	Number  int     `json:"number"`
}

var Items = make([]Item, 0)

func main() {
	e := echo.New()
	e.POST("/item", AddItem)
	e.GET("/item/:caption", GetItem)
	log.Fatal(e.Start(":8080"))
}

func AddItem(c echo.Context) error {
	var item Item
	if err := c.Bind(&item); err != nil {
		return c.String(http.StatusBadRequest, "Invalid data")
	}
	Items = append(Items, item)
	log.Println(Items)
	return c.String(http.StatusOK, "Item added successfully")
}

func findItem(caption string) *Item {
	for _, item := range Items {
		if item.Caption == caption {
			return &item
		}
	}
	return nil
}

func GetItem(c echo.Context) error {
	// item := findItem(c.QueryParam("caption"))
	// if item != nil {
	// 	return c.String(http.StatusOK, fmt.Sprintf("caption: %s, weight: %f, number: %d", item.Caption, item.Weight, item.Number))
	// }
	item := findItem(c.Param("caption"))
	if item != nil {
		return c.String(http.StatusOK, fmt.Sprintf("caption: %s, weight: %f, number: %d", item.Caption, item.Weight, item.Number))
	}
	return c.String(http.StatusBadRequest, "item not found")
}
