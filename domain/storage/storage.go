package storage

//Storage type is the struct in the Domain that contains all the fields for what a Storage Unit is.
type Storage struct {
	StorageID   int    `json:"StorageID"`
	UserID      int    `json:"UserID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	TempMatch   string `json:"TempMatch"`
}

//Storages type is a slice of the domain type Storage.
type Storages []Storage

//Contains methods and validators that a storage unit would know about
//...
//...
//... anything?
