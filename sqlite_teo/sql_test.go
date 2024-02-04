package sqlite_teo

import "testing"

func TestVideo(t *testing.T) {

	// todo!
	value := "Mateo_ortega"
	_, err := CreateDir(&value)

	if err != nil {
		t.Errorf("Could not create directory %v", err)
	}
	value = ("nul")
	_, err = CreateDir(&value)

	if err != nil {
		t.Fail()
	}

}
