package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type Memory struct {
	Id           uint64    `db:"id" form:"id"`
	Title        string    `db:"title" form:"title"`
	Details      string    `db:"details" form:"details"`
	Latitude     float64   `db:"latitude" form:"latitude"`
	Longitude    float64   `db:"longitude" form:"longitude"`
	Author       string    `db:"author" form:"author"`
	IsApproved   bool      `db:"is_approved"`
	ApprovalUuid string    `db:"approval_uuid"`
	InsertedAt   time.Time `db:"inserted_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

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

func GetRecentMemories() ([]*Memory, error) {
	var memories []*Memory

	err := DB.Select(&memories, "select * from memory where is_approved = true order by id desc limit 5")

	return memories, err
}

func GetAllMemories() ([]*Memory, error) {
	var memories []*Memory

	err := DB.Select(&memories, "select * from memory where is_approved = true order by id desc")	

	return memories, err
}

func GetMemory(id int) (Memory, error) {
	memory := Memory{}
	
	err := DB.Get(&memory, "select * from memory where id = $1", id)

	return memory, err
}

func GetMemoryByUuid(uuid string) (Memory, error) {
	memory := Memory{}

	err := DB.Get(&memory, "select * from memory where is_approved = false and approval_uuid = $1", uuid)

	return memory, err
}

func AddMemory(memory *Memory) (uint64, error) {
	return NamedInsert("insert into memory values (default, :title, :details, :latitude, :longitude, :author, false, :approval_uuid, now(), now()) returning id", memory)
}

func ApproveMemory(memory Memory) error {
	_, err := DB.NamedExec("update memory set is_approved = true where id = :id", memory)
	
	return err
}
