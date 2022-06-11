package cfg


import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

//структура для хранения конфигов, Viper для матчинга имен
type Cfg struct {
	Port   string
	DbName string
	DbUser string
	DbPass string
	DbHost string
	DbPort string
}

const name int = 344
const NatsURI string = "nats://localhost:4222"

func LoadAndStoreConfig() Cfg {
	v := viper.New() //создаем экземпляр ридера для Env
	v.SetEnvPrefix("SERV") // чтобы не создавать коллизий с переменными окружения которые уже есть
	v.SetDefault("PORT", "8080")
	v.SetDefault("DBUSER", "user")
	v.SetDefault("DBPASS", "example")
	v.SetDefault("DBHOST", "")
	v.SetDefault("DBPORT", "5400")
	v.SetDefault("DBNAME", "db01")
	v.AutomaticEnv() //собираем наши переменные с системных

	var cfg Cfg

	err := v.Unmarshal(&cfg) //закидываем переменные в cfg после анмаршалинга
	if err != nil {
		log.Panic(err)
	}
	return cfg
}

//метод для сборки строки соединения с БД
func (cfg *Cfg) GetDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName)
}

