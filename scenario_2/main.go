package main

import (
	"airline_ticket/common"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
)

func Insert(tx *sql.Tx, db *sql.DB, passenger_id int, trip_id int) error {
	seat_id := rand.Intn(120) + 1
	query := `UPDATE seats SET passenger_id = ? where  id = ?`
	_, err := tx.Exec(query, passenger_id, seat_id)

	fmt.Printf("Updating the seat with passenger: %d with seatid :%d", passenger_id, seat_id)
	return err
}

func main() {
	common.Prepare(true)
	db := common.GetDb()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	errCh := make(chan error, 120)
	var wg sync.WaitGroup

	for i := 1; i <= 120; i++ {
		wg.Add(1)
		//fmt.Print("updating for passenger : ", i)
		go func(db *sql.DB, passenger_id int, seat_id int) {
			defer wg.Done() //
			err := Insert(tx, db, passenger_id, seat_id)
			if err != nil {
				errCh <- fmt.Errorf("failed to update the seat: %d for passenger: %d:  %w", seat_id, passenger_id, err)
			}
		}(db, i, 1)
	}
	wg.Wait()

	close(errCh)
	if len(errCh) > 0 {
		log.Println("error occurred during transaction, rolling back")
		if err := tx.Rollback(); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}
	PrintLayout(db)
	log.Println("Transaction committed successfully!")
}

type Seat struct {
	seat_id      int
	passenger_id *int
}

func PrintLayout(db *sql.DB) {
	query := `select id, passenger_id from seats order by id`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var seats []Seat

	for rows.Next() {
		var seat Seat

		err := rows.Scan(&seat.seat_id, &seat.passenger_id)

		if err != nil {
			log.Fatal(err)
		}

		seats = append(seats, seat)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for index, seat := range seats {

		if seat.passenger_id == nil {
			fmt.Print("_ ")
		} else {
			//fmt.Println(*seat.passenger_id, seat.seat_id)
			fmt.Print("* ")
		}
		if index%20 == 0 {
			fmt.Println()
		}
		if index%60 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
