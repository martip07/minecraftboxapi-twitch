package procs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/martip07/minecraftboxapi/structs"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
)

func TwitchAuth() string {
	type tokenData struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `token_type:"id"`
	}

	cli := gentleman.New()
	cli.Method("POST")
	clientId := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")
	uriBase := "https://id.twitch.tv/oauth2/token?client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=client_credentials"
	res, err := cli.Request().URL(uriBase).Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
	}
	json := &tokenData{}
	res.JSON(json)
	fmt.Printf("Body: %#v\n", json.AccessToken)
	return json.AccessToken
}

func TwitchProc(_idGame string, _idLanguage string) structs.Streams {
	cli := gentleman.New()
	cli.SetHeader("App", "minecraftbox")
	auth := TwitchAuth()
	authHeader := "Bearer " + auth
	//fmt.Println(authHeader)
	cli.Use(headers.Set("Client-ID", os.Getenv("CLIENTID")))
	cli.Use(headers.Set("Authorization", authHeader))
	uriBase := "https://api.twitch.tv/helix/streams?first=20&game_id=" + _idGame + "&language=" + _idLanguage
	res, err := cli.Request().URL(uriBase).Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
	}

	fmt.Printf("Status: %d\n", res.StatusCode)
	var TwitchRes structs.Streams
	var TwitchStream structs.Stream
	var TwitchArr []structs.Stream
	resByte := res.Bytes()
	if resByte != nil {
		json.Unmarshal(resByte, &TwitchRes)
	}
	for _, StreamElement := range TwitchRes.StreamData {

		str := StreamElement.ThumbnailURL
		str = strings.Replace(str, "{width}", "350", -1)
		str = strings.Replace(str, "{height}", "220", -1)

		TwitchStream.GameID = StreamElement.GameID
		TwitchStream.StartedAt = StreamElement.StartedAt
		TwitchStream.StreamID = StreamElement.StreamID
		TwitchStream.StreamLanguage = StreamElement.StreamLanguage
		TwitchStream.StreamTitle = StreamElement.StreamTitle
		TwitchStream.StreamType = StreamElement.StreamType
		TwitchStream.TagIDS = StreamElement.TagIDS
		TwitchStream.ThumbnailURL = str
		TwitchStream.UserID = StreamElement.UserID
		TwitchStream.UserName = StreamElement.UserName
		TwitchStream.ViewerCount = StreamElement.ViewerCount

		TwitchArr = append(TwitchArr, TwitchStream)

	}

	TwitchRes.StreamData = TwitchArr
	return TwitchRes
}
