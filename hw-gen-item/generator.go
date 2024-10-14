package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var Captions = []string{"Box", "Ball", "Book", "Pen", "Phone", "Laptop", "Cup", "Notebook"}

type Item struct {
	Caption string  `json:"caption"`
	Weight  float32 `json:"weight"`
	Number  int     `json:"number"`
}

func GenerateRandomItem() (Item, error) {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)
	if len(Captions) > 0 {
		caption := Captions[0]
		Captions = Captions[1:]
		return Item{
			Caption: caption,
			Weight:  r.Float32() * 100, // случайный вес от 0 до 100
			Number:  r.Intn(100) + 1,   // случайное число от 1 до 100
		}, nil
	}
	return Item{}, fmt.Errorf("No more captions available")
}

func CreateItem(item Item, er error) error {
	if er != nil {
		return fmt.Errorf("No more captions available")
	}

	jsonData := fmt.Sprintf(`{"caption":"%s", "weight":%f, "number":%d}`, item.Caption, item.Weight, item.Number)
	resp, err := http.Post("http://localhost:8080/item", "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to send post request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add item, server returned: %s", resp.Status)
	}

	return nil
}

func GetItemInfo(caption string) (Item, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/item/%s", caption))
	if err != nil {
		return Item{}, fmt.Errorf("failed to send get request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Item{}, fmt.Errorf("failed to get item, server returned: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Item{}, fmt.Errorf("failed to read response body: %v", err)
	}
	responseString := string(body)
	parts := strings.Split(responseString, ", ")

	if len(parts) != 3 {
		return Item{}, fmt.Errorf("unexpected response format")
	}
	captionPart := strings.Split(parts[0], ": ")[1]
	weightPart := strings.Split(parts[1], ": ")[1]
	numberPart := strings.Split(parts[2], ": ")[1]
	weight, err := strconv.ParseFloat(weightPart, 32)
	if err != nil {
		return Item{}, fmt.Errorf("invalid weight format")
	}

	number, err := strconv.Atoi(numberPart)
	if err != nil {
		return Item{}, fmt.Errorf("invalid number format")
	}

	item := Item{
		Caption: captionPart,
		Weight:  float32(weight),
		Number:  number,
	}

	return item, nil
}

func main() {
	fmt.Println("Сколько создать предметов?")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	count, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка: введено некорректное число")
		return
	}

	var itemsCap []string
	for i := 0; i < count; i++ {
		item, err := GenerateRandomItem()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		itemsCap = append(itemsCap, item.Caption)
		err = CreateItem(item, nil)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
	for _, cap := range itemsCap {
		item, err := GetItemInfo(cap)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		sumWeight := item.Weight * float32(item.Number)
		fmt.Printf("%s %d шт общим весом %.2f\n", cap, item.Number, sumWeight)
	}
	fmt.Println("Создание предметов завершено.")
}
