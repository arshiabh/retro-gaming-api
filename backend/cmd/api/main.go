package main

func main() {
	app := application{}
	max := app.mount()
	app.run(max)
}
