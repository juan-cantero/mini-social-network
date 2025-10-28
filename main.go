package main

import (
	"fmt"
	"os"
)

func main() {
	// Connect to Neo4j
	sn, err := NewSocialNetwork("bolt://localhost:7687", "neo4j", "password123")
	if err != nil {
		fmt.Printf("Error connecting to Neo4j: %v\n", err)
		os.Exit(1)
	}
	defer sn.Close()

	fmt.Println("Connected to Neo4j successfully!")

	// Run the menu
	runMenu(sn)
}
