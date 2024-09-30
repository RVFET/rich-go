package main

import "github.com/rvfet/rich-go"

type Person struct {
	Name      string
	Age       int
	Contacts  map[string]any
	Languages []string
	Height    float64
	Notes     any
	IsAdmin   bool
	IsBanned  bool
}

func main() {
	rvfet := Person{
		// String
		Name: "Rafet Abbasli",
		// Integer
		Age: 23,
		// Map
		Contacts: map[string]any{
			"email":       "me@rvfet.com",
			"github":      "https://github.com/rvfet",
			"telegram":    "t.me/rvfet",
			"mobilePhone": 1234567890,
		},
		// Slice
		Languages: []string{"Go", "Python", "JavaScript"},
		// Boolean
		IsAdmin:  true,
		IsBanned: false,
		// Float
		Height: 1.90,
		// Nil
		Notes: nil,
	}

	rich.Print(rvfet)
	rich.Error("Could not connect to the server, please try again later.")
	rich.Success("The data has been saved successfully.")
	rich.Warning("Operation timed out, please try again later.")
	rich.Info("The application is running in the debug mode.")
	rich.Debug("The value of the variable is", 42)
}
