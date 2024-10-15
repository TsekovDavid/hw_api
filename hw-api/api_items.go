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
	e.GET("/item/:caption", GetItems)
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

func findItems(caption string) []Item {
	var foundItems []Item
	for _, item := range Items {
		if item.Caption == caption {
			foundItems = append(foundItems, item)
		}
	}
	return foundItems
}

func GetItems(c echo.Context) error {
	items := findItems(c.Param("caption"))
	if len(items) == 0 {
		return c.String(http.StatusBadRequest, "items not found")
	}
	var response string
	for _, item := range items {
		response += fmt.Sprintf("caption: %s, weight: %f, number: %d\n", item.Caption, item.Weight, item.Number)
	}
	return c.String(http.StatusOK, response)
}
