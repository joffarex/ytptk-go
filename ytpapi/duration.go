package ytpapi

import (
	"github.com/joffarex/ytptk-go/sliceutils"
	"google.golang.org/api/youtube/v3"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var playlistItemParts = []string{"contentDetails"}
var videoParts = []string{"contentDetails", "id", "liveStreamingDetails", "player", "status"}

func toSeconds(input string) int {
	regexPattern := "^P(?:(\\d+)D)?T(?:(\\d+)H)?(?:(\\d+)M)?(?:(\\d+)S)?$"

	match, err := regexp.MatchString(regexPattern, input)

	if (err != nil) || !match {
		log.Fatalf("No match: %v, Err: %v", match, err)
	}

	pattern := regexp.MustCompile(regexPattern)
	matches := pattern.FindAllStringSubmatch(input, -1)

	rawDays := matches[0][1]
	rawHours := matches[0][2]
	rawMinutes := matches[0][3]
	rawSeconds := matches[0][4]

	days, _ := strconv.Atoi(rawDays)
	hours, _ := strconv.Atoi(rawHours)
	minutes, _ := strconv.Atoi(rawMinutes)
	seconds, _ := strconv.Atoi(rawSeconds)

	totalSeconds := days*86400 + hours*3600 + minutes*60 + seconds

	return totalSeconds
}

func GetTotalPlaylistDuration(service *youtube.Service, playlistId string, pageToken string) int {
	playlistItemsResponse, playlistItemsError := service.PlaylistItems.List(playlistItemParts).PlaylistId(playlistId).MaxResults(50).PageToken(pageToken).Do()

	if playlistItemsError != nil {
		log.Fatalf("Failed to get playlist items list: %v", playlistItemsError.Error())
	}

	videoIds := sliceutils.Map(playlistItemsResponse.Items, func(item *youtube.PlaylistItem) string { return item.ContentDetails.VideoId })

	videosResponse, videosError := service.Videos.List(videoParts).Id(strings.Join(videoIds, ",")).Do()

	if videosError != nil {
		log.Fatalf("Failed to get playlist items list: %v", videosError.Error())
	}

	totalDuration := 0
	for _, video := range videosResponse.Items {
		totalDuration = totalDuration + toSeconds(video.ContentDetails.Duration)
	}

	if playlistItemsResponse.NextPageToken != "" {
		totalDuration = totalDuration + GetTotalPlaylistDuration(service, playlistId, playlistItemsResponse.NextPageToken)
	}

	return totalDuration
}
