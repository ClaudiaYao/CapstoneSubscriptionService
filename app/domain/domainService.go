package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ClaudiaYao/CapstoneSubscriptionService/app/data"
	"github.com/lithammer/shortuuid"
)

func (service *SubscriptionService) InsertNewSubscriptionRecord(ctx context.Context, payload SubscriptionServiceRequestDataDTO) (*SubscriptionServiceResponseDataDTO, error) {
	subReq := payload.SubscriptionRequest

	// this ensures that every time posting the request to create subscription, the id will be different.
	subInfo := data.Subscription{
		ID:         "Sub" + shortuuid.New(),
		UserID:     subReq.UserID,
		PlaylistID: subReq.PlaylistID,
		Customized: subReq.Customized,
		Status:     "Active",
		Frequency:  subReq.Frequency,
		StartDate:  subReq.StartDate,
		EndDate:    subReq.EndDate,
	}

	dishIncluded := payload.DishIncluded
	dishes := []data.SubscriptionDish{}

	for _, dishInfo := range dishIncluded {
		dish := data.SubscriptionDish{
			ID:             "SDish" + shortuuid.New(),
			DishID:         dishInfo.DishID,
			SubscriptionID: subInfo.ID,
			ScheduleTime:   dishInfo.ScheduleTime,
			Frequency:      dishInfo.Frequency,
			Note:           dishInfo.Note,
		}
		dishes = append(dishes, dish)
	}

	_, err := service.DBConnection.InsertSubscription(ctx, subInfo)
	if err != nil {
		return nil, errors.New(fmt.Sprint("error when inserting the subscription: ", err))
	}

	// dishesDelivery := []data.DishDelivery{}

	for _, dish := range dishes {
		dishID, err := service.DBConnection.InsertDishes(ctx, dish)

		if err != nil {
			return nil, errors.New(fmt.Sprint("error when inserting the subscription dishes:", err))
		}

		nextTime := dish.ScheduleTime
		for !nextTime.After(subInfo.EndDate) {
			dishDelivery := data.DishDelivery{
				ID:                 "DD" + shortuuid.New(),
				SubscriptionDishID: dishID,
				Status:             "Active",
				ExpectedTime:       nextTime,
				Note:               dish.Note,
			}
			service.DBConnection.InsertDishDelivery(ctx, dishDelivery)
			nextTime = nextDelivery(dish.Frequency, nextTime)
			// dishesDelivery = append(dishesDelivery, dishDelivery)
		}

	}

	response := &SubscriptionServiceResponseDataDTO{
		Subscription: subInfo,
		DishIncluded: dishes,
	}
	return response, err

}

func (service *SubscriptionService) CancelSubscriptionRelatedRecords(ctx context.Context, subscriptionID string) ([]data.DishDelivery, error) {

	sub, err := service.DBConnection.ChangeSubscriptionStatus(ctx, "Cancelled", subscriptionID)

	if err != nil {
		return nil, errors.New(fmt.Sprint("error when changing the subscription: ", err))
	}

	dishes, err := service.DBConnection.GetDishBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, errors.New(fmt.Sprint("error when querying the dishes: ", err))
	}

	dishDeliveryInfo := []data.DishDelivery{}

	for _, dish := range dishes {
		DeliveryInfo, err := service.DBConnection.ChangeDishDeliveryStatus(ctx, sub.Status, dish.DishID)

		if err != nil {
			return nil, errors.New(fmt.Sprint("error when updating the dish delivery status: ", err))
		}
		dishDeliveryInfo = append(dishDeliveryInfo, DeliveryInfo...)
	}

	return dishDeliveryInfo, err

}

func nextDelivery(frequency string, thisDelivery time.Time) time.Time {
	if frequency == "daily" {
		return thisDelivery.AddDate(0, 0, 1)
	} else if frequency == "weekly" {
		return thisDelivery.AddDate(0, 0, 7)

	} else if frequency == "monthly" {
		return thisDelivery.AddDate(0, 1, 0)
	}
	return thisDelivery.AddDate(0, 0, 1)
}
