package youtube

import (
	"testing"
)

func TestDeconstructYouTubeURL(t *testing.T) {
	tests := []struct {
		name      string
		rawURL    string
		wantVideo string
		wantTime  int
		wantErr   bool
	}{
		{
			name:      "CaseInsensitive",
			rawURL:    "https://www.YouTube.cOm/watch?v=RBYJTop1e-g&t=104s",
			wantVideo: "RBYJTop1e-g",
			wantTime:  104,
			wantErr:   false,
		},
		{
			name:      "WithSecondTimestamp",
			rawURL:    "https://www.youtube.com/watch?v=RBYJTop1e-g&t=100s",
			wantVideo: "RBYJTop1e-g",
			wantTime:  100,
			wantErr:   false,
		},
		{
			name:      "WithMinuteTimestamp",
			rawURL:    "https://www.youtube.com/watch?v=RBYJTop1e-g&t=101m",
			wantVideo: "RBYJTop1e-g",
			wantTime:  6060,
			wantErr:   false,
		},
		{
			name:      "WithHourTimestamp",
			rawURL:    "https://www.youtube.com/watch?v=RBYJTop1e-g&t=1h",
			wantVideo: "RBYJTop1e-g",
			wantTime:  3600,
			wantErr:   false,
		},
		{
			name:      "WithoutTimestamp",
			rawURL:    "https://www.youtube.com/watch?v=RBYJTop1e-g",
			wantVideo: "RBYJTop1e-g",
			wantTime:  0,
			wantErr:   false,
		},
		{
			name:      "WithTimestamp",
			rawURL:    "https://youtu.be/RBYJTop1e-g?t=21",
			wantVideo: "RBYJTop1e-g",
			wantTime:  21,
			wantErr:   false,
		},
		{
			name:      "WithoutTimestamp",
			rawURL:    "https://youtu.be/RBYJTop1e-g",
			wantVideo: "RBYJTop1e-g",
			wantTime:  0,
			wantErr:   false,
		},
		{
			name:      "InvalidURL",
			rawURL:    "not_a_valid_url",
			wantVideo: "",
			wantTime:  0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeconstructYouTubeURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeconstructYouTubeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.VideoID != tt.wantVideo {
					t.Errorf("DeconstructYouTubeURL() gotVideo = %v, want %v", got.VideoID, tt.wantVideo)
				}
				if got.TimestampSeconds != tt.wantTime {
					t.Errorf("DeconstructYouTubeURL() gotTime = %v, want %v", got.TimestampSeconds, tt.wantTime)
				}
			}
		})
	}
}

func TestDeconstructYouTubeURL_InvalidURL(t *testing.T) {
	_, err := DeconstructYouTubeURL("not_a_valid_url")
	if err == nil {
		t.Error("Expected an error for an invalid URL, but got nil")
	}
}
