package engine

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tidwall/gjson"
)

func TestGetData(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	// Set the global db variable to the mock database
	// Assuming you have a global variable `db` in your package
	// db.db = db // Uncomment this line if you have a global db variable

	// Define the test cases
	tests := []struct {
		name         string
		queryJSON    string
		expectedData []string
		expectedErr  error
		mockSetup    func()
	}{
		{
			name:         "Valid Query",
			queryJSON:    `{"collection": "test", "limit": 2, "match": {}}`,
			expectedData: []string{"record1", "record2"},
			expectedErr:  nil,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT record FROM test`).
					WillReturnRows(sqlmock.NewRows([]string{"record"}).
						AddRow("record1").
						AddRow("record2"))
			},
		},
		{
			name:         "Missing Collection",
			queryJSON:    `{"limit": 2}`,
			expectedData: nil,
			expectedErr:  fmt.Errorf(`{"error":"forgot collection name "}`),
			mockSetup:    func() {},
		},
		{
			name:         "Database Error",
			queryJSON:    `{"collection": "test"}`,
			expectedData: nil,
			expectedErr:  fmt.Errorf("some database error"),
			mockSetup: func() {
				mock.ExpectQuery(`SELECT record FROM test`).WillReturnError(fmt.Errorf("some database error"))
			},
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock expectations
			tt.mockSetup()

			// Parse the query JSON
			query := gjson.Parse(tt.queryJSON)

			// Call the function
			data, err := getData(query)

			// Check the results
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if !equalSlices(data, tt.expectedData) {
				t.Errorf("expected data: %v, got: %v", tt.expectedData, data)
			}
		})
	}
}

// Helper function to compare slices
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
