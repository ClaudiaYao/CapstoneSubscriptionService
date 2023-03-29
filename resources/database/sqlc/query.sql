-- name: GetSubscriptionByID :one
select * FROM subscription where id = $1;

-- name: GetSubscriptionByUserID :many
select * FROM subscription where user_id = $1;

-- name: GetDishBySubscriptionID :many
select * FROM subscription_dish where subscription_id = $1;

-- name: GetDisDeliveryCondition :many
select * FROM dish_delivery where subscription_dish_id = $1;

-- name: InsertSubscription :one
insert into subscription ("id", "user_id", "playlist_id",
  "customized", "status", "frequency", "start_date",
  "end_date", "receiver_name", "receiver_contact") values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
  returning *;

-- name: InsertDishes :one
insert into subscription_dish ("id", "dish_id", "subscription_id",
  "schedule_time", "frequency", "dish_options", "note")
  values ($1, $2, $3, $4, $5, $6, $7)
  returning *;

-- name: InsertDishDelivery :one
insert into dish_delivery ("id", "subscription_dish_id", "status",
  "expected_time", "delivery_time", "note")
  values ($1, $2, $3, $4, $5, $6)
  returning *;

-- name: ChangeSubscriptionStatus :one
update subscription set status = $1 where id = $2
returning *;

-- name: ChangeDishDeliveryStatus :many
update dish_delivery set status = $1 where subscription_dish_id = $2 and status = "(empty)"
returning *;

