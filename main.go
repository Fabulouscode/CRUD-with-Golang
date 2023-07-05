package main

type Address struct {
	Street string `json:"street" "bson:street"`
	City   string `json:"city" "bson:city"`
}

type User struct {
	Name    string  `json:"name" "bson:name"`
	Age     int     `json:"age" "bson:age"`
	Address Address `json:"address "bson:address""`
}

func main() {}
