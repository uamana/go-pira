package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go-pira/pkg/pira"
)

func main() {
	client, err := pira.Dial("/dev/tty.usbserial-A8ATQQ5Y", 115_200, 500*time.Millisecond)
	if err != nil {
		fmt.Println("Error dialing Pira:", err)
		os.Exit(1)
	}
	defer client.Close()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	basicData, err := client.GetBasicData()
	if err != nil {
		fmt.Println("Error getting basic data:", err)
		os.Exit(1)
	}
	jsonData, err := json.Marshal(basicData)
	if err != nil {
		fmt.Println("Error marshalling basic data:", err)
		os.Exit(1)
	}
	fmt.Println("Basic data:", string(jsonData))

	memory := pira.MemoryPart1{}
	err = client.Load(0x1A, &memory)
	if err != nil {
		fmt.Println("Error getting memory:", err)
		os.Exit(1)
	}
	jsonData, err = json.Marshal(memory)
	if err != nil {
		fmt.Println("Error marshalling memory:", err)
		os.Exit(1)
	}
	fmt.Println("Memory:", string(jsonData))

	memory2 := pira.MemoryPart2{}
	err = client.Load(0x48C, &memory2)
	if err != nil {
		fmt.Println("Error getting memory2:", err)
		os.Exit(1)
	}
	jsonData, err = json.Marshal(memory2)
	if err != nil {
		fmt.Println("Error marshalling memory2:", err)
		os.Exit(1)
	}
	fmt.Println("Memory2:", string(jsonData))

	os.Exit(0)
}
