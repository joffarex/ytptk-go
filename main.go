package main

import (
	"fmt"
	"github.com/joffarex/ytptk-go/dateutils"
	"github.com/joffarex/ytptk-go/googleauth"
	"github.com/joffarex/ytptk-go/ytpapi"
)

func main() {
	service := googleauth.AuthenticateWithGoogle()

	playlistId := "PL7atuZxmT954bCkC062rKwXTvJtcqFB8i"
	durationInSeconds := ytpapi.GetTotalPlaylistDuration(service, playlistId, "")
	fmt.Printf("PlaylistId: %s, Duration: %s", playlistId, dateutils.SecondsToDuration(durationInSeconds))
}
