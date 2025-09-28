package main

func main() {
	app := InitAPP()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
