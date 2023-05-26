package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"outdoorsy/db"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestCase struct {
	Method       string
	URL          string
	Body         io.Reader
	ExpectedCode int
}

// Runs before every test
func TestMain(m *testing.M) {
	// Get Config
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestRentalsEndpointSuccess(t *testing.T) {
	testCases := []TestCase{
		// Success cases
		{
			Method:       "GET",
			URL:          "/rentals/1",
			Body:         nil,
			ExpectedCode: http.StatusOK,
		},
		// Error cases
		{
			Method:       "GET",
			URL:          "/rentals/a",
			Body:         nil,
			ExpectedCode: http.StatusBadRequest,
		},
	}

	database, err := db.NewDatabase(nil)
	assert.NoError(t, err)

	router, err := NewRouter(database)
	assert.NoError(t, err)

	for i, testCase := range testCases {
		t.Logf("Running test case %d: %s", i, testCase.URL)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(testCase.Method, testCase.URL, testCase.Body)
		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, testCase.ExpectedCode, w.Code)

		if testCase.ExpectedCode == http.StatusOK {
			var response Rental
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "json should marshal to Rental struct")
		}
	}

	// TODO: check that response has all the necessary data

	// TODO: add error tests
}

type EndpointRentalsTestSuite struct {
	suite.Suite
	Router *gin.Engine
}

func (suite *EndpointRentalsTestSuite) SetupTest() {
	database, err := db.NewDatabase(nil)
	assert.NoError(suite.T(), err)

	router, err := NewRouter(database)
	assert.NoError(suite.T(), err)

	suite.Router = router
}

func (suite *EndpointRentalsTestSuite) TestGetRentals() {
	testCases := []TestCase{
		// Success cases
		{
			Method:       "GET",
			URL:          "/rentals?price_min=9000",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?price_min=9000&price_max=75000",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?limit=3&offset=6",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?ids=3,4,5",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?near=33.64,-117.93",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?sort=price",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?near=33.64,-117.93",
			Body:         nil,
			ExpectedCode: 200,
		},
		{
			Method:       "GET",
			URL:          "/rentals?near=33.64,-117.93&price_min=9000&price_max=75000&limit=3&offset=6&sort=price",
			Body:         nil,
			ExpectedCode: 200,
		},
		// Error Cases
		// TODO: use validation to surface a 400 instead of a 500
		{
			Method:       "GET",
			URL:          "/rentals?near=-117.93",
			Body:         nil,
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	for i, testCase := range testCases {
		suite.T().Logf("Running test case %d: %s", i, testCase.URL)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(testCase.Method, testCase.URL, testCase.Body)
		suite.Router.ServeHTTP(w, req)

		// Assertions
		suite.Equal(testCase.ExpectedCode, w.Code)

		if testCase.ExpectedCode == http.StatusOK {
			var response []Rental
			err := json.Unmarshal(w.Body.Bytes(), &response)
			suite.NoError(err, "json should marshal to Rental struct")

			suite.T().Log(w.Body.String())
		}

		// TODO: validate response body
	}
}

func TestEndpointRentalsTestSuite(t *testing.T) {
	suite.Run(t, new(EndpointRentalsTestSuite))
}
