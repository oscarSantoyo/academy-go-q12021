package service

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/oscarSantoyo/academy-go-q12021/model"

	"github.com/fatih/structs"
	"github.com/golobby/container"
	"github.com/labstack/gommon/log"
)

// CsvService is used as an interface for this service.
type CsvService interface {
	FilterByID(string) ([]model.Doc, error)
	DownloadCsvData(string) error
	SearchByConditions(map[string]string) ([]model.Doc, error)
}

// CsvServiceImpl is used as the implementation of this service.
type CsvServiceImpl struct{}

// DownloadCsvData downloads information from external API.
func (c CsvServiceImpl) DownloadCsvData(topic string) error {
	return getAndSaveData(topic)
}

// FilterByID reads records from a CSV and returns the filtered ones
func (c CsvServiceImpl) SearchByConditions(conditions map[string]string) ([]model.Doc, error) {
	file, err := os.Open(getConfigService().GetConfig().CSV.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	reader := csv.NewReader(file)
	records, errReader := reader.ReadAll()
	if errReader != nil {
		return nil, errReader
	}

	var recordsStruct []model.Doc
	for _, record := range records {

		docRecord := toDoc(record)
		fields := structs.Fields(docRecord)

		if ( filterRecord(fields, conditions) ) {
			recordsStruct = append(recordsStruct, toDoc(record))
		}

	}

	if len(records) == 0 {
		log.Info("The file is empty")
		return nil, errors.New("The file is empty")
	}
	defer file.Close()

	return recordsStruct, nil
}

// FilterByID reads records from a CSV and returns the filtered ones
func (c CsvServiceImpl) FilterByID(id string) ([]model.Doc, error) {
	file, err := os.Open(getConfigService().GetConfig().CSV.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	reader := csv.NewReader(file)
	records, errReader := reader.ReadAll()
	if errReader != nil {
		return nil, errReader
	}

	var recordsStruct []model.Doc
	conditions := make(map[string] string)
	conditions["Key"] = id

	for _, record := range records {

		// if record[0] != id {
			// continue
		// }
		docRecord := toDoc(record)
		fields := structs.Fields(docRecord)

		if ( filterRecord(fields, conditions) ) {
			recordsStruct = append(recordsStruct, toDoc(record))
		}

	}

	if len(records) == 0 {
		log.Info("The file is empty")
		return nil, errors.New("The file is empty")
	}
	defer file.Close()

	return recordsStruct, nil
}

func toDoc(record []string) model.Doc {
	return model.Doc{
		Key:       record[0],
		Title:     record[1],
		Type:      record[2],
		Published: record[3],
	}
}

func filterRecord(s []*structs.Field, m map[string]string) bool {
	for _, f := range s {
		log.Info("validating "+f.Name()+" in conditions "+m[f.Name()])
		if filter := m[f.Name()]; filter != "" {
			if f.Value().(string) != filter {
				return false
			}
		}

	}
	return true
}

func deleteCsv() {
	err := os.Remove(getConfigService().GetConfig().CSV.FileName)
	if err != nil {
		log.Error("Cannot delete file", err)
	}
}

func getAndSaveData(topic string) error {
	// the+lord+of+the+rings&page=1
	// <search-term>  <page>

	deleteCsv() // before write results clears file
	escapedTopic := strings.ReplaceAll(topic, 	" ", "+")
	apiUrl := strings.ReplaceAll(getConfigService().GetConfig().External.ApiUrl, "<topic>", escapedTopic)
	file, err := openOrCreate(getConfigService().GetConfig().CSV.FileName)

	if err != nil {
		log.Error("Unable to use file", err)
		return err
	}

	for i := 1; true; i++ {
		response, err := http.Get(strings.ReplaceAll(apiUrl, "<page>", strconv.Itoa(i)))

		if err != nil {
			log.Info(err.Error())
			return err
		}


		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Error(err)
			return err
		}

		var responseObject model.SearchResponse

		json.Unmarshal(responseData, &responseObject)

		if (len(responseObject.Doc) == 0) {
			break
		}

		saveCsv(responseObject, *file)
	}

	defer file.Close()

	return nil
}

func saveCsv(responseObject model.SearchResponse, file os.File) {
	writer := csv.NewWriter(&file)
	for _, doc := range responseObject.Doc {
		str := interfaceToString(structs.Values(doc))
		err := writer.Write(str)
		if err != nil {
			log.Error("cannot write CSV", err)
		}
	}
	defer writer.Flush()
}
func openOrCreate(path string) (*os.File, error) {
	file, _ := os.Open(path)

	if file != nil {
		log.Info("File opened")
		return file, nil
	}

	file, err := os.Create(path)

	if err != nil {
		log.Error("File was not able to create")
		return nil, err
	}
	return file, nil
}

func interfaceToString(record []interface{}) []string {
	a := make([]string, len(record))

	for i, row := range record {
		a[i] = row.(string)
	}
	return a
}

func getConfigService() ConfigService {
	var config ConfigService
	container.Make(&config)
	return config
}
