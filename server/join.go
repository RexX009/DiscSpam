package server

import (
	"Raid-Client/cloudflare"
	"Raid-Client/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/patrickmn/go-cache"
)

var C *cache.Cache
var green = color.FgGreen.Render
var white = color.FgWhite.Render
var red = color.FgRed.Render

func JoinServer(inviteCode string, token string) error {
	code := ""
	if strings.Contains(inviteCode, "https://discord") {
		j := strings.Split(inviteCode, "/")
		code = j[3]
	} else {
		code = inviteCode
	}
	payload := map[string]string{"": ""}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	request, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%s", code), payloadBuf)
	if err != nil {
		return err
	}
	cf := cloudflare.Cookie()
	xprop := utils.XSuperProperties()
	request.Header = http.Header{
		"Accept":             []string{"*/*"},
		"Accept-language":    []string{"en-GB"},
		"Authorization":      []string{token},
		"Content-length":     []string{"2"},
		"Content-type":       []string{"application/json"},
		"Cookie":             []string{cf},
		"Origin":             []string{"https://discord.com"},
		"Sec-fetch-dest":     []string{"empty"},
		"Sec-fetch-mode":     []string{"cors"},
		"Sec-fetch-site":     []string{"same-origin"},
		"User-agent":         []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.15 Chrome/83.0.4103.122 Electron/9.3.5 Safari/537.36"},
		"X-debug-options":    []string{"bugReporterEnabled"},
		"X-super-properties": []string{xprop},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var responseJson map[string]interface{}
	if err := json.Unmarshal(body, &responseJson); err != nil {
		return err
	}
	ParseServerID := responseJson["guild"].(map[string]interface{})
	ServerID := ParseServerID["id"].(string)
	ServerName := ParseServerID["name"].(string)
	C.Set("JoinServerInvite", ServerID, cache.NoExpiration)
	fmt.Printf("%s %s %s", white(token), green("| Successfully Joined"), white(ServerName))

	return nil
}
