package sqlite_teo

import (
	"testing"
)

func TestVideo(t *testing.T) {

	// todo!
	value := "Mateo_ortega"
	err := CreateDir(&value)

	if err != nil {
		t.Errorf("Could not create directory %v", err)
	}
	value = ("nul")
	err = CreateDir(&value)

	if err != nil {
		t.Fail()
	}

}
