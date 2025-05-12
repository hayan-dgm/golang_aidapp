package models

type Family struct {
	ID            int    `json:"id"`
	FullName      string `json:"fullName"`
	NationalID    string `json:"nationalID"`
	FamilyBookID  string `json:"familyBookID"`
	PhoneNumber   string `json:"phoneNumber"`
	FamilyMembers int    `json:"familyMembers"`
	Children      int    `json:"children"`
	Babies        int    `json:"babies"`
	Adults        int    `json:"adults"`
	Milk          int    `json:"milk"`
	Diapers       int    `json:"diapers"`
	Basket        int    `json:"basket"`
	Clothing      int    `json:"clothing"`
	Drugs         int    `json:"drugs"`
	Other         string `json:"other"`
	Taken         bool   `json:"taken"`
}
