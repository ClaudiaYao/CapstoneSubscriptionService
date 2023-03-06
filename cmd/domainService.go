package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lithammer/shortuuid"
)

func (app *SubscriptionService) GenerateNewSubscription(ctx context.Context, subscription Subscription, dishes []SubscriptionDish) ([]DishDelivery, error) {
	_, err := app.DBConnection.InsertSubscription(ctx, subscription)
	if err != nil {

		return nil, errors.New(fmt.Sprint("error when inserting the subscription: ", err))
	}

	dishesDelivery := []DishDelivery{}

	for _, dish := range dishes {
		dishID, err := app.DBConnection.InsertDishes(ctx, dish)

		if err != nil {
			return nil, errors.New(fmt.Sprint("error when inserting the subscription dishes:", err))
		}

		nextTime := dish.ScheduleTime
		for !nextTime.After(subscription.EndDate) {
			dishDelivery := DishDelivery{
				ID:                 "DD" + shortuuid.New(),
				SubscriptionDishID: dishID,
				Status:             "Active",
				ExpectedTime:       nextTime,
				Note:               dish.Note,
			}
			app.DBConnection.InsertDishDelivery(ctx, dishDelivery)
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
