package global

type Config struct {
	App struct {
		Port string
	}

	Object struct {
		Minio struct {
			Endpoint  string
			AccessKey string
			SecretKey string
			Token     string
			UseSsl    bool
		}
	}
}
