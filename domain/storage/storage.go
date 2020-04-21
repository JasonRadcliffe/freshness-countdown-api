package storage

type storage struct {
	StorageID   int    `json:"StorageID"`
	UserID      int    `json:"UserID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type storages []storage
