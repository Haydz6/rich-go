package main

import (
	"fmt"
	"time"

	"github.com/Haydz6/rich-go/client"
)

func main() {
	client.AuthenticationUpdate.Add(1)
	go func() {
		for {
			client.AuthenticationUpdate.Wait()
			client.AuthenticationUpdate.Add(1)
		}
	}()

	err := client.Login("1105722413905346660")
	if err != nil {
		panic(err)
	}

	now := time.Now()

	for {
		err = client.SetActivity(&client.Activity{
			State:      "Heyy!!!",
			Details:    "I'm running on rich-go :)",
			LargeImage: "largeimageid",
			LargeText:  "This is the large image :D",
			SmallImage: "smallimageid",
			SmallText:  "And this is the small image",
			Party: &client.Party{
				ID:         "-1",
				Players:    15,
				MaxPlayers: 24,
			},
			Timestamps: &client.Timestamps{
				Start: &now,
			},
			Buttons: []*client.Button{
				&client.Button{
					Label: "GitHub",
					Url:   "https://github.com/Haydz6/rich-go",
				},
			},
		})

		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 5)
	}

	// Discord will only show the presence if the app is running
	// Sleep for a few seconds to see the update
	fmt.Println("Sleeping...")
	time.Sleep(time.Second * 200)
}
