package models

import (
	"errors"
	"log"
	"time"

	"example.com/rest-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

//var events = []Event{}

func (e *Event) Save() error {
	if db.DB == nil {
		return errors.New("database connection is nil")
	}
	query :=
		`INSERT INTO events 
	(name, description, location, dateTime, user_id)
	VALUES (?,?,?,?,?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted id: %v", err)
		return err
	}
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			log.Printf("Error scanning rows: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	if db.DB == nil {
		return nil, errors.New("database connection is nil")
	}
	query := "SELECT * FROM events WHERE id =?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		log.Printf("Error scanning row: %v", err)
		return nil, err
	}
	return &event, nil
}

func (e Event) Update() error {
	if db.DB == nil {
		return errors.New("database connection is nil")
	}
	query := `UPDATE events 
	SET name =?, description =?, 
	location =?, dateTime =?
	WHERE id =?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing updateQuery: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	if err != nil {
		log.Printf("Error executing updateQuery: %v", err)
		return err
	}
	return err
}

func (e Event) Delete() error {
	if db.DB == nil {
		return errors.New("database connection is nil")
	}
	query := "DELETE FROM events WHERE id =?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing deleteQuery: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID)
	return err
}

func (e Event) RegisterUser(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?,?)"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing registrationQuery: %v", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id =? AND user_id =?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing cancelRegistrationQuery: %v", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err

}
