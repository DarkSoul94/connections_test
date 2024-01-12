package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	cryptorand "crypto/rand"

	"github.com/DarkSoul94/connections_test/models"
	"github.com/DarkSoul94/connections_test/repo"
	"github.com/oklog/ulid"
)

func main() {
	fmt.Println("=================")
	//test 1 - 1 db connect
	dbConn := repo.ConnectDB()
	repo.RunMigrations(dbConn)

	var wg1 sync.WaitGroup
	wg1.Add(2)

	fmt.Println("Test 1 start")
	test1Start := time.Now()

	go transactionFunc(dbConn, &wg1)
	go userFunc(dbConn, &wg1)
	wg1.Wait()

	test1End := time.Now()
	fmt.Println("Test 1 end")

	fmt.Printf("test 1 duration: %s\n", test1End.Sub(test1Start).String())

	dbConn.Close()

	fmt.Println("=================")
	//test 2 - individual db connections
	db1Conn := repo.ConnectDB()
	db2Conn := repo.ConnectDB()

	var wg2 sync.WaitGroup
	wg2.Add(2)

	fmt.Println("Test 2 start")
	test2Start := time.Now()

	go transactionFunc(db1Conn, &wg2)
	go userFunc(db2Conn, &wg2)

	wg2.Wait()

	test2End := time.Now()
	fmt.Println("Test 2 end")

	fmt.Printf("test 2 duration: %s\n", test2End.Sub(test2Start).String())
	fmt.Println("=================")
}

func NewULID() ulid.ULID {
	entropy := cryptorand.Reader
	t := time.Now().UTC()
	ent1 := ulid.Monotonic(entropy, 0)
	return ulid.MustNew(ulid.Timestamp(t), ent1)
}

func transactionFunc(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	tx.Commit()
}

func userFunc(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	newUser := models.User{
		ID:   NewULID(),
		Name: "Slava",
		Age:  29,
	}
	fmt.Printf("new user %v\n", newUser)

	err := repo.CreateUser(db, newUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("user created")

	existUser, err := repo.GetUser(db, newUser.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("user in db %v\n", existUser)

	if !existUser.Compare(newUser) {
		panic("users not compare")
	}
	fmt.Println("user is compare")

	err = repo.DeleteUser(db, newUser.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("user deleted")
}
