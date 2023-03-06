package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid"
)

// open address file and store them in a slice of structs
type address struct {
	unit_number   string
	address_line1 string
	address_line2 string
	postal_code   int
}

func main() {
	categoryCodes := generateCategory()
	restaurantIDs := generateRestaurant()
	// since Dish table has a foreign reference key to RestaurantID, so we
	// pass restaurantIDs to the dish generation function
	dishIDs := generateDish(restaurantIDs)
	playlistIDs := generatePlaylist(categoryCodes)

	playlistDishRelation := generatePlaylistDishRelation(playlistIDs, dishIDs)

	// when generating subscription data, subscriptionDish file and DishDelivery file is gnerated
	// together because they are closely related and need to refer to the data from playlist
	// service
	generateSubscription(playlistIDs, playlistDishRelation, dishIDs)

}

func generateDishDeliveryRecords(subscriptionDishID string, ExpectedTime time.Time) {
	write_f, err := os.OpenFile("Generated/DishesDelivery.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer write_f.Close()

	statusChoice := []string{"Completed", "Cancelled", "Pending"}
	id := "DD" + shortuuid.New()
	status := statusChoice[rand.Intn(len(statusChoice))]
	deliveryTime := ExpectedTime.Add(time.Minute * time.Duration(rand.Intn(500)))
	note := "on time"

	// Parse the data record
	new_text := id + "|" + subscriptionDishID + "|" + status + "|" +
		ExpectedTime.Format(time.RFC3339) + "|" +
		deliveryTime.Format(time.RFC3339) + "|" + note

	_, err = write_f.WriteString(new_text + "\n")
	if err != nil {
		log.Fatal("error occurs when writing to file:", err)
	}

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

func generateSubscriptionDishes(subscriptionID string, startDate time.Time, endDate time.Time, frequency string, dishIDs []string) []string {
	write_f, err := os.OpenFile("Generated/subscriptionDishes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer write_f.Close()

	// customized means the user chooses own playlist

	subscriptionDishIDs := []string{}

	noteChoice := []string{
		"perferendis volupt",
		"veniam ut exc",
		"fugit, amet",
		"nulla consequatur",
	}
	for _, dishID := range dishIDs {
		id := "SDish" + shortuuid.New()
		scheduleTime := startDate.AddDate(0, 0, rand.Intn(3))
		note := noteChoice[rand.Intn(len(noteChoice))]

		new_text := id + "|" + dishID + "|" + subscriptionID + "|" +
			scheduleTime.Format(time.RFC3339) +
			"|" + frequency + "|" + note

		nextTime := scheduleTime
		for nextTime.Before(endDate) {
			generateDishDeliveryRecords(id, nextTime)
			nextTime = nextDelivery(frequency, nextTime)
		}

		_, err := write_f.WriteString(new_text + "\n")
		if err != nil {
			log.Fatal("error occurs when writing to file:", err)
		}
		subscriptionDishIDs = append(subscriptionDishIDs, id)

	}

	return subscriptionDishIDs

}

func generateSubscription(playlistIDs []string, playlistDishRelations map[string][]string, dishIDs []string) []string {

	write_f, err := os.Create("Generated/subscription.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the method
	defer write_f.Close()

	write_dishes_f, err := os.Create("Generated/subscriptionDishes.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the method
	write_dishes_f.Close()

	write_dishes_delivery, err := os.Create("Generated/dishesDelivery.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the method
	write_dishes_delivery.Close()

	statusInfo := []string{
		"Active",
		"Cancelled",
		"Pause",
		"Finish",
	}

	frequencyChoices := []string{
		"Daily",
		"Weekly",
		"Monthly",
	}

	customizedChoice := []string{
		"true",
		"false",
	}

	subscriptionIDs := []string{}

	// Get each playlist id and then form the table content
	for _, playlistID := range playlistIDs {
		chosen := rand.Intn(10)
		// for each playlistID, generate some subscriptions
		for i := 0; i < chosen; i++ {
			id := "Sub" + shortuuid.New()
			userID := "user" + strconv.Itoa(rand.Intn(20))
			customized := customizedChoice[rand.Intn(len(customizedChoice))]
			frequency := frequencyChoices[rand.Intn(len(frequencyChoices))]
			status := statusInfo[rand.Intn(len(statusInfo))]
			startDate := time.Now().AddDate(0, 0, -rand.Intn(15))
			endDate := startDate.AddDate(0, 0, 7*rand.Intn(3))

			new_text := ""
			if customized == "false" {
				new_text = id + "|" + userID + "|" + playlistID + "|" +
					customized + "|" + status + "|" + frequency + "|" + startDate.Format(time.RFC3339) +
					"|" + endDate.Format(time.RFC3339)

				generateSubscriptionDishes(id, startDate, endDate, frequency, playlistDishRelations[playlistID])

			} else {
				new_text = id + "|" + userID + "|" + "(empty)" + "|" +
					customized + "|" + status + "|" + frequency + "|" + startDate.Format(time.RFC3339) +
					"|" + endDate.Format(time.RFC3339)

				randomDishes := []string{}
				for i := 0; i < rand.Intn(5); i++ {
					randomDishes = append(randomDishes, dishIDs[i])
				}

				generateSubscriptionDishes(id, startDate, endDate, frequency, randomDishes)
			}

			_, err := write_f.WriteString(new_text + "\n")
			if err != nil {
				log.Fatal("error occurs when writing to file:", err)
			}
			subscriptionIDs = append(subscriptionIDs, id)
		}

	}
	return subscriptionIDs
}

func generatePlaylist(categoryCodes []string) []string {

	read_f, err := os.Open("Initial/playlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer read_f.Close()

	write_f, err := os.Create("Generated/playlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer write_f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(read_f)
	categoryNumber := len(categoryCodes)

	dietaryInfo := []string{
		"perferendis voluptatibus veniam",
		"veniam ut excepturi nulla conse",
		"fugit, amet quia nulla culpa",
		"nulla consequatur natus tempore officiis",
	}
	dietaryNumber := len(dietaryInfo)

	statusInfo := []string{
		"Active",
		"Expired",
		"Pending",
	}

	playlistIDs := []string{}

	// Get each playlist name and then form the table content
	for scanner.Scan() {
		id := "Play" + shortuuid.New()
		// do something with a line
		name := scanner.Text()
		categoryCode := categoryCodes[rand.Intn(categoryNumber)]
		dietary := dietaryInfo[rand.Intn(dietaryNumber)]
		statusInfo := statusInfo[rand.Intn(len(statusInfo))]
		startDateRandom := time.Now().AddDate(0, 0, rand.Intn(60)-30)
		startDate := startDateRandom.Format("2006-01-02")
		endDateRandom := startDateRandom.AddDate(0, 2, 15)
		end_date := endDateRandom.Format("2006-01-02")
		popularity := 1 + rand.Intn(5)

		new_text := id + "|" + name + "|" + categoryCode + "|" +
			dietary + "|" + statusInfo + "|" + startDate +
			"|" + end_date + "|" + strconv.Itoa(popularity)

		_, err := write_f.WriteString(new_text + "\n")
		if err != nil {
			log.Fatal("error occurs when writing to file:", err)
		}
		playlistIDs = append(playlistIDs, id)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return playlistIDs

}

func generatePlaylistDishRelation(playlistIDs []string, dishIDs []string) map[string][]string {

	write_f, err := os.Create("Generated/playlist_dish.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer write_f.Close()

	// read the file line by line using scanner

	totalDish := len(dishIDs)
	if totalDish < 8 {
		fmt.Println("too less dishes in the database.")
	}

	playlistDishRelation := map[string][]string{}

	// Get each playlist name and then form the table content
	for _, playlistID := range playlistIDs {

		// do something with a line

		dishNum := 4 + rand.Intn(6)
		chosen := 0

		dishes := []string{}

		for chosen < dishNum {
			id := "PD" + shortuuid.New()
			dishID := dishIDs[rand.Intn(totalDish)]
			new_text := id + "|" + dishID + "|" + playlistID

			_, err := write_f.WriteString(new_text + "\n")
			if err != nil {
				log.Fatal("error occurs when writing to file:", err)
			}
			dishes = append(dishes, dishID)
			chosen += 1

		}
		playlistDishRelation[playlistID] = dishes

	}
	return playlistDishRelation
}

func generateDish(restaurantIDs []string) []string {

	read_f, err := os.Open("Initial/dish.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer read_f.Close()

	write_f, err := os.Create("Generated/dish.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer write_f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(read_f)
	restaurantNumber := len(restaurantIDs)
	dishIDs := []string{}

	cuisineStyles := []string{
		"perferendis voluptatibus veniam",
		"veniam ut excepturi nulla conse",
		"fugit, amet quia nulla culpa",
		"nulla consequatur natus tempore officiis",
	}
	cuisineStyleNumber := len(cuisineStyles)

	ingredients := []string{
		"voluptatibus, veniam",
		"veniam, ut, excepturi",
		"fugit, amet quia, nulla",
		"consequatur, natus, tempore, officiis",
	}
	ingredientNumber := len(ingredients)

	// Get each dish name and then form the table content
	for scanner.Scan() {
		id := "Dish" + shortuuid.New()
		// do something with a line
		name := scanner.Text()
		restaurantID := restaurantIDs[rand.Intn(restaurantNumber)]
		comment := ""
		price := fmt.Sprintf("%.2f", 4.0+rand.Float32()*20)
		if err != nil {
			fmt.Println(err)
		}
		cuisineStyle := cuisineStyles[rand.Intn(cuisineStyleNumber)]
		ingredient := ingredients[rand.Intn(ingredientNumber)]

		new_text := id + "|" + name + "|" + restaurantID + "|" +
			price + "|" + cuisineStyle + "|" + ingredient +
			"|" + comment

		_, err := write_f.WriteString(new_text + "\n")
		if err != nil {
			log.Fatal("error occurs when writing to file:", err)
		}
		dishIDs = append(dishIDs, id)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return dishIDs

}
func generateCategory() []string {
	// open file
	read_f, err := os.Open("Initial/category.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer read_f.Close()

	// create file
	features := []string{
		"Temporibus aliquid, obcaecati soluta consequatur a veritatis ad omnis",
		"Rerum sequi, earum delectus quidem tenetur est dicta exercitationem eius labore ipsa",
		"Quibusdam perferendis voluptatibus veniam ut excepturi nulla",
	}

	write_f, err := os.Create("Generated/category.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer write_f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(read_f)
	feature_number := len(features)

	categoryCodes := []string{}
	for scanner.Scan() {
		// do something with a line
		line := scanner.Text()
		code := line[:3]
		new_text := code + "|" + line + "|" + features[rand.Intn(feature_number)]
		_, err := write_f.WriteString(new_text + "\n")
		if err != nil {
			log.Fatal("error occurs when writing to file:", err)
		}

		categoryCodes = append(categoryCodes, code)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return categoryCodes
}

func GenerateAddress() []address {

	var addresses []address
	// read the file line by line using scanner
	read_address, err := os.Open("Initial/address.txt")

	if err != nil {
		log.Fatal("read address file:", err)
	}
	defer read_address.Close()

	scanner := bufio.NewScanner(read_address)

	for scanner.Scan() {
		line := scanner.Text()
		adds := strings.Split(line, ",")
		if len(adds) == 3 {
			postal, err := strconv.Atoi(strings.TrimSpace(adds[2]))
			if err != nil {
				log.Fatal("postal code format is incorrect: ", err)
			}
			new_address := address{unit_number: strings.TrimSpace(adds[0]),
				address_line1: strings.TrimSpace(adds[1]),
				address_line2: "(empty)", postal_code: postal}

			addresses = append(addresses, new_address)
		} else if len(adds) == 4 {
			postal, err := strconv.Atoi(strings.TrimSpace(adds[3]))
			if err != nil {
				log.Fatal("postal code format is incorrect.", err)
			}
			new_address := address{unit_number: strings.TrimSpace(adds[0]),
				address_line1: strings.TrimSpace(adds[1]),
				address_line2: strings.TrimSpace(adds[2]),
				postal_code:   postal}

			addresses = append(addresses, new_address)
		} else {
			log.Fatal("format of the address.txt is incorrect.")
		}

	}
	return addresses
}

func generateRestaurant() []string {
	addresses := GenerateAddress()

	read_f, err := os.Open("Initial/restaurant.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer read_f.Close()

	write_f, err := os.Create("Generated/restaurant.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer write_f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(read_f)
	address_number := len(addresses)

	restaurantIDs := []string{}

	for scanner.Scan() {
		id := "RES" + shortuuid.New()
		// do something with a line
		name := scanner.Text()
		restaurant := addresses[rand.Intn(address_number)]

		new_text := id + "|" + name + "|" + restaurant.unit_number + "|" +
			restaurant.address_line1 + "|" + restaurant.address_line2 + "|" + strconv.Itoa(restaurant.postal_code)

		_, err := write_f.WriteString(new_text + "\n")
		if err != nil {
			log.Fatal("error occurs when writing to file:", err)
		}
		restaurantIDs = append(restaurantIDs, id)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return restaurantIDs
}

// first_names := []string{"James", "Robert", "John", "Michael", "David", "William", "Richard",
// 	"Joseph", "Thomas", "Charles", "Christopher", "Daniel", "Matthew", "Anthony", "Mark",
// 	"Donald", "Steven", "Paul", "Andrew", "Joshua", "Kenneth", "Kevin", "Brian", "George",
// 	"Timothy", "Patricia", "Jennifer", "Linda", "Elizabeth", "Barbara", "Susan",
// 	"Jessica", "Sarah", "Karen", "Lisa", "Nancy", "Betty", "Margaret", "Sandra",
// 	"Ashley", "Kimberly", "Emily", "Donna", "Michelle", "Carol", "Amanda", "Dorothy",
// 	"Melissa", "Deborah", "Stephanie", "Rebecca"}
// len_first := len(first_names)

// last_names := []string{"Tan", "Lim", "Lee", "Ng", "Ong", "Wong", "Goh",
// 	"Chua", "Chan", "Koh", "Teo", "Ang", "Yeo", "Tay", "Ho", "Low", "Toh", "Sim",
// 	"Chong", "Chia", "Seah"}
// len_last := len(last_names)

// email_dn := []string{"@gmail.com", "@hotmail.com", "@yahoo.com", "@dental.com"}
// len_dn := len(email_dn)

// // all the methods to generate unique user id: https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html
// // use a shorter version

// for count := 0; count < 200; count++ {
// 	id := shortuuid.New()

// 	first_name := first_names[rand.Intn(len_first)]
// 	last_name := last_names[rand.Intn(len_last)]
// 	user_name := first_name + last_name + id[:4]
// 	email_add := first_name + "." + last_name + id[:4] + email_dn[rand.Intn(len_dn)]
// 	phone_num := strconv.Itoa(70000000 + rand.Intn(10000000))
// 	bPassword, err := bcrypt.GenerateFromPassword([]byte(user_name), bcrypt.MinCost)
// 	if err != nil {
// 		log.Panic("could not generate test data password successfully.")
// 		return
// 	}
// 	user := User{UserID: id, UserName: user_name, FirstName: first_name,
// 		LastName: last_name, Password: bPassword, EmailAddress: email_add,
// 		ContactPhone: phone_num}
// 	UserTestMap[id] = user

// }
// SaveToTestUserJSON()
