package main

import (
	"fmt"
	"log"
	"os"
	"poker"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, closeFile, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	defer closeFile()

	if err != nil {
		log.Fatal(err)
	}

	poker.NewCLI(os.Stdin, os.Stdout, poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)).PlayPoker()
}
