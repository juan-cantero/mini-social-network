package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// AddPerson adds a new person to the database
// Creates a node with label :Persona and properties: nombre, ciudad, hobby
func (sn *SocialNetwork) AddPerson(name, city, hobby string) error {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	// Check if person already exists (unique constraint)
	checkQuery := `MATCH (p:Persona {nombre: $nombre}) RETURN p`
	result, err := session.Run(sn.ctx, checkQuery, map[string]interface{}{
		"nombre": name,
	})
	if err != nil {
		return fmt.Errorf("error checking if person exists: %w", err)
	}

	if result.Next(sn.ctx) {
		return fmt.Errorf("person '%s' already exists", name)
	}

	// Create the person
	query := `CREATE (:Persona {nombre: $nombre, ciudad: $ciudad, hobby: $hobby})`
	_, err = session.Run(sn.ctx, query, map[string]interface{}{
		"nombre": name,
		"ciudad": city,
		"hobby":  hobby,
	})

	if err != nil {
		return fmt.Errorf("error creating person: %w", err)
	}

	return nil
}

// ListAllPeople returns all people in the database
func (sn *SocialNetwork) ListAllPeople() ([]Person, error) {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	query := `MATCH (p:Persona) RETURN p.nombre AS nombre, p.ciudad AS ciudad, p.hobby AS hobby ORDER BY p.nombre`
	result, err := session.Run(sn.ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error listing people: %w", err)
	}

	var people []Person
	for result.Next(sn.ctx) {
		record := result.Record()
		nombre, _ := record.Get("nombre")
		ciudad, _ := record.Get("ciudad")
		hobby, _ := record.Get("hobby")

		people = append(people, Person{
			Name:  nombre.(string),
			City:  ciudad.(string),
			Hobby: hobby.(string),
		})
	}

	return people, nil
}

// SearchPerson searches for a person by name
func (sn *SocialNetwork) SearchPerson(name string) (*Person, error) {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	query := `MATCH (p:Persona {nombre: $nombre}) RETURN p.nombre AS nombre, p.ciudad AS ciudad, p.hobby AS hobby`
	result, err := session.Run(sn.ctx, query, map[string]interface{}{
		"nombre": name,
	})
	if err != nil {
		return nil, fmt.Errorf("error searching person: %w", err)
	}

	if !result.Next(sn.ctx) {
		return nil, nil // Person not found
	}

	record := result.Record()
	nombre, _ := record.Get("nombre")
	ciudad, _ := record.Get("ciudad")
	hobby, _ := record.Get("hobby")

	return &Person{
		Name:  nombre.(string),
		City:  ciudad.(string),
		Hobby: hobby.(string),
	}, nil
}

// CreateFriendship creates a bidirectional friendship relationship
// Relationship type: :AMIGO_DE
func (sn *SocialNetwork) CreateFriendship(name1, name2 string) error {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	if name1 == name2 {
		return fmt.Errorf("cannot create friendship with oneself")
	}

	// Check if both people exist
	checkQuery := `
		MATCH (a:Persona {nombre: $nombre1})
		MATCH (b:Persona {nombre: $nombre2})
		RETURN a, b
	`
	result, err := session.Run(sn.ctx, checkQuery, map[string]interface{}{
		"nombre1": name1,
		"nombre2": name2,
	})
	if err != nil {
		return fmt.Errorf("error checking if people exist: %w", err)
	}

	if !result.Next(sn.ctx) {
		return fmt.Errorf("one or both people not found")
	}

	// Check if friendship already exists
	friendshipQuery := `
		MATCH (a:Persona {nombre: $nombre1})-[:AMIGO_DE]-(b:Persona {nombre: $nombre2})
		RETURN a, b
	`
	friendResult, err := session.Run(sn.ctx, friendshipQuery, map[string]interface{}{
		"nombre1": name1,
		"nombre2": name2,
	})
	if err != nil {
		return fmt.Errorf("error checking friendship: %w", err)
	}

	if friendResult.Next(sn.ctx) {
		return fmt.Errorf("friendship already exists between '%s' and '%s'", name1, name2)
	}

	// Create bidirectional friendship
	query := `
		MATCH (a:Persona {nombre: $nombre1}), (b:Persona {nombre: $nombre2})
		CREATE (a)-[:AMIGO_DE]->(b), (b)-[:AMIGO_DE]->(a)
	`
	_, err = session.Run(sn.ctx, query, map[string]interface{}{
		"nombre1": name1,
		"nombre2": name2,
	})

	if err != nil {
		return fmt.Errorf("error creating friendship: %w", err)
	}

	return nil
}

