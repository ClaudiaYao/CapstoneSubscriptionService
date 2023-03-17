package data

import (
	"context"
	"database/sql"
)

// const dbTimeout = time.Second * 3

// C: Although DataService struct only contains one *sql.DB, using this struct
// C: Could allow to create own service
type DataQuery struct {
	DBConn *sql.DB
}

const getDisDeliveryCondition = `
select id, subscription_dish_id, status, expected_time, delivery_time, note FROM dish_delivery where subscription_dish_id = $1
`

func (dq *DataQuery) GetDishDeliveryCondition(ctx context.Context, subscriptionDishID string) ([]DishDelivery, error) {
	rows, err := dq.DBConn.QueryContext(ctx, getDisDeliveryCondition, subscriptionDishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DishDelivery
	for rows.Next() {
		var i DishDelivery
		if err := rows.Scan(
			&i.ID,
			&i.SubscriptionDishID,
			&i.Status,
			&i.ExpectedTime,
			&i.DeliveryTime,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDishBySubscriptionID = `
select id, dish_id, subscription_id, schedule_time, frequency, note FROM subscription_dish where subscription_id = $1
`

func (dq *DataQuery) GetDishBySubscriptionID(ctx context.Context, subscriptionID string) ([]SubscriptionDish, error) {
	rows, err := dq.DBConn.QueryContext(ctx, getDishBySubscriptionID, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SubscriptionDish
	for rows.Next() {
		var i SubscriptionDish
		if err := rows.Scan(
			&i.ID,
			&i.DishID,
			&i.SubscriptionID,
			&i.ScheduleTime,
			&i.Frequency,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubscriptionByID = `
select id, user_id, playlist_id, customized, status, frequency, start_date, end_date FROM subscription where id = $1
`

func (dq *DataQuery) GetSubscriptionByID(ctx context.Context, id string) (Subscription, error) {
	row := dq.DBConn.QueryRowContext(ctx, getSubscriptionByID, id)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PlaylistID,
		&i.Customized,
		&i.Status,
		&i.Frequency,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}

const getSubscriptionByUserID = `
select id, user_id, playlist_id, customized, status, frequency, start_date, end_date FROM subscription where user_id = $1
`

func (dq *DataQuery) GetSubscriptionByUserID(ctx context.Context, userID string) ([]Subscription, error) {
	rows, err := dq.DBConn.QueryContext(ctx, getSubscriptionByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.PlaylistID,
			&i.Customized,
			&i.Status,
			&i.Frequency,
			&i.StartDate,
			&i.EndDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertDishDelivery = `
insert into dish_delivery ("id", "subscription_dish_id", "status",
  "expected_time", "delivery_time", "note")
  values ($1, $2, $3, $4, $5, $6)
  returning id, subscription_dish_id, status, expected_time, delivery_time, note
`

func (dq *DataQuery) InsertDishDelivery(ctx context.Context, arg DishDelivery) (string, error) {
	row := dq.DBConn.QueryRowContext(ctx, insertDishDelivery,
		arg.ID,
		arg.SubscriptionDishID,
		arg.Status,
		arg.ExpectedTime,
		arg.DeliveryTime,
		arg.Note,
	)
	var i DishDelivery
	err := row.Scan(
		&i.ID,
		&i.SubscriptionDishID,
		&i.Status,
		&i.ExpectedTime,
		&i.DeliveryTime,
		&i.Note,
	)
	return i.ID, err
}

const insertDishes = `
insert into subscription_dish ("id", "dish_id", "subscription_id",
  "schedule_time", "frequency", "note")
  values ($1, $2, $3, $4, $5, $6)
  returning id, dish_id, subscription_id, schedule_time, frequency, note
`

func (dq *DataQuery) InsertDishes(ctx context.Context, arg SubscriptionDish) (string, error) {
	row := dq.DBConn.QueryRowContext(ctx, insertDishes,
		arg.ID,
		arg.DishID,
		arg.SubscriptionID,
		arg.ScheduleTime,
		arg.Frequency,
		arg.Note,
	)
	var i SubscriptionDish
	err := row.Scan(
		&i.ID,
		&i.DishID,
		&i.SubscriptionID,
		&i.ScheduleTime,
		&i.Frequency,
		&i.Note,
	)
	return i.ID, err
}

const insertSubscription = `
insert into subscription ("id", "user_id", "playlist_id",
  "customized", "status", "frequency", "start_date",
  "end_date" ) values ($1, $2, $3, $4, $5, $6, $7, $8)
  returning id, user_id, playlist_id, customized, status, frequency, start_date, end_date
`

func (dq *DataQuery) InsertSubscription(ctx context.Context, arg Subscription) (string, error) {
	row := dq.DBConn.QueryRowContext(ctx, insertSubscription,
		arg.ID,
		arg.UserID,
		arg.PlaylistID,
		arg.Customized,
		arg.Status,
		arg.Frequency,
		arg.StartDate,
		arg.EndDate,
	)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PlaylistID,
		&i.Customized,
		&i.Status,
		&i.Frequency,
		&i.StartDate,
		&i.EndDate,
	)
	return i.ID, err
}

const changeDishDeliveryStatus = `
update dish_delivery set status = $1 where subscription_dish_id = $2 and delivery_time is NULL
returning id, subscription_dish_id, status, expected_time, delivery_time, note`

func (dq *DataQuery) ChangeDishDeliveryStatus(ctx context.Context, toStatus string, subscriptionDishID string) ([]DishDelivery, error) {
	rows, err := dq.DBConn.QueryContext(ctx, changeDishDeliveryStatus, toStatus, subscriptionDishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishesDelivery []DishDelivery
	for rows.Next() {
		var dishDeliver DishDelivery
		if err := rows.Scan(
			&dishDeliver.ID,
			&dishDeliver.SubscriptionDishID,
			&dishDeliver.Status,
			&dishDeliver.ExpectedTime,
			&dishDeliver.DeliveryTime,
			&dishDeliver.Note,
		); err != nil {
			return nil, err
		}
		dishesDelivery = append(dishesDelivery, dishDeliver)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dishesDelivery, nil
}

const changeSubscriptionStatus = `
update subscription set status = $1 where id = $2
returning id, user_id, playlist_id, customized, status, frequency, start_date, end_date
`

func (dq *DataQuery) ChangeSubscriptionStatus(ctx context.Context, toStatus, subscriptionID string) (Subscription, error) {
	row := dq.DBConn.QueryRowContext(ctx, changeSubscriptionStatus, toStatus, subscriptionID)
	var sub Subscription
	err := row.Scan(
		&sub.ID,
		&sub.UserID,
		&sub.PlaylistID,
		&sub.Customized,
		&sub.Status,
		&sub.Frequency,
		&sub.StartDate,
		&sub.EndDate,
	)
	return sub, err
}
