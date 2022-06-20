package db

import (
	"testing"
)

func TestGetAllCities_1(t *testing.T) {
	db := NewFileDB("./cities_test.txt")

	cities, err := db.GetAllCities()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cities == nil {
		t.Fatalf("expected non-nil array")
	}

	for i, city := range cities {
		if (city.Id != int32(i + 1)) {
			t.Fatalf("expected city to have id: %v, got id: %v", i + 1, city.Id)
		}
	}
}

func TestGetAllCities_2(t *testing.T) {
	defer func() { recover() }()

	db := NewFileDB("wrong-file-path.txt")
	db.GetAllCities()

	t.Fatalf("expected GetAllCities() to panic")
}

func TestGetCityById_1(t *testing.T) {
	db := NewFileDB("./cities_test.txt")

	city, err := db.GetCityById(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if (city.Id != 3) {
		t.Fatalf("expected city to have id: %v, got %v", 3, city)
	}
}

func TestGetCityById_2(t *testing.T) {
	db := NewFileDB("./cities_test.txt")

	_, err := db.GetCityById(9)
	if err != ErrorCityNotFound {
		t.Fatalf("expected error: %v, got error: %v", ErrorCityNotFound, err)
	}
}