package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ClaudiaYao/CapstoneSubscriptionService/app/data"
	"github.com/lithammer/shortuuid"
)

func (service *SubscriptionService) GenerateNewSubscription(ctx context.Context, subscription data.Subscription, dishes []data.SubscriptionDish) ([]data.DishDelivery, error) {
	_, err := service.DBConnection.InsertSubscription(ctx, subscription)
	if err != nil {

		return nil, errors.New(fmt.Sprint("error when inserting the subscription: ", err))
	}

	dishesDelivery := []data.DishDelivery{}

	for _, dish := range dishes {
		dishID, err := service.DBConnection.InsertDishes(ctx, dish)

		if err != nil {
			return nil, errors.New(fmt.Sprint("error when inserting the subscription dishes:", err))
		}

		nextTime := dish.ScheduleTime
		for !nextTime.After(subscription.EndDate) {
			dishDelivery := data.DishDelivery{
				ID:                 "DD" + shortuuid.New(),
				SubscriptionDishID: dishID,
				Status:             "Active",
				ExpectedTime:       nextTime,
				Note:               dish.Note,
			}
			service.DBConnection.InsertDishDelivery(ctx, dishDelivery)
			nextTime = nextDelivery(dish.Frequency, nextTime)
			dishesDelivery = append(dishesDelivery, dishDelivery)
		}

	}
	return dishesDelivery, err

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
