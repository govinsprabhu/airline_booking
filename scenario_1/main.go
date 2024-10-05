package main

import (
	"airline_ticket/common"
	"database/sql"
	"log"
	"math/rand"
)

func Insert(db *sql.DB, passenger_id int, trip_id int) {
	seat_id := rand.Intn(120) + 1
	query := `UPDATE seats SET passenger_id = ? where  id = ?`
	_, err := db.Exec(query, passenger_id, seat_id)

	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	common.Prepare()
	db := common.GetDb()
	for i := 1; i <= 120; i++ {
		Insert(db, i, 1)
	}

	defer db.Close()

}
