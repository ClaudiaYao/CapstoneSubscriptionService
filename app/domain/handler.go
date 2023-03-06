package domain

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ClaudiaYao/CapstoneSubscriptionService/app/data"
	"github.com/go-chi/chi"
	"github.com/lithammer/shortuuid"
)

type SubscriptionService struct {
	DBConnection *data.DataQuery
	// Model        *SubscriptionServiceResponseDataDTO
}

// this SubscriptionServiceDataDTO represents the data returned to the client
type SubscriptionServiceResponseDataDTO struct {
	Subscription data.Subscription
	DishIncluded []data.SubscriptionDish
	DishDelivery []data.DishDelivery
}

type SubscriptionServiceRequestDataDTO struct {
	SubscriptionRequest SubscriptionRequested
	DishIncluded        []SubscriptionDishRequested
}

type SubscriptionRequested struct {
	UserID     string    `json:"userID"`
	PlaylistID string    `json:"playlistID"`
	Customized bool      `json:"customized"`
	Frequency  string    `json:"frequency"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
}

type SubscriptionDishRequested struct {
	DishID       string    `json:"dishID"`
	ScheduleTime time.Time `json:"scheduleTime"`
	Frequency    string    `json:"frequency"`
	Note         string    `json:"Note"`
}

// C: this PlaylistService is responsible for transfering information request/response
// C: the database operation is conducted by its member *sql.DB
// C: when designing API or micro-service, the service request passes data via JSON
func (service *SubscriptionService) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// how to define the structure of SubscriptionServiceData is depending on the front end
	var requestPayload SubscriptionServiceRequestDataDTO

	err := service.readJSON(w, r, &requestPayload)
	if err != nil {
		service.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	subscriptionReq := requestPayload.SubscriptionRequest

	subscriptionInfo := data.Subscription{
		ID:         "Sub" + shortuuid.New(),
		UserID:     subscriptionReq.UserID,
		PlaylistID: subscriptionReq.PlaylistID,
		Customized: subscriptionReq.Customized,
		Status:     "Active",
		Frequency:  subscriptionReq.Frequency,
		StartDate:  subscriptionReq.StartDate,
		EndDate:    subscriptionReq.EndDate,
	}

	dishIncluded := requestPayload.DishIncluded
	dishes := []data.SubscriptionDish{}

	for _, dishInfo := range dishIncluded {
		dish := data.SubscriptionDish{
			ID:             "SDish" + shortuuid.New(),
			DishID:         dishInfo.DishID,
			SubscriptionID: subscriptionInfo.ID,
			ScheduleTime:   dishInfo.ScheduleTime,
			Frequency:      dishInfo.Frequency,
			Note:           dishInfo.Note,
		}
		dishes = append(dishes, dish)
	}

	subscriptionID, err := service.GenerateNewSubscription(r.Context(), subscriptionInfo, dishes)

	if err != nil {
		service.errorJSON(w, errors.New("invalid query"), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("subscription is created: %s", subscriptionID),
	}

	service.writeJSON(w, http.StatusAccepted, responsePayload)
}

func (app *SubscriptionService) Welcome(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusAccepted, "Welcome to Subscription service!")
}

func (app *SubscriptionService) GetDishBySubscriptionID(w http.ResponseWriter, r *http.Request) {

	subscriptionId := chi.URLParam(r, "subscription_id")
	subscriptionDishes, err := app.DBConnection.GetDishBySubscriptionID(r.Context(), subscriptionId)
	if err != nil {
		app.errorJSON(w, errors.New(fmt.Sprint("invalid query", err)), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "dishes are retrieved",
		Data:    subscriptionDishes,
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)

}

func (service *SubscriptionService) GetDishDeliveryStatus(w http.ResponseWriter, r *http.Request) {

	dishID := chi.URLParam(r, "dish_id")
	dishDeliveryStatus, err := service.DBConnection.GetDishBySubscriptionID(r.Context(), dishID)
	if err != nil {
		service.errorJSON(w, errors.New("invalid query"), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "dishes are retrieved",
		Data:    dishDeliveryStatus,
	}

	// C: this means the success response
	service.writeJSON(w, http.StatusAccepted, responsePayload)

}

func (service *SubscriptionService) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {

	// planType := entities.SourceB2C
	// source := strings.TrimSpace(r.URL.Query().Get("source"))

	id := chi.URLParam(r, "id")
	subscription, err := service.DBConnection.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		service.errorJSON(w, errors.New("invalid query"), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "subscription are retrieved",
		Data:    subscription,
	}

	service.writeJSON(w, http.StatusAccepted, responsePayload)

}

func (service *SubscriptionService) GetSubscriptionByUserID(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "user_id")

	subscriptions, err := service.DBConnection.GetSubscriptionByUserID(r.Context(), userID)
	if err != nil {
		service.errorJSON(w, errors.New("invalid query"), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "subscriptions are retrieved",
		Data:    subscriptions,
	}

	service.writeJSON(w, http.StatusAccepted, responsePayload)
}

// 	// validate the user against the database
// 	user, err := app.Models.User.GetByEmail(requestPayload.Email)
// 	if err != nil {
// 		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
// 		return
// 	}

// 	valid, err := user.PasswordMatches(requestPayload.Password)
// 	if err != nil || !valid {
// 		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
// 		return
// 	}

// 	// log authentication
// 	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}

// 	payload := jsonResponse{
// 		Error:   false,
// 		Message: fmt.Sprintf("Logged in user %s", user.Email),
// 		Data:    user,
// 	}

// 	app.writeJSON(w, http.StatusAccepted, payload)
// }

// func (app *Config) logRequest(name, data string) error {
// 	var entry struct {
// 		Name string `json:"name"`
// 		Data string `json:"data"`
// 	}

// 	entry.Name = name
// 	entry.Data = data

// 	jsonData, _ := json.MarshalIndent(entry, "", "\t")
// 	logServiceURL := "http://logger-service/log"

// 	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err
// 	}

// 	client := &http.Client{}
// 	_, err = client.Do(request)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
