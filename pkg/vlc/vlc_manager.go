package vlc

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"vlc_remote/pkg/logger"
)

type VlcPlayer struct {
	VlcConnectStatus  int      `json:"vlc_connect_status"`
	VlcRunning        bool     `json:"vlc_running"`
	LastSyncTime      int64    `json:"last_sync_time"`
	VlcError          string   `json:"vlc_error"`
	CurrentPlaying    string   `json:"current_playing"`
	Volume            string   `json:"volume"`
	Speed             float64  `json:"speed"`
	AudioTrack        string   `json:"audio_track"`
	AudioTrackList    []string `json:"audio_track_list"`
	SubtitleTrack     string   `json:"subtitle_track"`
	SubtitleTrackList []string `json:"subtitle_track_list"`

	vlcUrl      string
	vlcPassword string
	vlcUser     string
}

type vlcResponse struct {
	XMLName      xml.Name `xml:"root"`
	Text         string   `xml:",chardata"`
	Volume       string   `xml:"volume"`
	Version      string   `xml:"version"`
	Videoeffects struct {
		Text       string `xml:",chardata"`
		Brightness string `xml:"brightness"`
		Saturation string `xml:"saturation"`
		Gamma      string `xml:"gamma"`
		Contrast   string `xml:"contrast"`
		Hue        string `xml:"hue"`
	} `xml:"videoeffects"`
	Loop          string `xml:"loop"`
	Equalizer     string `xml:"equalizer"`
	Length        string `xml:"length"`
	Apiversion    string `xml:"apiversion"`
	Subtitledelay string `xml:"subtitledelay"`
	State         string `xml:"state"`
	Fullscreen    string `xml:"fullscreen"`
	Currentplid   string `xml:"currentplid"`
	Audiofilters  struct {
		Text    string `xml:",chardata"`
		Filter0 string `xml:"filter_0"`
	} `xml:"audiofilters"`
	Position    string `xml:"position"`
	Audiodelay  string `xml:"audiodelay"`
	Time        string `xml:"time"`
	Repeat      string `xml:"repeat"`
	Random      string `xml:"random"`
	Rate        string `xml:"rate"`
	Information struct {
		Text     string `xml:",chardata"`
		Category struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"category"`
	} `xml:"information"`
	Stats string `xml:"stats"`
}

func Init(vlcHost, vlcPort, vlcApiLogin, vlcApiPass string) VlcPlayer {
	return VlcPlayer{
		vlcUrl:      fmt.Sprintf("http://%s:%s/requests/", vlcHost, vlcPort),
		vlcPassword: vlcApiPass,
		vlcUser:     vlcApiLogin,
	}
}

func (vp *VlcPlayer) IsAlive() bool {
	res := true
	resp, err := vp.apiCall(http.MethodGet, "status.xml", map[string]interface{}{})
	if err != nil || resp.StatusCode != 200 {
		res = false
	}
	vp.VlcRunning = res
	return vp.VlcRunning
}

func (vp *VlcPlayer) Pause() bool {
	resp, err := vp.apiCall(http.MethodGet, "status.xml", map[string]interface{}{
		"command": "pl_pause",
	})
	return err == nil && resp != nil
}

func (vp *VlcPlayer) Play() bool {
	resp, err := vp.apiCall(http.MethodGet, "status.xml", map[string]interface{}{
		"command": "pl_play",
	})
	return err == nil && resp != nil
}

func (vp *VlcPlayer) apiCall(method string, path string, params map[string]interface{}) (*http.Response, error) {
	client := &http.Client{}

	var str []byte
	if method != http.MethodGet && len(params) > 0 {
		str, _ = json.Marshal(params)
	}

	req, _ := http.NewRequest(method, vp.vlcUrl+path, bytes.NewBuffer(str))

	if method == http.MethodGet && len(params) > 0 {
		q := req.URL.Query()
		for i, value := range params {
			q.Add(i, value.(string))
		}
		req.URL.RawQuery = q.Encode()
	}

	req.SetBasicAuth("", vp.vlcPassword)

	response, err := client.Do(req)
	if err != nil {
		logger.Error("error while connect to vlc", err)
	}
	if response != nil {
		vp.VlcConnectStatus = response.StatusCode
		if response.StatusCode == 401 {
			logger.Error("invalid vlc login/pass", nil)
		}
	}
	vp.parseResp(response.Body)
	return response, err
}

func (vp *VlcPlayer) parseResp(body io.ReadCloser) {
	data, err := ioutil.ReadAll(body)
	now := time.Now()
	vp.LastSyncTime = now.Unix()
	if err != nil {
		vp.VlcError = "Invalid response from vlc"
		return
	}
	var result vlcResponse
	eU := xml.Unmarshal(data, &result)
	if eU != nil {
		vp.VlcError = "Invalid structure response from vlc"
		return
	}

	vp.Volume = result.Volume
	vp.CurrentPlaying = result.State
}
