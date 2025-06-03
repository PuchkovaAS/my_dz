package main

func main() {
	conf := configs.LoadConfig()

	dataBase := db.NewDb(conf)
	dataBase.AutoMigrate(&link.Link{})
}
