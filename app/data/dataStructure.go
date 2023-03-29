package data

import (
	"errors"
	"time"
)

type DishDelivery struct {
	ID                 string    `json:"id"`
	SubscriptionDishID string    `json:"subscriptionDishID"`
	Status             string    `json:"status"`
	ExpectedTime       time.Time `json:"expectedTime"`
	DeliveryTime       time.Time `json:"deliveryTime,omitempty"`
	Note               string    `json:"note,omitempty"`
}

type Subscription struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userID"`
	PlaylistID      string    `json:"playlistID,omitempty"`
	Customized      bool      `json:"customized"`
	Status          string    `json:"status"`
	Frequency       string    `json:"frequency"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate,omitempty"`
	ReceiverName    string    `json:"receiverName"`
	ReceiverContact string    `json:"receiverContact"`
}

type SubscriptionDish struct {
	ID             string    `json:"id"`
	DishID         string    `json:"dishID"`
	SubscriptionID string    `json:"subscriptionID"`
	ScheduleTime   time.Time `json:"scheduleTime"`
	Frequency      string    `json:"frequency"`
	DishOptions    string    `json:"dishOptions,omitempty"`
	Note           string    `json:"note,omitempty"`
}

type SubscriptionDishDTO struct {
	ID             string     `json:"id"`
	DishID         string     `json:"dishID"`
	SubscriptionID string     `json:"subscriptionID"`
	ScheduleTime   time.Time  `json:"scheduleTime"`
	Frequency      string     `json:"frequency"`
	DishOptions    [][]string `json:"dishOptions,omitempty"`
	Note           string     `json:"note,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// This struct includes all the data returned to the request
// DishIncluded is a map structure, the key is the DishID
// RestaurantInfo is a map structure, the key is the RestaurantID
// RestaurantAddress is a map structure, the key is the AddressID

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)
