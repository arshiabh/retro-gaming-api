package main

import "log"

func main() {
	config := Load()

	//faster with pointer
	app := &application{
		config: *config,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
