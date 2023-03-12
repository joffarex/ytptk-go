package ytpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joffarex/ytptk-go/sliceutils"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"google.golang.org/api/youtube/v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type InputClientMainApp struct {
	GraftUrl string `json:"graftUrl"`
}

type InputClient struct {
	Hl             string             `json:"hl"`
	Gl             string             `json:"gl"`
	ClientName     string             `json:"clientName"`
	ClientVersion  string             `json:"clientVersion"`
	MainAppWebInfo InputClientMainApp `json:"mainAppWebInfo"`
}

type InputContext struct {
	Client InputClient `json:"client"`
}

type Input struct {
	VideoId string       `json:"videoId"`
	Context InputContext `json:"context"`
}

type Format struct {
	Url      string `json:"url"`
	MimeType string `json:"mimeType"`
}

type AdaptiveFormat struct {
	Format
	Quality       string `json:"quality"`
	ContentLength string `json:"contentLength"`
	Height        int    `json:"height"`
	AudioQuality  string `json:"audioQuality"`
}

type StreamingData struct {
	Formats         []Format         `json:"formats"`
	AdaptiveFormats []AdaptiveFormat `json:"adaptiveFormats"`
}

type Output struct {
	StreamingData StreamingData `json:"StreamingData"`
}

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	Total   int64 // Total # bytes to be transferred
	Current int64 // Current # amount of bytes transferred
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += int64(n)
	fmt.Printf("Chunk size: %d | Current Sum: %d | Total: %d | Left: %d\n", n, wc.Current, wc.Total, wc.Total-wc.Current)
	//fmt.Printf("Read %d bytes for a total of %d\n", n, wc.Current)
	return n, nil
}

func DownloadVideoToDestination(service *youtube.Service, videoId string, destination string) {
	input := Input{VideoId: videoId,
		Context: InputContext{
			Client: InputClient{
				Hl:            "en",
				Gl:            "US",
				ClientName:    "WEB",
				ClientVersion: "2.20210721.00.00",
				MainAppWebInfo: InputClientMainApp{
					GraftUrl: "/watch?v=" + videoId,
				},
			},
		},
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(input)

	res, _ := http.Post("https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8", "application/json", &buf)

	bytes, _ := io.ReadAll(res.Body)
	var output Output
	_ = json.Unmarshal(bytes, &output)

	audio := sliceutils.Filter(output.StreamingData.AdaptiveFormats, func(item AdaptiveFormat) bool {
		return strings.Contains(item.MimeType, "audio/mp4") && item.AudioQuality == "AUDIO_QUALITY_MEDIUM"
	})[0]
	downloadAudio, downloadAudioErr := http.Get(audio.Url)
	if downloadAudioErr != nil {
		log.Fatalf("Download Audio error: %v", downloadAudioErr.Error())
	}

	fullHD := sliceutils.Filter(output.StreamingData.AdaptiveFormats, func(item AdaptiveFormat) bool {
		return strings.Contains(item.MimeType, "video/mp4") && item.Height == 1080 && item.Quality == "hd1080"
	})[0]
	downloadFullHd, downloadFullHdErr := http.Get(fullHD.Url)
	if downloadFullHdErr != nil {
		log.Fatalf("Download Audio error: %v", downloadFullHdErr.Error())
	}

	osUser, _ := user.Current()
	path := filepath.Join(osUser.HomeDir, "ytptk", "_temp", destination)
	os.MkdirAll(path, 0700)
	audioFile, _ := os.Create(filepath.Join(path, "audio.mp4"))
	fullHDFile, _ := os.Create(filepath.Join(path, "fullHD.mp4"))

	var audioSrc io.Reader
	audioSrc = downloadAudio.Body
	_, audioErr := io.Copy(audioFile, audioSrc)

	if audioErr != nil {
		log.Fatalf("Audio error: %v", audioErr.Error())
	}

	var fullHDSrc io.Reader
	fullHDSrc = downloadFullHd.Body
	_, fullHDErr := io.Copy(fullHDFile, fullHDSrc)

	if fullHDErr != nil {
		log.Fatalf("Full HD error: %v", fullHDErr.Error())
	}

	inputAudio := ffmpeg_go.Input(filepath.Join(path, "audio.mp4"))
	inputFullHD := ffmpeg_go.Input(filepath.Join(path, "fullHD.mp4"))
	encodingErr := ffmpeg_go.Concat([]*ffmpeg_go.Stream{inputAudio, inputFullHD}, ffmpeg_go.KwArgs{"v": 1, "a": 1}).Output(filepath.Join(osUser.HomeDir, "ytptk", destination+".mp4")).Run()

	if encodingErr != nil {
		log.Fatalf("Encoding error: %v", encodingErr.Error())
	}

	// -------------------------------------------------
	//download, _ := http.Get(output.StreamingData.Formats[0].Url)
	//
	//osUser, _ := user.Current()
	//path := filepath.Join(osUser.HomeDir, "ytptk")
	//os.MkdirAll(path, 0700)
	//outFile, _ := os.Create(filepath.Join(path, destination+".mp4"))
	//
	//var src io.Reader
	//src = download.Body
	//
	//videosResponse, error := service.Videos.List([]string{"fileDetails"}).Id(videoId).Do()
	//fmt.Printf("\nVideo: %v\n", error)
	//
	//src = io.TeeReader(src, &WriteCounter{Total: int64(videosResponse.Items[0].FileDetails.FileSize)})
	//
	//_, _ = io.Copy(outFile, src)
}
