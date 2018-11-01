package main

import (
	"log"

	"github.com/Charlemagne3/golang-api-skeleton/server/config"
)

var conf *config.Configuration

func init() {
	conf = config.GetConfiguration()
}

func main() {
	log.Printf("%s %s \n", conf.AppName, conf.AppVersion)
	s := makeServer()
	s.ListenAndServe()
}
