package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Person represents a person in the social network
type Person struct {
	Name  string
	City  string
	Hobby string
}

// SocialNetwork handles all operations for the mini social network
type SocialNetwork struct {
	driver neo4j.DriverWithContext
	ctx    context.Context
}

// NewSocialNetwork creates a new social network instance
func NewSocialNetwork(uri, username, password string) (*SocialNetwork, error) {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	// Verify connectivity
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &SocialNetwork{
		driver: driver,
		ctx:    ctx,
	}, nil
}

// Close closes the database connection
func (sn *SocialNetwork) Close() error {
	return sn.driver.Close(sn.ctx)
}
