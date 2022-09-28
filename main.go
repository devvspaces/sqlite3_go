package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "test.db"

type Activities struct {
	db                     *sql.DB
	insert, retrieve, list *sql.Stmt
}

type Activity struct {
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
	ID          uint64    `json:"id"`
}

const create string = `
CREATE TABLE IF NOT EXISTS activities (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  description TEXT
  );
`

func NewActivities() (*Activities, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	// Create stmts too
	insert, err := db.Prepare("INSERT INTO activities VALUES(NULL,?,?)")
	if err != nil {
		return nil, err
	}

	retrieve, err := db.Prepare("SELECT id, time, description FROM activities WHERE id=?")
	if err != nil {
		return nil, err
	}

	list, err := db.Prepare("SELECT id, time, description FROM activities")
	if err != nil {
		return nil, err
	}

	return &Activities{
		db:       db,
		insert:   insert,
		retrieve: retrieve,
		list:     list,
	}, nil
}

func (c *Activities) InsertActivity(activity Activity) (int, error) {
	result, err := c.insert.Exec(activity.Time, activity.Description)
	if err != nil {
		return 0, err
	}
	defer c.insert.Close()

	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *Activities) RetrieveActivity(id int) (Activity, error) {
	row := c.retrieve.QueryRow(id)
	defer c.retrieve.Close()

	act := Activity{}
	if err := row.Scan(&act.ID, &act.Time, &act.Description); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Id %v not found", id)
		}
		return Activity{}, err
	}

	return act, nil
}

func (c *Activities) List() ([]Activity, error) {

	rows, err := c.list.Query()
	if err != nil {
		return nil, err
	}
	defer c.list.Close()
	defer rows.Close()

	activities := []Activity{}
	for rows.Next() {
		activity := Activity{}
		err = rows.Scan(&activity.ID, &activity.Time, &activity.Description)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func main() {

	conn, err := NewActivities()
	if err != nil {
		log.Fatalf("Error occured: \n %v", err)
	}
	defer conn.db.Close()

	// Create new activity
	// activity := Activity{
	// 	Time:        time.Now(),
	// 	Description: "This is a test description",
	// }

	// id, err := conn.InsertActivity(activity)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("New id", id)

	// Get an activity
	// activity, err := conn.RetrieveActivity(4)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("The activity is %v", activity)

	// List activities
	// activities, err := conn.List()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for idx := range activities {
	// 	activity := activities[idx]
	// 	fmt.Printf("ID: %v,\tDescription: %v,\tTime: %v\n", activity.ID, activity.Description[:20], activity.Time)
	// }

}
