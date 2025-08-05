package main

import (
	db "jwtauth/database"
)

func main() {
	db.MongoConnect()
}
