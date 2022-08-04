package main

import "testing"

func TestRepo_validate(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "letters+digits",
			dir:     "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			wantErr: false,
		},
		{
			name:    "special characters",
			dir:     "-_",
			wantErr: false,
		},
		{
			name:    "letters+digits+special characters",
			dir:     "a-0_",
			wantErr: false,
		},
		{
			name:    "invalid characters",
			dir:     "a$1.;",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repo{
				Name: tt.dir,
				Url:  "",
			}
			if err := r.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Repo.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
