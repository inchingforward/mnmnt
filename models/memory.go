package models

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

var (
	DB *sqlx.DB
	slugifyExpr = regexp.MustCompile("[^a-z0-9]+")
)

// A Memory represents a single memory that is tied to a location.
type Memory struct {
	ID           uint64    `db:"id" form:"id"`
	AddressText  string    `db:"address_text" form:"address_text"`
	Title        string    `db:"title" form:"title"`
	Slug         string    `db:"slug"`
	Details      string    `db:"details" form:"details"`
	Latitude     float64   `db:"latitude" form:"latitude"`
	Longitude    float64   `db:"longitude" form:"longitude"`
	Author       string    `db:"author" form:"author"`
	IsApproved   bool      `db:"is_approved"`
	ApprovalUUID string    `db:"approval_uuid"`
	EditUUID     string    `db:"edit_uuid"`
	InsertedAt   time.Time `db:"inserted_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// NamedInsert executes the query insert statement and returns the
// generated sequence id.
func NamedInsert(query string, arg interface{}) (uint64, error) {
	rows, err := DB.NamedQuery(query, arg)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, rows.Err()
	}

	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetRecentMemories returns the 5 most recent memories that are approved.
func GetRecentMemories() ([]*Memory, error) {
	var memories []*Memory

	err := DB.Select(&memories, "select * from memory where is_approved = true order by id desc limit 10")

	return memories, err
}

// GetAllMemories returns all approved memories.
func GetAllMemories() ([]*Memory, error) {
	var memories []*Memory

	err := DB.Select(&memories, "select * from memory where is_approved = true order by id desc")

	return memories, err
}

// GetMemoryByID returns an individual memory by Memory ID.  The memory will not be returned
// if it is not approved.
func GetMemoryByID(id int) (Memory, error) {
	memory := Memory{}

	err := DB.Get(&memory, "select * from memory where id = $1 and is_approved = true", id)

	return memory, err
}

// GetMemoryByEditUUID returns an individual memory by its edit UUID.  The memory will
// not be returned if it is not approved.
func GetMemoryByEditUUID(uuid string) (Memory, error) {
	memory := Memory{}

	err := DB.Get(&memory, "select * from memory where is_approved = true and edit_uuid = $1", uuid)

	return memory, err
}

// GetMemoryByApprovalUUID returns an indiviual memory by its approval UUID.  The memory must
// not already be approved.
func GetMemoryByApprovalUUID(uuid string) (Memory, error) {
	memory := Memory{}

	err := DB.Get(&memory, "select * from memory where is_approved = false and approval_uuid = $1", uuid)

	return memory, err
}

// AddMemory inserts memory.  The sequence ID generated by the database is set on memory.
func AddMemory(memory *Memory) error {
	if memory.AddressText == "" {
		return errors.New("A place or address is required.")
	}

	if memory.Title == "" {
		return errors.New("A Title is required.")
	}

	if memory.Details == "" {
		return errors.New("Memory details are required.")
	}

	if memory.Latitude == 0 && memory.Longitude == 0 {
		return errors.New("A valid memory location is required.")
	}

	if memory.Author == "" {
		memory.Author = "Anonymous"
	}

	memory.ApprovalUUID = uuid.NewV4().String()
	memory.EditUUID = uuid.NewV4().String()
	memory.Slug = slugify(memory.Title)

	id, err := NamedInsert("insert into memory values (default, :title, :details, :latitude, :longitude, :author, false, :approval_uuid, now(), now(), :edit_uuid, :address_text) returning id", memory)
	if err != nil {
		return err
	}

	memory.ID = id

	log.Printf("New memory \"%v\" (id: %v) created.\n", memory.Title, memory.ID)

	return nil
}

// UpdateDetails updates the details and update date for the given memory.
func UpdateDetails(memory Memory) error {
	_, err := DB.NamedExec("update memory set details = :details, updated_at = now() where id = :id", memory)

	return err
}

// ApproveMemory sets memory as approved.
func ApproveMemory(memory Memory) error {
	_, err := DB.NamedExec("update memory set is_approved = true where id = :id", memory)

	return err
}

func slugify(str string) string {
	// FIXME: check the database for duplicates, consider passing in the memory 
	// and setting the value on the memory
    return strings.Trim(slugifyExpr.ReplaceAllString(strings.ToLower(str), "-"), "-")
}