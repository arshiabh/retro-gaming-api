package main

func main() {
	config := Load()
	app := application{
		cfg: *config,
	}
	max := app.mount()
	app.run(max)
}
