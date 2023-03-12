package main

import (
	"fmt"
	"github.com/joffarex/ytptk-go/dateutils"
	"github.com/joffarex/ytptk-go/googleauth"
	"github.com/joffarex/ytptk-go/ytpapi"
)

func main() {
	service := googleauth.AuthenticateWithGoogle()

	playlistId := "PL9ukzXzOMKcC6t4ecuc11Ud0--n5tirr4"
	durationInSeconds := ytpapi.GetTotalPlaylistDuration(service, playlistId, "")
	fmt.Printf("PlaylistId: %s, Duration: %s", playlistId, dateutils.SecondsToDuration(durationInSeconds))

	ytpapi.DownloadVideoToDestination(service, "bTE73vo5jxc", "testfilename")
}
