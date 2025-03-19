package integ

import (
	"testing"
	"time"
)

func TestCreateRecord(t *testing.T) {
	testRecord := TestEntity{
		Column1: 12,
		Column2: 3.14,
		Column3: "Hello, World!",
		Column4: time.Now(),
		Column5: true,
	}

	result := db.Create(&testRecord)
	if result.Error != nil {
		t.Fatalf("Failed to create record: %v", result.Error)
	}
}

func TestReadRecord(t *testing.T) {
	var readRecord TestEntity
	result := db.First(&readRecord, "column3 = ?", "Hello, World!")
	if result.Error != nil {
		t.Fatalf("Failed to read record: %v", result.Error)
	}
	if readRecord.Column2 != 3.14 {
		t.Errorf("Expected Column2 to be 3.14, got %v", readRecord.Column2)
	}
	if readRecord.Column3 != "Hello, World!" {
		t.Errorf("Expected Column3 to be 'Hello, World!', got %v", readRecord.Column3)
	}
}

func TestUpdateRecord(t *testing.T) {
	var readRecord TestEntity
	db.Where("column3 = ?", "Hello, World!").First(&readRecord)

	readRecord.Column3 = "Updated Value"
	result := db.Where("column3 = ?", "Hello, World!").Save(&readRecord)
	if result.Error != nil {
		t.Fatalf("Failed to update record: %v", result.Error)
	}

	var updatedRecord TestEntity
	db.First(&updatedRecord, "column1 = ?", readRecord.Column1)
	if updatedRecord.Column3 != "Updated Value" {
		t.Errorf("Expected Column3 to be 'Updated Value', got %v", updatedRecord.Column3)
	}
}

func TestDeleteRecord(t *testing.T) {
	var readRecord TestEntity
	db.Where("column3 = ?", "Updated Value").First(&readRecord)

	result := db.Where("column3 = ?", "Updated Value").Delete(&readRecord)
	if result.Error != nil {
		t.Fatalf("Failed to delete record: %v", result.Error)
	}
}
