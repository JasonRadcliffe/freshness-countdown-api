package dishes

type dish struct {
	DishID      int    `json:"DishID"`
	UserID      int    `json:"UserID"`
	StorageID   int    `json:"StorageID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	CreatedDate string `json:"TimeCreated"`
	ExpireDate  string `json:"TimeExpires"`
	Priority    string `json:"Priority"`
}

type dishes []dish
