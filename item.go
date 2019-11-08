package main

import "fmt"

type Item struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (Item) ResourceType() string {
	return "item"
}

func (x Item) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("no name given")
	}
	return nil
}
