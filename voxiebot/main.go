package main

func main() {
	account1 := NewInstance("aeooedysn", "02f8E9B288a12")
	if err := account1.Authorize(); err != nil {
		panic(err)
	}

	account1.StartMatchmaking()
}