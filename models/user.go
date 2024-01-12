package models

import "github.com/oklog/ulid"

type User struct {
	ID   ulid.ULID
	Name string
	Age  int
}

func (u User) Compare(otherUser User) bool {
	return u.ID.Compare(otherUser.ID) == 0 && u.Name == otherUser.Name && u.Age == otherUser.Age
}
