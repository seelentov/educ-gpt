package unit

import (
	"educ-gpt/config/data"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func setup() error {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print(err)
		return err
	}

	if err := data.SwitchToMock(); err != nil {
		log.Print(err)
		return err
	}

	if err := data.SwitchToMockRedis(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func tearDown() error {
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := tearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}
