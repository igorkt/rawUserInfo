package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	const numberOfUsers = 100
	const workerCount = 3

	startTime := time.Now()

	userCh := make(chan User, numberOfUsers)
	// userFiles := make(chan User, numberOfUsers)

	wg := &sync.WaitGroup{}

	generateUsers(numberOfUsers, userCh)

	wg.Add(numberOfUsers)
	for i := 0; i < numberOfUsers; i++ {
		go saveUserInfo(<-userCh, wg)
	}

	wg.Wait()
	close(userCh)

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func worker(id int, users <-chan User, results chan<- User) {
	for user := range users {
		fmt.Printf("worker #%d finished\n", id)
		results <- user
	}
}

func saveUserInfo(user User, wg *sync.WaitGroup) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
	wg.Done()
}

func generateUsers(count int, userCh chan User) {
	for i := 0; i < count; i++ {
		go generateUser(i, userCh)
	}
}

func generateUser(i int, userCh chan User) {
	userCh <- User{
		id:    i + 1,
		email: fmt.Sprintf("user%d@company.com", i+1),
		logs:  generateLogs(rand.Intn(1000)),
	}
	fmt.Printf("generated user %d\n", i+1)
	time.Sleep(time.Millisecond * 100)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
