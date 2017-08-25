package dataobjects

import "time"

type User struct {
	Id		int		`gorm:"primary_key; AUTO_INCREMENT" json:"id"`
	Username	string		`gorm:"size:255" json:"username"`
	Email		string		`gorm:"size:100" json:"email"`
	Password	string		`json:"password"`
	Role 		string		`json:"role"`
	IsActivated 	string		`json:"activated"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt 	time.Time	`json:"updatedAt"`
	Country		Country		`gorm:"ForeignKey:Id;AssociationForeignKey:CountryId" json:"country"` //ISO "ALPHA-2 Code"
	CountryId	int		`json:"country_id,omitempty"`
	City		Country		`gorm:"ForeignKey:Id;AssociationForeignKey:CityId" json:"city"` //ISO "ALPHA-2 Code"
	CityId		int		`json:"city_id,omitempty"`
}

type Country struct {
	Id	int	`gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name	string	`gorm:"size:150" json:"name"`
}


type City struct {
	Id		int		`gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name		string		`gorm:"size:150"`
	CountryId	int		`json:"countryId"`
}

