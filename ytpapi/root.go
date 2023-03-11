package ytpapi

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
)

func ChannelsListByUsername(service *youtube.Service, part []string, forUsername string) {
	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()

	if err != nil {
		log.Fatalf("Failed to get channel list: %v", err.Error())
	}

	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}
