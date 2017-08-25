package db

import (
	"github.com/jinzhu/gorm"
	. "../dataobjects"
	"time"
)

/*
	Data that will be used to initialise the database in the AutoMigration
 */
 var (
	 users  	[]User
	 countries  []Country
	 cities 	[]City
	 dbInitialised bool = false
 )

func InitialiseDB(db *gorm.DB) {
	//Test if there is data in the database
	checkDatabaseEmpty(db)
	if dbInitialised {
		/*
			COUNTRIES
	 	*/
		insertCountries(db)
		/*
			SERVICES
		*/
		insertCities(db)
		/*
			PRODUCTS
		*/
		insertUsers(db)
	}
}

func checkDatabaseEmpty(db *gorm.DB){
	countUsers, countCountries, countCities := 0,0,0
	db.Find(&users).Count(&countUsers)
	db.Find(&countries).Count(&countCountries)
	db.Find(&cities).Count(&countCities)

	if countUsers == 0 && countCountries == 0 && countCities == 0  {
		dbInitialised = true
	}
}

func insertCountries(db *gorm.DB){
	var query = "INSERT INTO countries VALUES(?,?)"
	/*
		SOUTH AMERICA
	 */
	db.Exec(query,1,"ARGENTINA")
	db.Exec(query,2,"BOLIVIA")
	db.Exec(query,3,"BRAZIL")
	db.Exec(query,4,"CHILE")
	db.Exec(query,5,"COLOMBIA")
	db.Exec(query,6,"ECUADOR")
	db.Exec(query,7,"FRENCH GUIANA")
	db.Exec(query,8,"GUYANA")
	db.Exec(query,9,"PARAGUAY")
	db.Exec(query,10,"PERU")
	db.Exec(query,11,"SURINAME")
	db.Exec(query,12,"URUGUAY")
	db.Exec(query,13,"VENEZUELA")
	/*
		CENTRAL AMERICA AND NORTH AMERICA
	 */
	db.Exec(query,14,"BELIZE")
	db.Exec(query,15,"COSTA RICA")
	db.Exec(query,16,"EL SALVADOR")
	db.Exec(query,17,"GUATEMALA")
	db.Exec(query,18,"HONDURAS")
	db.Exec(query,19,"MEXICO")
	db.Exec(query,20,"NICARAGUA")
	db.Exec(query,21,"PANAMA")
	/*
		OCEANIA
	 */
	db.Exec(query,22,"AUSTRALIA")
	db.Exec(query,23,"NEW ZELAND")
}

func insertCities(db *gorm.DB){
	var query = "INSERT INTO cities VALUES(?,?,?)"
	/*
		TYPE OF SERVICES
	 */
	db.Exec(query,1,"PERTH",22)
	db.Exec(query,2,"SIDNEY",22)
	db.Exec(query,3,"MELBOURNE",22)
	db.Exec(query,4,"BRISBANE",22)
	db.Exec(query,5,"ADELAIDE",22)
}

func insertUsers(db *gorm.DB){
	var query = "INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?)"
	/*
		TYPE OF PRODUCTS
	 */
	db.Exec(query,1,"admin","capcarde@gmail.com","PASSWORD","ADMIN","true",time.Now(),time.Now(),22,1)
}
