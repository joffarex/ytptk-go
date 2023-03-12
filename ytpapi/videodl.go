package ytpapi

import (
	"bytes"
	"encoding/json"
	"github.com/joffarex/ytptk-go/sliceutils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"google.golang.org/api/youtube/v3"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"os/user"
	"path/filepath"
	jsinterp "rogchap.com/v8go"
	"strings"
	"sync"
	"time"
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

func handleScopeAsync(wg *sync.WaitGroup, progress *mpb.Progress, jsinterpContext *jsinterp.Context, url string, path string, scope string) *ffmpeg.Stream {
	wg.Add(1)

	barName := "[" + scope + "]"
	bar := progress.AddBar(
		0,
		mpb.PrependDecorators(
			decor.Name(barName, decor.WC{W: len(barName) + 1, C: decor.DidentRight}),
			decor.CountersKiloByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.Name(" | ETA: "),
			decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncWidth), "done"),
			decor.Name(" | Speed: "),
			decor.EwmaSpeed(decor.UnitKB, "% .2f", 60),
		),
	)

	go func() {
		client := &http.Client{}

		parsedUrl, _ := neturl.Parse(url)
		params := parsedUrl.Query()

		// Taken from https://www.youtube.com/s/player/e06dea74/player_ias.vflset/en_US/base.js which by itself was taken
		// from https://github.com/yt-dlp/yt-dlp/blob/cf9fd52fabe71d6e7c30d3ea525029ffa561fc9c/test/test_youtube_signature.py#L94
		// Generally we are trying to replicate this logic here https://github.com/yt-dlp/yt-dlp/blob/cf9fd52fabe71d6e7c30d3ea525029ffa561fc9c/yt_dlp/extractor/youtube.py#L3024
		jsinterpContext.RunScript("var iha=function(_){var n=_.split(\"\"),e=[1085621920,-1413232877,245246939,-214478372,-1947691414,-1440840439,function(_,n){n=(n%_.length+_.length)%_.length;var e=_[0];_[0]=_[n],_[n]=e},-28067232,function(_,n){for(n=(n%_.length+_.length)%_.length;n--;)_.unshift(_.pop())},-1253364940,1108783393,1993082094,-2141271272,-1236697005,2066803531,null,-122169779,1673105371,952868508,1462692941,function(){for(var _=64,n=[];++_-n.length-32;){switch(_){case 58:_-=14;case 91:case 92:case 93:continue;case 123:_=47;case 94:case 95:case 96:continue;case 46:_=95}n.push(String.fromCharCode(_))}return n},1737827501,2066803531,99253448,function(_,n){n=(n%_.length+_.length)%_.length,_.splice(0,1,_.splice(n,1,_[0])[0])},661884814,n,1987089113,1657437625,1981674291,function(_){for(var n=_.length;n;)_.push(_.splice(--n,1)[0])},\"Y8eupT\",-188239353,-811835968,n,-795230947,null,1894196126,function(_,n){n=(n%_.length+_.length)%_.length,_.splice(-n).reverse().forEach(function(n){_.unshift(n)})},n,-259629667,function(_,n,e){var t=e.length;_.forEach(function(_,n,$){this.push($[n]=e[(e.indexOf(_)-e.indexOf(this[n])+n+t--)%e.length])},n.split(\"\"))},245246939,null,-1804836242,421896565,function(_,n){n=(n%_.length+_.length)%_.length,_.splice(n,1)},1108783393,function(_){_.reverse()},-1755683472,function(_,n){_.push(n)}];e[15]=e,e[36]=e,e[43]=e;try{e[24](e[36],e[23]),e[38](e[43],e[44]),e[12](e[30],e[9]),e[42](e[43],e[8]),e[3](e[40],e[2]),e[34](e[47]),e[50](e[32]),e[28](e[43],e[21]),e[28](e[38],e[37]),e[12](e[30],e[44]),e[50](e[38],e[39]),e[10](e[19],e[11]),e[45](e[30],e[35],e[24]()),e[42](e[47],e[5]),e[25](e[13],e[34]),e[47](e[26],e[1]),e[47](e[2],e[12]),e[51](e[36],e[35]),e[1](e[8],e[25]),e[24](e[42],e[19]),e[22](e[0],e[28]),e[24](e[3],e[27]),e[12](e[3],e[45]),e[10](e[50],e[26]),e[12](e[16],e[43]),e[12](e[50],e[0]),e[2](e[31],e[29]),e[9](e[13]),e[25](e[5],e[50]),e[27](e[31],e[49]),e[25](e[13],e[20]),e[37](e[5],e[40]),e[30](e[31],e[32]),e[28](e[46]),e[13](e[30],e[19]),e[37](e[1],e[7]),e[25](e[6],e[45]),e[13](e[39],e[5])}catch(t){return\"enhanced_except_z5QBsOv-_w8_\"+_}return n.join(\"\")};", "n-resolver.js")
		nResolved, _ := jsinterpContext.RunScript("iha("+params.Get("n")+")", "main.js")

		params.Del("n")
		params.Add("n", nResolved.String())

		parsedUrl.RawQuery = params.Encode()

		request, requestError := http.NewRequest(http.MethodGet, parsedUrl.String(), nil)
		if requestError != nil {
			log.Fatalf("Request %s error: %v", scope, requestError.Error())
		}
		// ----------------------
		download, downloadErr := client.Do(request)
		if downloadErr != nil {
			log.Fatalf("Download %s error: %v", scope, downloadErr.Error())
		}

		file, _ := os.Create(filepath.Join(path, scope+".mp4"))

		bar.SetTotal(download.ContentLength, false)
		var src io.Reader
		src = download.Body
		proxyReader := bar.ProxyReader(src)
		defer proxyReader.Close()
		_, audioErr := io.Copy(file, proxyReader)

		if audioErr != nil {
			log.Fatalf("%s error: %v", scope, audioErr.Error())
		}

		defer wg.Done()
	}()

	inputAudio := ffmpeg.Input(filepath.Join(path, scope+".mp4"))

	return inputAudio
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

	osUser, _ := user.Current()
	path := filepath.Join(osUser.HomeDir, "ytptk", "_temp", destination)
	os.MkdirAll(path, 0700)

	// --------------- URLs ---------------
	audio := sliceutils.Filter(output.StreamingData.AdaptiveFormats, func(item AdaptiveFormat) bool {
		return strings.Contains(item.MimeType, "audio/mp4") && item.AudioQuality == "AUDIO_QUALITY_MEDIUM"
	})[0]
	video := sliceutils.Filter(output.StreamingData.AdaptiveFormats, func(item AdaptiveFormat) bool {
		return strings.Contains(item.MimeType, "video/mp4") && item.Quality == "hd1080"
	})[0]

	var wg sync.WaitGroup

	progress := mpb.New(mpb.WithWidth(64), mpb.WithWaitGroup(&wg), mpb.WithRefreshRate(180*time.Millisecond))

	ctx := jsinterp.NewContext()

	ffmpegInputAudio := handleScopeAsync(&wg, progress, ctx, audio.Url, path, "audio")
	ffmpegInputVideo := handleScopeAsync(&wg, progress, ctx, video.Url, path, "video")

	wg.Wait()

	encodingErr := ffmpeg.Concat([]*ffmpeg.Stream{ffmpegInputAudio, ffmpegInputVideo}, ffmpeg.KwArgs{"v": 1, "a": 1}).Output(filepath.Join(osUser.HomeDir, "ytptk", destination+".mp4")).Run()

	if encodingErr != nil {
		log.Fatalf("Encoding error: %v", encodingErr.Error())
	}
}
