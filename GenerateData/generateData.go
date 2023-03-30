package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid"
)

type allIDInfo struct {
	categories           []string
	restaurants          []string
	dishes               []string
	playlists            []string
	playlistDishRelation map[string][]string
	subscriptions        []string
}

// open address file and store them in a slice of structs
type address struct {
	unit_number   string
	address_line1 string
	address_line2 string
	postal_code   int
}

func main() {
	categoryCodes := generateCategory()
	restaurantIDs := generateRestaurant(categoryCodes)
	// // since Dish table has a foreign reference key to RestaurantID, so we
	// // pass restaurantIDs to the dish generation function
	dishIDs := generateDish(restaurantIDs)
	playlistIDs := generatePlaylist(categoryCodes)

	playlistDishRelation := generatePlaylistDishRelation(playlistIDs, dishIDs)
	// // when generating subscription data, subscriptionDish file and DishDelivery file is gnerated
	// // together because they are closely related and need to refer to the data from playlist
	// // service
	subscriptionIDs := generateSubscription(playlistIDs, playlistDishRelation, dishIDs)

	keepData(categoryCodes, restaurantIDs, dishIDs, playlistIDs, playlistDishRelation, subscriptionIDs)

}

func keepData(categories []string, restaurantsIDs []string, dishIDs []string, playlistIDs []string, playlistDishRelations map[string][]string, subscriptionIDs []string) {
	write_f, err := os.Create("Generated/all_IDs.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the method
	defer write_f.Close()

	data := strings.Join(categories, "|") + "\n"
	data += strings.Join(restaurantsIDs, "|") + "\n"
	data += strings.Join(dishIDs, "|") + "\n"
	data += strings.Join(playlistIDs, "|") + "\n"
	data += strings.Join(subscriptionIDs, "|") + "\n"

	_, err = write_f.WriteString(data)
	if err != nil {
		log.Fatal("error occurs when writing to file:", err)
	}

	write_rel, err := os.Create("Generated/playlistDishRel.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the method
	defer write_rel.Close()
	data = ""
	for key, value := range playlistDishRelations {
		data += key + "|" + strings.Join(value, "|") + "\n"
	}
	_, err = write_rel.WriteString(data)
	if err != nil {
		log.Fatal("error occurs when writing to file:", err)
	}

}

func extractIDs() (allIDInfo, error) {
	readFile, err := os.Open("Generated/all_IDs.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	i := 0
	ids := allIDInfo{}

	for fileScanner.Scan() {
		items := strings.Split(fileScanner.Text(), "|")

		switch i {
		case 0:
			ids.categories = items
		case 1:
			ids.restaurants = items
		case 2:
			ids.dishes = items
		case 3:
			ids.playlists = items
		case 4:
			ids.subscriptions = items
		}
		i += 1

	}

	readFileRel, err := os.Open("Generated/playlistDishRel.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer readFileRel.Close()

	fileScanner = bufio.NewScanner(readFileRel)

	fileScanner.Split(bufio.ScanLines)

	playlistDishRel := map[string][]string{}

	for fileScanner.Scan() {
		items := strings.Split(fileScanner.Text(), "|")
		playlistDishRel[items[0]] = items[1:]
	}

	ids.playlistDishRelation = playlistDishRel
	return ids, err

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

	new_text := ""

	// Parse the data record
	if status == "Completed" {
		new_text = id + "|" + subscriptionDishID + "|" + status + "|" +
			ExpectedTime.Format(time.RFC3339) + "|" +
			deliveryTime.Format(time.RFC3339) + "|" + note
	} else {
		new_text = id + "|" + subscriptionDishID + "|" + status + "|" +
			ExpectedTime.Format(time.RFC3339) + "|" +
			"\\N" + "|" + "not deliver yet"
	}
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

		var dishOptions = make([][]string, 2)
		if rand.Intn(2) < 1 {
			dishOptions = [][]string{
				{"Mentaico Source", "Yes"},
				{"Wasabi", "No"},
			}
		} else {
			dishOptions = [][]string{
				{"More source", "No"},
				{"Pepper and Chili", "No"},
			}
		}

		optionsB, err := json.Marshal(dishOptions)
		if err != nil {
			log.Fatal(err)
		}

		new_text := id + "|" + dishID + "|" + subscriptionID + "|" +
			scheduleTime.Format(time.RFC3339) +
			"|" + frequency + "|" + string(optionsB) + "|" + note

		nextTime := scheduleTime
		for nextTime.Before(endDate) {
			generateDishDeliveryRecords(id, nextTime)
			nextTime = nextDelivery(frequency, nextTime)
		}

		_, err = write_f.WriteString(new_text + "\n")
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
		"Pending",
		"Expired",
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

	names := []string{"Tody Liang", "Jone Tew", "James"}

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
			receiverName := names[rand.Intn(len(names))]
			receiverContact := strconv.Itoa(80000000 + rand.Intn(200000))

			new_text := ""
			if customized == "false" {
				new_text = id + "|" + userID + "|" + playlistID + "|" +
					customized + "|" + status + "|" + frequency + "|" + startDate.Format(time.RFC3339) +
					"|" + endDate.Format(time.RFC3339) + "|" + receiverName +
					"|" + receiverContact

				generateSubscriptionDishes(id, startDate, endDate, frequency, playlistDishRelations[playlistID])

			} else {
				new_text = id + "|" + userID + "|" + "(empty)" + "|" +
					customized + "|" + status + "|" + frequency + "|" + startDate.Format(time.RFC3339) +
					"|" + endDate.Format(time.RFC3339) + "|" + receiverName +
					"|" + receiverContact

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

	playlistIDs := []string{}

	// Get each playlist name and then form the table content
	for scanner.Scan() {
		id := "Play" + shortuuid.New()
		// do something with a line
		name := scanner.Text()
		categoryCode := categoryCodes[rand.Intn(categoryNumber)]
		dietary := dietaryInfo[rand.Intn(dietaryNumber)]

		startDateRandom := time.Now().AddDate(0, 0, rand.Intn(10)-10)
		startDate := startDateRandom.Format("2006-01-02")
		endDateRandom := startDateRandom.AddDate(0, 2, 15)
		end_date := endDateRandom.Format("2006-01-02")
		popularity := 1 + rand.Intn(5)

		statusInfo := ""
		if startDateRandom.After(time.Now()) {
			statusInfo = "Pending"
		} else if endDateRandom.Before(time.Now()) {
			statusInfo = "Expired"
		} else if endDateRandom.After(time.Now()) {
			statusInfo = "Active"
		}

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
		str := strings.Split(scanner.Text(), ",")
		name := strings.TrimSpace(str[0])
		imageUrl := strings.TrimSpace(str[1])

		restaurantID := restaurantIDs[rand.Intn(restaurantNumber)]
		comment := "extra information"
		price := fmt.Sprintf("%.2f", 4.0+rand.Float32()*20)
		if err != nil {
			fmt.Println(err)
		}
		cuisineStyle := cuisineStyles[rand.Intn(cuisineStyleNumber)]
		ingredient := ingredients[rand.Intn(ingredientNumber)]

		var dishOptions = make([][]string, 2)
		if rand.Intn(2) < 1 {
			dishOptions = [][]string{
				{"Mentaico Source", "Yes", "No"},
				{"Wasabi", "Yes", "No"},
			}
		} else {
			dishOptions = [][]string{
				{"More source", "Yes", "No"},
				{"Pepper and Chili", "Yes", "No"},
			}
		}

		optionsB, err := json.Marshal(dishOptions)
		if err != nil {
			log.Fatal(err)
		}

		new_text := id + "|" + name + "|" + restaurantID + "|" +
			price + "|" + cuisineStyle + "|" + ingredient + "|" +
			string(optionsB) + "|" + comment + "|" + imageUrl

		_, err = write_f.WriteString(new_text + "\n")
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

func generateRestaurant(categories []string) []string {
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

	urls := []string{
		"https://example1.com",
		"https://example2.com",
		"https://example3.com",
	}
	operationStartOptions := []string{"0600", "0700", "0800", "0900", "1100", "1200"}
	operationEndOptions := []string{"1900", "2100", "2300", "2400", "0200"}

	// read the file line by line using scanner
	scanner := bufio.NewScanner(read_f)
	address_number := len(addresses)

	restaurantIDs := []string{}

	for scanner.Scan() {
		id := "RES" + shortuuid.New()
		// do something with a line
		name := scanner.Text()
		restaurant := addresses[rand.Intn(address_number)]
		openStart := operationStartOptions[rand.Intn(len(operationStartOptions))]
		openEnd := operationEndOptions[rand.Intn(len(operationEndOptions))]

		operationHours := make([][]string, 7)
		for i := range operationHours {
			if i < 5 {
				operationHours[i] = []string{openStart, openEnd}
			} else {
				openEnd = operationEndOptions[rand.Intn(len(operationEndOptions))]
				operationHours[i] = []string{openStart, openEnd}
			}
		}

		operationB, err := json.Marshal(operationHours)
		if err != nil {
			log.Fatal(err)
		}

		logoUrl := urls[rand.Intn(len(urls))]
		headerUrl := urls[rand.Intn(len(urls))]
		tag := categories[rand.Intn(len(categories))]

		new_text := id + "|" + name + "|" + restaurant.unit_number + "|" +
			restaurant.address_line1 + "|" + restaurant.address_line2 + "|" +
			strconv.Itoa(restaurant.postal_code) + "|" + tag + "|" + string(operationB) +
			"|" + logoUrl + "|" + headerUrl

		_, err = write_f.WriteString(new_text + "\n")
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
