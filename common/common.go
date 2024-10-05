package common

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser = "root"
	dbPass = "password"
	dbName = "travel_db"
)

var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Elizabeth", "David", "Barbara", "Richard", "Margaret", "Joseph", "Susan",
	"Thomas", "Dorothy", "Charles", "Lisa", "Christopher", "Nancy", "Daniel", "Karen",
	"Paul", "Betty", "Mark", "Helen", "Donald", "Sandra", "George", "Donna", "Kenneth",
	"Carol", "Steven", "Ruth", "Edward", "Sharon", "Brian", "Michelle", "Ronald", "Laura",
	"Anthony", "Sarah", "Kevin", "Kimberly", "Jason", "Deborah", "Jeff", "Jessica",
	"Gary", "Shirley", "Timothy", "Cynthia", "Jose", "Angela", "Larry", "Melissa", "Scott",
	"Brenda", "Frank", "Amy", "Eric", "Anna", "Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
	"Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
	"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young",
	"Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green",
	"Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts",
	"Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins",
	"Reyes", "Stewart", "Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz",
	"Morgan", "Cooper", "Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim",
	"Cox", "Ward", "Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett",
	"Gray", "Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel",
	"Myers", "Long", "Ross", "Foster", "Jimenez",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
	"Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
	"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young",
	"Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green",
	"Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts",
	"Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins",
	"Reyes", "Stewart", "Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz",
	"Morgan", "Cooper", "Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim",
	"Cox", "Ward", "Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett",
	"Gray", "Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel",
	"Myers", "Long", "Ross", "Foster", "Jimenez", "James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Elizabeth", "David", "Barbara", "Richard", "Margaret", "Joseph", "Susan",
	"Thomas", "Dorothy", "Charles", "Lisa", "Christopher", "Nancy", "Daniel", "Karen",
	"Paul", "Betty", "Mark", "Helen", "Donald", "Sandra", "George", "Donna", "Kenneth",
	"Carol", "Steven", "Ruth", "Edward", "Sharon", "Brian", "Michelle", "Ronald", "Laura",
	"Anthony", "Sarah", "Kevin", "Kimberly", "Jason", "Deborah", "Jeff", "Jessica",
	"Gary", "Shirley", "Timothy", "Cynthia", "Jose", "Angela", "Larry", "Melissa", "Scott",
	"Brenda", "Frank", "Amy", "Eric", "Anna",
}

func GetDb() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPass, dbName))

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Prepare(createNew bool) {
	db := GetDb()
	defer db.Close()
	if createNew {
		createTable(db)
		populatePassengers(db)
	}

}

func populatePassengers(db *sql.DB) {
	stmet, err := db.Prepare("INSERT INTO  passengers (first_name, last_name) VALUES (?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmet.Close()

	for i := 0; i < 120; i++ {
		firstName := firstNames[i]
		lastName := lastNames[i]
		result, err := stmet.Exec(firstName, lastName)

		if err != nil {
			log.Fatal(err)
		}
		_, err = result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("Last Inserted ID:", id)
	}

	fmt.Println("Populated 120 passengers with details")

	stmet, err = db.Prepare("INSERT INTO  seats (seat_number) VALUES (?)")

	if err != nil {
		log.Fatal(err)
	}

	seats := []string{"A", "B", "C", "D", "E", "F"}
	fmt.Println("Inserting the seats")
	for i := 1; i <= 20; i++ {
		for _, seat := range seats {
			seat_number := fmt.Sprintf("%d-%s", i, seat)
			//fmt.Println(seat_number)
			_, err = stmet.Exec(seat_number)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

}

func createTable(db *sql.DB) {
	tableNames := [3]string{"seats", "trip", "passengers"}
	for i := 0; i < 3; i++ {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableNames[i])
		result, err := db.Exec(query)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}
	fmt.Println("Deleting completed")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS passengers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			first_name VARCHAR(100),
			last_name VARCHAR(100)
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS trip (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100)
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS seats (
			id INT AUTO_INCREMENT PRIMARY KEY,
			seat_number VARCHAR(10),
			passenger_id INT,
			trip_id INT,
			FOREIGN KEY (passenger_id) REFERENCES passengers(id),
			FOREIGN KEY (trip_id) REFERENCES trip(id)
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
}
