package db

import (
	"bufio"
	"log"
	"os"

	"github.com/gbandres98/pack-and-go/model"
)

type fileDB struct {
	filePath string
}

func NewFileDB(filePath string) *fileDB {
	return &fileDB{filePath: filePath}
}

func (fileDB *fileDB) GetAllCities() ([]model.City, error) {
	file, err := os.Open(fileDB.filePath)
	if err != nil {
		log.Panicf("could not open cities db file: %v", err)
	}
	defer file.Close()

	result := []model.City{}
	nextId := int32(1)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, model.City{Id: nextId, Name: scanner.Text()})
		nextId++
	}

	return result, scanner.Err()
}

func (fileDB *fileDB) GetCityById(id int32) (model.City, error) {
	cities, err := fileDB.GetAllCities()
	if err != nil {
		return model.City{}, err
	}

	for _, city := range cities {
		if (city.Id == id) {
			return city, nil
		}
	}

	return model.City{}, ErrorCityNotFound
}