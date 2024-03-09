package utils

import "testing"

func TestGetIDFromPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr error
	}{
		{
			name:    "valid path",
			path:    "/movies/123",
			want:    "123",
			wantErr: nil,
		},
		{
			name:    "valid path with trailing slash",
			path:    "/movies/123/",
			want:    "123",
			wantErr: nil,
		},
		{
			name:    "invalid path with extra segment",
			path:    "/movies/123/456",
			want:    "",
			wantErr: ErrInvalidPath,
		},
		{
			name:    "invalid path with missing segment",
			path:    "/movies",
			want:    "",
			wantErr: ErrInvalidPath,
		},
		{
			name:    "invalid path with missing segment and trailing slash",
			path:    "/movies/",
			want:    "",
			wantErr: ErrInvalidPath,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetIDFromPath(test.path)

			if got != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}

			if err != test.wantErr {
				t.Errorf("got %v, want %v", err, test.wantErr)
			}
		})
	}
}
