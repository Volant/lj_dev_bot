package rating

import (
	"database/sql"
	"log"
)

func GetRating(db_file string, username string) (int, error) {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select rating from users where username = ?", username)
	if err != nil {
		log.Printf("%s", err)
		return 0, err
	}
	defer rows.Close()

	var rating int = 1
	if rows.Next() {
		rows.Scan(&rating)
		log.Printf("was: %d", rating)
		rating++
		log.Printf("after: %d", rating)
		rows.Close()
		_, err := db.Exec("update users set rating = ? where username = ?", rating, username)
		if err != nil {
			log.Print(err)
		}
	} else {
		db.Exec("insert into users (rating, username) values (?, ?)", rating, username)
	}

	return rating, nil
}
