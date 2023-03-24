package domain

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ClaudiaYao/CapstoneSubscriptionService/app/data"
	auth "github.com/ClaudiaYao/CapstoneSubscriptionService/app/domain/auth"
	"github.com/go-chi/chi"
)

type SubscriptionService struct {
	DBConnection *data.DataQuery
	JwtMaker     *auth.JWTMaker
	JwtVerifier  *auth.JWTVerifier
	AppConfig    *AppConfiguration
}

type AppConfiguration struct {
	TokenExpireSecs                  int
	ServicePort                      string
	EmailServiceContainerName        string
	PlaylistServiceContainerName     string
	SubscriptionServiceContainerName string
	LoginServiceContainerName        string
}

// this SubscriptionServiceDataDTO represents the data returned to the client
type SubscriptionServiceResponseDataDTO struct {
	Subscription data.Subscription
	DishIncluded []data.SubscriptionDish
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
	EndDate    time.Time `json:"endDate,omitempty"`
}

type SubscriptionDishRequested struct {
	DishID       string    `json:"dishID"`
	ScheduleTime time.Time `json:"scheduleTime"`
	Frequency    string    `json:"frequency"`
	Note         string    `json:"Note,omitempty"`
}

// C: this PlaylistService is responsible for transfering information request/response
// C: the database operation is conducted by its member *sql.DB
// C: when designing API or micro-service, the service request passes data via JSON
func (service *SubscriptionService) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// how to define the structure of SubscriptionServiceData is depending on the front end

	// add validation for the request data
	var requestPayload SubscriptionServiceRequestDataDTO

	err := service.readJSON(w, r, &requestPayload)
	if err != nil {
		service.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	subReqServiceDTO, err := service.InsertNewSubscriptionRecord(r.Context(), requestPayload)

	if err != nil {
		service.errorJSON(w, errors.New("invalid query"), http.StatusBadRequest)
		return
	}

	mailMsg := data.MailPayload{
		From:    "test@example.com",
		To:      "user@example.com",
		Subject: "subscription is success.",
		Message: "detailed information. TODO",
	}

	toAddress, err := service.SendEmail(r.Context(), mailMsg)

	if err != nil {
		service.errorJSON(w, errors.New("email sending failure"), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "subscription is created and mail sent to " + toAddress,
		Data:    subReqServiceDTO,
	}

	service.writeJSON(w, http.StatusAccepted, responsePayload)
}

func (service *SubscriptionService) Welcome(w http.ResponseWriter, r *http.Request) {
	service.writeJSON(w, http.StatusAccepted, "Welcome to Subscription service!")
}

func (service *SubscriptionService) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionID := chi.URLParam(r, "subscription_id")

	CancelledDishes, err := service.CancelSubscriptionRelatedRecords(r.Context(), subscriptionID)

	if err != nil {
		service.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("subscription %s is cancelled", subscriptionID),
		Data:    CancelledDishes,
	}

	service.writeJSON(w, http.StatusAccepted, responsePayload)
}

func (service *SubscriptionService) GetDishBySubscriptionID(w http.ResponseWriter, r *http.Request) {

	subscriptionId := chi.URLParam(r, "subscription_id")
	subscriptionDishes, err := service.DBConnection.GetDishBySubscriptionID(r.Context(), subscriptionId)
	if err != nil {
		service.errorJSON(w, errors.New(fmt.Sprint("invalid query", err)), http.StatusBadRequest)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "dishes are retrieved",
		Data:    subscriptionDishes,
	}

	service.writeJSON(w, http.StatusAccepted, responsePayload)

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

func (service *SubscriptionService) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("authorization")

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			service.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			service.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			service.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		payload, err := service.JwtVerifier.GetMetaData(accessToken)
		if err != nil {
			service.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		// add userID to the header of the request
		ctx := context.WithValue(r.Context(), "userID", payload.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
