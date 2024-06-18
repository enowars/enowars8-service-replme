package server

func Init(apiKey string) {
	engine := NewRouter(apiKey)
	engine.Run(":3000")
}
