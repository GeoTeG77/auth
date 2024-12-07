package service

import (
	"auth/internal/models"
	"auth/internal/repository"
	"auth/storage"
	"flag"
	"fmt"
	"os"
	"testing"
)

var testCount int

func init() {
	flag.IntVar(&testCount, "count", 1, "Number of test numbers")
}

func SetupTestEnv(t *testing.T) (*TokenManager, error) {
	os.Setenv("JWT_SECRET", "secret_key")
	os.Setenv("ACCESS_TTL", "15m")
	os.Setenv("REFRESH_TTL", "168h")
	db, err := storage.InitDatabase("postgres://your_db_user:your_db_password@postgres_auth:5432/your_db_name?sslmode=disable")
	if err != nil {
		t.Fatal("Failed to create Test Connect DB")
	}

	err = db.CreateStatements()
	if err != nil {
		t.Fatal("Failed to create Test Statements")
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		t.Fatal("Failed to create Test Repository")
	}

	TestService, err := NewTokenManager(repo)
	if err != nil {
		t.Fatal("Failed to create Test Service")
	}

	return TestService, nil
}

func CreateTestData(t *testing.T, testCount int) ([]models.User, error) {
	testData := make([]models.User, testCount)
	for idx := range testData {
		testData[idx].Email = fmt.Sprintf("example%d@tester.com", idx)
		testData[idx].IP = fmt.Sprintf("192.168.0.%d", idx)
	}
	return testData, nil
}

func TestTokenManager(t *testing.T) {
	flag.Parse()
	t.Logf("Test with %d test-data", testCount)
	TestService, err := SetupTestEnv(t)
	if err != nil {
		t.Fatal("Failed test enviroment setup")
	}
	defer TestService.repo.Storage.DB.Close()

	TestData, err := CreateTestData(t, testCount)
	if err != nil {
		t.Fatal("Failed to create test data")
	}
	TestData, err = TestService.InsertData(TestData)
	if err != nil {
		t.Fatal("Failed to test Insert Data")
	}

	for idx, user := range TestData {
		_, err = TestService.GetToken(user.GUID, user.IP)
		if err != nil {
			t.Errorf("Failed to test data %d with %s, %s", idx, user.GUID, user.IP)
		}
		_, err = TestService.CheckToken(user.RefreshTokenHash, user.IP)
		if err != nil {
			t.Errorf("Failed to test data %d with %s, %s", idx, user.RefreshTokenHash, user.IP)
		}

		user.IP = fmt.Sprintf("192.168.0.%d", idx)
		user.GUID = "0"

		_, err = TestService.GetToken(user.GUID, user.IP)
		if err == nil {
			t.Errorf("Failed to test data %d with %s, %s", idx, user.GUID, user.IP)
		}
		_, err = TestService.CheckToken(user.RefreshTokenHash, user.IP)
		if err == nil {
			t.Errorf("Failed to test data %d with %s, %s", idx, user.RefreshTokenHash, user.IP)
		}

	}

}
