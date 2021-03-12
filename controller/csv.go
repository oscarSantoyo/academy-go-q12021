package controller

import (
	"github.com/golobby/container"
	"github.com/oscarSantoyo/academy-go-q12021/service"
)

var csvService service.CsvService

func LoadCsvData(topic string) error {
	return getCsvService().DownloadCsvData(topic)
}

func getCsvService () service.CsvService {
	if csvService == nil {
		container.Make(&csvService)
	}
	return csvService
}
