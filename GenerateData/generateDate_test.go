package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type SubServiceRequestDataDTO struct {
	SubscriptionRequest SubRequested
	DishIncluded        []SubDishRequested
}

type SubRequested struct {
	UserID          string    `json:"userID"`
	PlaylistID      string    `json:"playlistID"`
	Customized      bool      `json:"customized"`
	Frequency       string    `json:"frequency"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate,omitempty"`
	ReceiverName    string    `json:"receiverName"`
	ReceiverContact string    `json:"receiverContact"`
}

type SubDishRequested struct {
	DishID       string     `json:"dishID"`
	ScheduleTime time.Time  `json:"scheduleTime"`
	Frequency    string     `json:"frequency"`
	DishOptions  [][]string `json:"dishOptions`
	Note         string     `json:"Note,omitempty"`
}

func TestGenerateNewSubscriptionDTO(t *testing.T) {
	infoIDs, err := extractIDs()

	if err != nil {
		log.Fatal("could not extract ID information.")
	}

	frequencyChoices := []string{
		"Daily",
		"Weekly",
	}

	customizedChoice := []bool{
		true,
		false,
	}

	subReq := SubRequested{
		UserID:          "user6",
		PlaylistID:      infoIDs.playlists[rand.Intn((len(infoIDs.playlists)))],
		Customized:      customizedChoice[rand.Intn(len(customizedChoice))],
		Frequency:       frequencyChoices[rand.Intn(len(frequencyChoices))],
		StartDate:       time.Now().AddDate(0, 0, -rand.Intn(15)),
		EndDate:         time.Now().AddDate(0, 0, 7*rand.Intn(3)),
		ReceiverName:    "Tony Lew",
		ReceiverContact: "12348766",
	}

	dishes := []SubDishRequested{}

	if subReq.Customized == true {
		n := len(infoIDs.dishes)
		for i := 0; i < rand.Intn(5); i++ {
			dish := SubDishRequested{
				DishID:       infoIDs.dishes[rand.Intn(n)],
				ScheduleTime: subReq.StartDate.AddDate(0, 0, rand.Intn(3)),
				Frequency:    subReq.Frequency,
				DishOptions: [][]string{
					{"Spice", "Yes"},
					{"Wasabi", "No"},
				},
				Note: "test dish " + strconv.Itoa(i),
			}
			dishes = append(dishes, dish)

		}

	} else {

		dishIDs := infoIDs.playlistDishRelation[subReq.PlaylistID]

		for _, dishID := range dishIDs {
			dish := SubDishRequested{
				DishID:       dishID,
				ScheduleTime: subReq.StartDate.AddDate(0, 0, rand.Intn(3)),
				Frequency:    subReq.Frequency,
				DishOptions: [][]string{
					{"Spice", "No"},
					{"Wasabi", "No"},
				},
				Note: "playlist dish",
			}
			dishes = append(dishes, dish)

		}

	}

	subReqServiceDTO := SubServiceRequestDataDTO{
		SubscriptionRequest: subReq,
		DishIncluded:        dishes,
	}

	jsonResult, _ := json.MarshalIndent(subReqServiceDTO, "", " ")

	fmt.Println("================Generated Subscription Data for Posting====================")
	fmt.Println(string(jsonResult))
	fmt.Println("===========================================================================")
}
