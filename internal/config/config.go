package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var flags = pflag.NewFlagSet("flags", pflag.ExitOnError)

func init() {
	println("Start to cfg...")
	flags.StringP("ServerAddress", "s", ":9000", "ServerAddress")
	flags.StringP("PostgresUrl", "d", "postgres://postgres:password@localhost:5432/users?sslmode=disable", "PostgresUrl")
	flags.StringP("KafkaBrokers", "k", "localhost:9095", "KafkaBrokers")
	flags.StringP("ProducerTopicID", "p", "ProduceTopic", "ProducerTopicID")
	flags.StringP("ConsumerTopicID", "c", "ConsumeTopic", "ConsumerTopicID")
	flags.StringP("KafkaGroup", "g", "gorik-group", "KafkaGroup")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing flags: %v\n", err)
	}

	bindViperPFlag("ServerAddress")
	bindViperPFlag("PostgresUrl")
	bindViperPFlag("KafkaBrokers")
	bindViperPFlag("ProducerTopicID")
	bindViperPFlag("ConsumerTopicID")
	bindViperPFlag("KafkaGroup")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindViperEnv("ServerAddress", "SERVER_ADDRESS")
	bindViperEnv("PostgresUrl", "POSTGRES_URL")
	bindViperEnv("KafkaBrokers", "KAFKA_BROKERS")
	bindViperEnv("ProducerTopicID", "PRODUCER_TOPICID")
	bindViperEnv("ConsumerTopicID", "CONSUMER_TOPICID")
	bindViperEnv("KafkaGroup", "KAFKA_GROUP")

}

func bindViperPFlag(flagName string) {
	err := viper.BindPFlag(flagName, flags.Lookup(flagName))
	if err != nil {

		log.Printf("Error parsing flags: %v\n", err)

	}
}

func bindViperEnv(viperKey, envName string) {

	err := viper.BindEnv(viperKey, envName)
	if err != nil {
		log.Printf("Error parsing env vars: %v\n", err)

	}

}
func In() {}
