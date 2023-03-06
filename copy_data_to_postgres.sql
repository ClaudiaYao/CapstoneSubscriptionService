
\c subscription;
\COPY subscription FROM myData/subscription.txt WITH (FORMAT text, DELIMITER '|');
\COPY subscription_dish FROM myData/subscriptionDishes.txt WITH (FORMAT text, DELIMITER '|');
\COPY dish_delivery FROM myData/dishesDelivery.txt WITH (FORMAT text, DELIMITER '|');
