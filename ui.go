package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// runMenu displays the main menu and handles user interaction
func runMenu(sn *SocialNetwork) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== MINI SOCIAL NETWORK ===")
		fmt.Println("1. Add person")
		fmt.Println("2. List all people")
		fmt.Println("3. Search person")
		fmt.Println("4. Create friendship")
		fmt.Println("5. View friends of a person")
		fmt.Println("6. Delete friendship")
		fmt.Println("7. City-based recommendations")
		fmt.Println("8. Recommendations by hobby")
		fmt.Println("9. Statistics")
		fmt.Println("0. Exit")
		fmt.Print("\nChoose an option: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			addPerson(sn, reader)
		case "2":
			listAllPeople(sn)
		case "3":
			searchPerson(sn, reader)
		case "4":
			createFriendship(sn, reader)
		case "5":
			viewFriends(sn, reader)
		case "6":
			deleteFriendship(sn, reader)
		case "7":
			cityRecommendations(sn, reader)
		case "8":
			hobbyRecommendations(sn, reader)
		case "9":
			showStatistics(sn)
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func addPerson(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter city: ")
	city, _ := reader.ReadString('\n')
	city = strings.TrimSpace(city)

	fmt.Print("Enter hobby: ")
	hobby, _ := reader.ReadString('\n')
	hobby = strings.TrimSpace(hobby)

	if name == "" || city == "" || hobby == "" {
		fmt.Println("Error: All fields are required!")
		return
	}

	err := sn.AddPerson(name, city, hobby)
	if err != nil {
		fmt.Printf("Error adding person: %v\n", err)
		return
	}

	fmt.Printf("Person '%s' added successfully!\n", name)
}

func listAllPeople(sn *SocialNetwork) {
	people, err := sn.ListAllPeople()
	if err != nil {
		fmt.Printf("Error listing people: %v\n", err)
		return
	}

	if len(people) == 0 {
		fmt.Println("No people found in the network.")
		return
	}

	fmt.Println("\n--- All People ---")
	for _, p := range people {
		fmt.Printf("- %s (City: %s, Hobby: %s)\n", p.Name, p.City, p.Hobby)
	}
}

func searchPerson(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter name to search: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	person, err := sn.SearchPerson(name)
	if err != nil {
		fmt.Printf("Error searching person: %v\n", err)
		return
	}

	if person == nil {
		fmt.Printf("Person '%s' not found.\n", name)
		return
	}

	fmt.Printf("\nFound: %s (City: %s, Hobby: %s)\n", person.Name, person.City, person.Hobby)
}

func createFriendship(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter first person's name: ")
	name1, _ := reader.ReadString('\n')
	name1 = strings.TrimSpace(name1)

	fmt.Print("Enter second person's name: ")
	name2, _ := reader.ReadString('\n')
	name2 = strings.TrimSpace(name2)

	err := sn.CreateFriendship(name1, name2)
	if err != nil {
		fmt.Printf("Error creating friendship: %v\n", err)
		return
	}

	fmt.Printf("Friendship created between '%s' and '%s'!\n", name1, name2)
}

func viewFriends(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter person's name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	friends, err := sn.GetFriends(name)
	if err != nil {
		fmt.Printf("Error getting friends: %v\n", err)
		return
	}

	if len(friends) == 0 {
		fmt.Printf("%s has no friends yet.\n", name)
		return
	}

	fmt.Printf("\n--- Friends of %s ---\n", name)
	for _, f := range friends {
		fmt.Printf("- %s (City: %s, Hobby: %s)\n", f.Name, f.City, f.Hobby)
	}
}

func deleteFriendship(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter first person's name: ")
	name1, _ := reader.ReadString('\n')
	name1 = strings.TrimSpace(name1)

	fmt.Print("Enter second person's name: ")
	name2, _ := reader.ReadString('\n')
	name2 = strings.TrimSpace(name2)

	err := sn.DeleteFriendship(name1, name2)
	if err != nil {
		fmt.Printf("Error deleting friendship: %v\n", err)
		return
	}

	fmt.Printf("Friendship deleted between '%s' and '%s'!\n", name1, name2)
}

func cityRecommendations(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter person's name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	recommendations, err := sn.GetCityRecommendations(name)
	if err != nil {
		fmt.Printf("Error getting recommendations: %v\n", err)
		return
	}

	if len(recommendations) == 0 {
		fmt.Printf("No city-based recommendations for %s.\n", name)
		return
	}

	fmt.Printf("\n--- City-based Recommendations for %s ---\n", name)
	for _, r := range recommendations {
		fmt.Printf("- %s (City: %s, Hobby: %s)\n", r.Name, r.City, r.Hobby)
	}
}

func hobbyRecommendations(sn *SocialNetwork, reader *bufio.Reader) {
	fmt.Print("Enter person's name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	recommendations, err := sn.GetHobbyRecommendations(name)
	if err != nil {
		fmt.Printf("Error getting recommendations: %v\n", err)
		return
	}

	if len(recommendations) == 0 {
		fmt.Printf("No hobby-based recommendations for %s.\n", name)
		return
	}

	fmt.Printf("\n--- Hobby-based Recommendations for %s ---\n", name)
	for _, r := range recommendations {
		fmt.Printf("- %s (City: %s, Hobby: %s)\n", r.Name, r.City, r.Hobby)
	}
}

func showStatistics(sn *SocialNetwork) {
	stats, err := sn.GetStatistics()
	if err != nil {
		fmt.Printf("Error getting statistics: %v\n", err)
		return
	}

	fmt.Println("\n=== STATISTICS ===")
	fmt.Printf("Total People: %d\n", stats["people"])
	fmt.Printf("Total Friendships: %d\n", stats["friendships"])
}
