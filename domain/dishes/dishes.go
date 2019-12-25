package dishes

type food struct {
	FoodID         int    `json:"FoodID"`
	LocationID     int    `json:"LocationID"`
	LocationString string `json:"LocationString"`
	Description    string `json:"FoodDescription"`
	TimeCreated    string `json:"TimeCreated"`
	TimeExpires    string `json:"TimeExpires"`
}

type foods []food
