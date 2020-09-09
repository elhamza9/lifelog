package auth_test

import (
	"os"
	"strings"
	"testing"

	"github.com/elhamza90/lifelog/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
)

// hashEnvVarName specifies the name of the environment
// variable that holds the hash of the password used for test
const hashEnvVarName string = "LFLG_TEST_PASS_HASH"

// authenticator is the instance of the service to be tested
var authenticator auth.Service

func TestMain(m *testing.M) {
	authenticator = auth.NewService(hashEnvVarName)
	os.Exit(m.Run())
}

func TestAuthenticate(t *testing.T) {
	const testPass string = "correct_pass"
	hash, err := bcrypt.GenerateFromPassword([]byte(testPass), 10)
	if err != nil {
		t.Fatalf("Error generating bcrypt hash: %s", err)
	}
	os.Setenv(hashEnvVarName, string(hash))
	tests := map[string]struct {
		pass        string
		expectedErr error
	}{
		"Short Password":           {"pass", auth.ErrPasswordLength},
		"Long Password":            {strings.Repeat("pass", 100), auth.ErrPasswordLength},
		"Valid Incorrect Password": {"valid_wrong_pass", auth.ErrIncorrectCredentials},
		"Correct Password":         {testPass, nil},
	}
	// Subtests execution
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if err := authenticator.Authenticate(test.pass); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
		})
	}
}