// GetFriends returns all friends of a person
func (sn *SocialNetwork) GetFriends(name string) ([]Person, error) {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	query := `
		MATCH (p:Persona {nombre: $nombre})-[:AMIGO_DE]->(friend:Persona)
		RETURN friend.nombre AS nombre, friend.ciudad AS ciudad, friend.hobby AS hobby
		ORDER BY friend.nombre
	`
	result, err := session.Run(sn.ctx, query, map[string]interface{}{
		"nombre": name,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting friends: %w", err)
	}

	var friends []Person
	for result.Next(sn.ctx) {
		record := result.Record()
		nombre, _ := record.Get("nombre")
		ciudad, _ := record.Get("ciudad")
		hobby, _ := record.Get("hobby")

		friends = append(friends, Person{
			Name:  nombre.(string),
			City:  ciudad.(string),
			Hobby: hobby.(string),
		})
	}

	return friends, nil
}

// DeleteFriendship deletes a bidirectional friendship
func (sn *SocialNetwork) DeleteFriendship(name1, name2 string) error {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	query := `
		MATCH (a:Persona {nombre: $nombre1})-[r:AMIGO_DE]-(b:Persona {nombre: $nombre2})
		DELETE r
	`
	result, err := session.Run(sn.ctx, query, map[string]interface{}{
		"nombre1": name1,
		"nombre2": name2,
	})
	if err != nil {
		return fmt.Errorf("error deleting friendship: %w", err)
	}

	summary, err := result.Consume(sn.ctx)
	if err != nil {
		return fmt.Errorf("error consuming result: %w", err)
	}

	if summary.Counters().RelationshipsDeleted() == 0 {
		return fmt.Errorf("friendship not found between '%s' and '%s'", name1, name2)
	}

	return nil
}

// GetCityRecommendations suggests people from the same city who are not friends
func (sn *SocialNetwork) GetCityRecommendations(name string) ([]Person, error) {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	// Find people from the same city who are not already friends
	query := `
		MATCH (p:Persona {nombre: $nombre})
		MATCH (candidato:Persona {ciudad: p.ciudad})
		WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
		RETURN candidato.nombre AS nombre, candidato.ciudad AS ciudad, candidato.hobby AS hobby
		ORDER BY candidato.nombre
	`
	result, err := session.Run(sn.ctx, query, map[string]interface{}{
		"nombre": name,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting city recommendations: %w", err)
	}

	var recommendations []Person
	for result.Next(sn.ctx) {
		record := result.Record()
		nombre, _ := record.Get("nombre")
		ciudad, _ := record.Get("ciudad")
		hobby, _ := record.Get("hobby")

		recommendations = append(recommendations, Person{
			Name:  nombre.(string),
			City:  ciudad.(string),
			Hobby: hobby.(string),
		})
	}

	return recommendations, nil
}

// GetHobbyRecommendations suggests people with the same hobby who are not friends
func (sn *SocialNetwork) GetHobbyRecommendations(name string) ([]Person, error) {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	// Find people with the same hobby who are not already friends
	query := `
		MATCH (p:Persona {nombre: $nombre})
		MATCH (candidato:Persona {hobby: p.hobby})
		WHERE candidato <> p AND NOT (p)-[:AMIGO_DE]-(candidato)
		RETURN candidato.nombre AS nombre, candidato.ciudad AS ciudad, candidato.hobby AS hobby
		ORDER BY candidato.nombre
	`
	result, err := session.Run(sn.ctx, query, map[string]interface{}{
		"nombre": name,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting hobby recommendations: %w", err)
	}

	var recommendations []Person
	for result.Next(sn.ctx) {
		record := result.Record()
		nombre, _ := record.Get("nombre")
		ciudad, _ := record.Get("ciudad")
		hobby, _ := record.Get("hobby")

		recommendations = append(recommendations, Person{
			Name:  nombre.(string),
			City:  ciudad.(string),
			Hobby: hobby.(string),
		})
	}

	return recommendations, nil
}

// GetStatistics returns basic statistics about the network
func (sn *SocialNetwork) GetStatistics() (map[string]int, error) {
	session := sn.driver.NewSession(sn.ctx, neo4j.SessionConfig{})
	defer session.Close(sn.ctx)

	// Count total people
	peopleQuery := `MATCH (p:Persona) RETURN count(p) AS total`
	peopleResult, err := session.Run(sn.ctx, peopleQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("error counting people: %w", err)
	}

	var totalPeople int64
	if peopleResult.Next(sn.ctx) {
		record := peopleResult.Record()
		total, _ := record.Get("total")
		totalPeople = total.(int64)
	}

	// Count total friendships (divided by 2 because relationships are bidirectional)
	friendshipsQuery := `MATCH ()-[r:AMIGO_DE]->() RETURN count(r) AS total`
	friendshipsResult, err := session.Run(sn.ctx, friendshipsQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("error counting friendships: %w", err)
	}

	var totalFriendships int64
	if friendshipsResult.Next(sn.ctx) {
		record := friendshipsResult.Record()
		total, _ := record.Get("total")
		totalFriendships = total.(int64) / 2 // Divide by 2 for bidirectional relationships
	}

	return map[string]int{
		"people":      int(totalPeople),
		"friendships": int(totalFriendships),
	}, nil
}
