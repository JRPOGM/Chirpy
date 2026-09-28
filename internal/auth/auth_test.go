package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	password3 := "originalPassword789!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	hash3, _ := HashPassword(password3)
	test := []struct {
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
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && macth != tt.macthPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}