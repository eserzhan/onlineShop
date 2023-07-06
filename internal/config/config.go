package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)


type (
	Config struct {
		Environment string
		HTTP        HTTPConfig
		Postgres 	DB
		Auth        AuthConfig
	}



	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verificationCodeLength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}
	
	DB struct{
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
		Dbname string `mapstructure:"dbname"`
		Sslmode string `mapstructure:"sslmode"`
		Username string `mapstructure:"username"`
		Password string
	}
	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

)



func InitConfig() (*Config, error) {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("main")
	// viper.AddConfigPath("/app/configs")
	
	// viper.SetConfigFile("main")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err 
	}

	var cfg Config
	if err := viper.UnmarshalKey("db", &cfg.Postgres); err != nil {
		return nil, err 
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return nil, err 
	}


	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("auth.verificationCodeLength", &cfg.Auth.VerificationCodeLength); err != nil {
		return nil, err
	}


	if err := parseEnv(&cfg); err != nil {
		return nil, err 
	}

	return &cfg, nil 
}

func parseEnv(cfg *Config) error {
	err := godotenv.Load("../.env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	if err := viper.BindEnv("PASSWORD"); err != nil {
		return err 
	}

	if err := viper.BindEnv("JWT.SigningKey"); err != nil {
		return err 
	}

	if err := viper.BindEnv("Password.salt"); err != nil {
		return err 
	}

	cfg.Postgres.Password = viper.GetString("PASSWORD")
	cfg.Auth.JWT.SigningKey = viper.GetString("JWT.SigningKey")
	cfg.Auth.PasswordSalt = viper.GetString("Password.salt")

	return nil 
}