package config 

import(
	"os"
	"github.com/joho/godotenv"
)
type ConfigEnvs struct{
	SecretJwt string
	Port string
	GoogleClientId string
	GoogleClientSecret string

	Host string
	DbUser string
	DbPassword string
	DbPort string
	DbName string 
	SslMode string
	ReplicateKey string
	BackendServerBaseUrl string
	ReplicateWebhookSecret string

	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2PublicURL       string

	Enviroment	string
}

var Envs *ConfigEnvs
func LoadEnvs()  {
	// Intentamos cargar el archivo .env local.
	// Usamos "_" para ignorar el error, ya que en Railway/Prod 
	// el archivo no existirá y no queremos que la app se detenga.
	_ = godotenv.Load()

	Envs = &ConfigEnvs{
		SecretJwt: getEnv("JWT_SECRET", "default_secret_key"),
		Port:      getEnv("PORT", "8080"),
		GoogleClientId: getEnv("GOOGLE_CLIENT_ID", "default"), 
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET","default"),

		Enviroment: getEnv("ENVIROMENT", "development"),

		Host:       getEnv("HOST", "localhost"),
		DbUser:     getEnv("DB_USER", "postgres"),
		DbPassword: getEnv("DB_PASSWORD", ""),
		DbName:     getEnv("DB_NAME", "postgres"),
		DbPort:     getEnv("DB_PORT", "5432"),
		SslMode:    getEnv("DB_SSLMODE", "disable"),

		ReplicateKey: getEnv("REPLICATE_KEY","default"),
		BackendServerBaseUrl: getEnv("BACKEND_SERVER_URL","default"), 

		ReplicateWebhookSecret:getEnv("REPLICATE_WEBHOOK_SECRET","default"),

		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:      getEnv("R2_BUCKET_NAME", ""),
		R2PublicURL:       getEnv("R2_PUBLIC_URL", ""),

	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}