package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type Storage struct {
	db sql.DB
}

func (storage *Storage) GetLastUpdateId() (id int) {
	row := storage.db.QueryRow("SELECT id FROM `update` ORDER BY created DESC LIMIT 1")
	var updateId string
	scanErr := row.Scan(&updateId)
	if scanErr != nil && scanErr.Error() == "sql: no rows in result set" {
		return 0
	}
	fmt.Println(updateId)
	checkErr(scanErr)
	i, scanErr := strconv.Atoi(updateId)
	checkErr(scanErr)
	return i
}

func (storage *Storage) SaveUpdateId(updateId int) {
	stmt, e := storage.db.Prepare("INSERT INTO `update` (`id`) VALUES (?)")
	checkErr(e)
	stmt.Exec(updateId)
}

func (storage *Storage) SaveChatId(chatId int64) {
	stmt, e := storage.db.Prepare("INSERT INTO `chat` (`id`) VALUES (?)")
	checkErr(e)
	stmt.Exec(chatId)
}

func (storage *Storage) GetChatIds() (ids []int64) {
	rows, err := storage.db.Query("SELECT id FROM `chat`")
	if err != nil && err.Error() == "sql: no rows in result set" {
		return []int64{}
	}
	checkErr(err)

	var result []int64
	for rows.Next() {
		var chatId int64
		err = rows.Scan(&chatId)
		checkErr(err)
		result = append(result, chatId)
	}
	return result
}

func (storage *Storage) createSchema() {
	storage.createChatTable()
	storage.createUpdateTable()
}

func (storage *Storage) createUpdateTable() {
	row := storage.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='update';")
	var updateTable string
	scanErr := row.Scan(&updateTable)

	if scanErr == nil {
		return
	}
	if scanErr.Error() != "sql: no rows in result set" {
		log.Panic(scanErr)
	}

	sql := "CREATE TABLE `update` (`id` VARCHAR(64) PRIMARY KEY, `created` TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"
	stmt, e := storage.db.Prepare(sql)
	checkErr(e)
	_, e2 := stmt.Exec()
	checkErr(e2)
}

func (storage *Storage) createChatTable() {
	row := storage.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='chat';")
	var chatTable string
	scanErr := row.Scan(&chatTable)

	if scanErr == nil {
		return
	}
	fmt.Println(scanErr.Error())
	if scanErr.Error() != "sql: no rows in result set" {
		log.Panic(scanErr)
	}

	sql := "CREATE TABLE `chat` (`id` VARCHAR(64) PRIMARY KEY)"
	stmt, e := storage.db.Prepare(sql)
	checkErr(e)
	_, e2 := stmt.Exec()
	checkErr(e2)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NewStorage(path string) (storage Storage) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Panic(err)
	}
	storage = Storage{db: *db}
	storage.createSchema()
	return storage
}
