package dish

//Dish type is the struct in the Domain that contains all the fields for what a Dish is.
type Dish struct {
	DishID      int    `json:"DishID"`
	UserID      int    `json:"UserID"`
	StorageID   int    `json:"StorageID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	CreatedDate string `json:"TimeCreated"`
	ExpireDate  string `json:"TimeExpires"`
	Priority    string `json:"Priority"`
	DishType    string `json:"DishType"`
	Portions    int    `json:"Portions"`
	TempMatch   string `json:"TempMatch"`
}

//Dishes type is a slice of the domain type Dish.
type Dishes []Dish

//Contains methods and validators that a dish would know about
//isExpired()
//get new dish with title()
//change expiration()
//how long remaining()
