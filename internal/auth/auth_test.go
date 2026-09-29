package auth

import (
	"time"
	"net/http"
	"testing"
	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	password3 := "originalPassword789!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	hash3, _ := HashPassword(password3)
	tests := []struct {
		name			string
		password		string
		hash			string
		wantError		bool
		matchPassword	bool
	}{
		{
			name:			"Correct password",
			password:		password1,
			hash:			hash1,
			wantError:		false,
			matchPassword:	true,
		},
		{
			name:			"Correct password",
			password:		password2,
			hash:			hash2,
			wantError:		false,
			matchPassword:	true,
		},
		{
			name:			"Correct password",
			password:		password3,
			hash:			hash3,
			wantError:		false,
			matchPassword:	true,
		},
		{
			name:			"Incorrect password",
			password:		password1,
			hash:			hash2,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Incorrect password",
			password:		password1,
			hash:			hash3,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Incorrect password",
			password:		password2,
			hash:			hash1,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Incorrect password",
			password:		password2,
			hash:			hash3,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Incorrect password",
			password:		password3,
			hash:			hash2,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Incorrect password",
			password:		password3,
			hash:			hash1,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Empty password",
			password:		"",
			hash:			hash1,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Empty password",
			password:		"",
			hash:			hash2,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Empty password",
			password:		"",
			hash:			hash3,
			wantError:		false,
			matchPassword:	false,
		},
		{
			name:			"Invalid hash",
			password:		password1,
			hash:			"This isn't a real hash",
			wantError:		true,
			matchPassword:	false,
		},
		{
			name:			"Invalid hash",
			password:		password3,
			hash:			"This isn't a real hash",
			wantError:		true,
			matchPassword:	false,
		},
		{
			name:			"Invalid hash",
			password:		password2,
			hash:			"This isn't a real hash",
			wantError:		true,
			matchPassword:	false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantError {
				t.Errorf("CheckPasswordHash() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)
	validToken2, _ := MakeJWT(userID, "truth", time.Hour)
	validToken3, _ := MakeJWT(userID, "power", time.Hour)
	validToken4, _ := MakeJWT(userID, "consequences", time.Hour)
	tests := []struct {
		name			string
		tokenString		string
		tokenSecret		string
		wantUserID		uuid.UUID
		wantError		bool
	}{
		{
			name:			"Valid token",
			tokenString:	validToken,
			tokenSecret:	"secret",
			wantUserID:		userID,
			wantError:		false,
		},
		{
			name:			"Invalid token",
			tokenString:	"invalid.token.string",
			tokenSecret:	"secret",
			wantUserID:		uuid.Nil,
			wantError:		true,
		},
		{
			name:			"Wrong secret",
			tokenString:	validToken,
			tokenSecret:	"wrong_secret",
			wantUserID:		uuid.Nil,
			wantError:		true,
		},
		{
			name:			"True token",
			tokenString:	validToken2,
			tokenSecret:	"truth",
			wantUserID:		userID,
			wantError:		false,
		},
		{
			name:			"Unreal token",
			tokenString:	"unreal.token.string",
			tokenSecret:	"truth",
			wantUserID:		uuid.Nil,
			wantError:		true,
		},
		{
			name:			"False truth",
			tokenString:	validToken2,
			tokenSecret:	"false_truth",
			wantUserID:		uuid.Nil,
			wantError:		true,
		},
		{
			name:			"Valid token",
			tokenString:	validToken3,
			tokenSecret:	"power",
			wantUserID:		userID,
			wantError:		false,
		},
		{
			name:			"Valid token",
			tokenString:	validToken4,
			tokenSecret:	"consequences",
			wantUserID:		userID,
			wantError:		false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateJWT() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name		string
		headers		http.Header
		wantToken	string
		wantError	bool
	}{
		{
			name:		"Valid Bearer token",
			headers:	http.Header{"Authorization": []string{"Bearer valid_token"}},
			wantToken:	"valid_token",
			wantError:	false,
		},
		{
			name:		"Missing Authorization header",
			headers:	http.Header{},
			wantToken:	"",
			wantError:	true,
		},
		{
			name:		"Malformed Authroization header",
			headers:	http.Header{"Authorization": []string{"InvalidBearer token"}},
			wantToken:	"",
			wantError:	true,
		},
		{
			name:		"Valid Bearer token",
			headers:	http.Header{"Authorization": []string{"Bearer true_token"}},
			wantToken:	"true_token",
			wantError:	false,
		},
		{
			name:		"Valid Bearer token",
			headers:	http.Header{"Authorization": []string{"Bearer correct_token"}},
			wantToken:	"correct_token",
			wantError:	false,
		},
		{
			name:		"Misspelled Authroization header",
			headers:	http.Header{"Amorization": []string{"Bearer valid_token"}},
			wantToken:	"",
			wantError:	true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantError {
				t.Errorf("GetBearerToken() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}