package pg

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestNewPG(t *testing.T) {
	godotenv.Load("../.env")
	_, err := NewPG()
	if err != nil {
		t.Fatalf("Error occured during open of new db\n%s", err)
	}
}
