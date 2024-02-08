package youtube

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type YouTubeURL struct {
	VideoID          string
	TimestampSeconds int
}

func IsYoutubeURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return parsedURL.Host == "www.youtube.com" || parsedURL.Host == "youtu.be"
}

func DeconstructYouTubeURL(rawURL string) (*YouTubeURL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	regexPattern := `(?i)(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))([\w-]{11})`
	regex := regexp.MustCompile(regexPattern)
	videoIDMatches := regex.FindStringSubmatch(rawURL)

	if len(videoIDMatches) < 2 {
		return nil, fmt.Errorf("unable to extract video ID from URL")
	}

	videoID := videoIDMatches[1]

	timestamp := time.Duration(0)
	timestampRegexPattern := `t=([0-9]+)(s|m|h)?`
	timestampRegex := regexp.MustCompile(timestampRegexPattern)
	timestampMatches := timestampRegex.FindStringSubmatch(parsedURL.RawQuery)

	if len(timestampMatches) >= 2 {
		durationStr := strings.ToLower(timestampMatches[1])
		unitStr := strings.ToLower(timestampMatches[2])

		// Remove the unit from the duration string
		durationStr = regexp.MustCompile(unitStr).ReplaceAllString(durationStr, "")
		durationInt := 0

		// Add the unit back to the duration string
		if unitStr == "s" || unitStr == "" {
		} else if unitStr == "m" {
			fmt.Sscanf(durationStr, "%d", &durationInt)
			durationInt *= 60
			durationStr = fmt.Sprintf("%d", durationInt)
		} else if unitStr == "h" {
			fmt.Sscanf(durationStr, "%d", &durationInt)
			durationInt *= 60 * 60
			durationStr = fmt.Sprintf("%d", durationInt)
		}

		durationStr += "s"

		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing duration: %v", err)
		}

		timestamp = duration
	}

	return &YouTubeURL{
		VideoID:          videoID,
		TimestampSeconds: int(timestamp.Truncate(time.Second).Seconds()),
	}, nil
}
