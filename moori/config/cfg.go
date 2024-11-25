package config

import (
	meiliCaller "moori/delivery/meilisearch"
	"moori/storage/mysql"
	qdrantDb "moori/storage/qdrant"
	"os"
	"strconv"
)

func New() (mysql.Config, qdrantDb.Config, meiliCaller.Config) {
	mysqlCfg := mysql.Config{
		Host:     os.Getenv("MYsqlHost"),
		Port:     os.Getenv("MYsqlPort"),
		Username: os.Getenv("MYsqlUsername"),
		Password: os.Getenv("MYsqlPassword"),
		Database: os.Getenv("MYsqlDatabase"),
	}

	qdrantPort, _ := strconv.Atoi(os.Getenv("QdrantPort"))
	QDRANTUseTLS, _ := strconv.ParseBool(os.Getenv("QDRANTUseTLS"))

	qdrantCfg := qdrantDb.Config{
		Host:   os.Getenv("QDRANTHost"),
		Port:   qdrantPort,
		ApiKey: os.Getenv("QDRANTApiKey"),
		UseTLS: QDRANTUseTLS,
	}

	meiliCfg := meiliCaller.Config{
		Address:  os.Getenv("MeiliAddress"),
		ApiKey:   os.Getenv("MeiliApiKey"),
		Index:    os.Getenv("MeiliIndex"),
		Document: nil,
	}
	return mysqlCfg, qdrantCfg, meiliCfg
}
