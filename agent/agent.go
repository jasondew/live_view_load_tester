package agent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

  "golang.org/x/net/websocket"
	"golang.org/x/net/html"
)

func Execute() {
	pageUrl := "http://fanatics.localhost:4000/breaking/44a9547d-a1a6-4a6e-94a5-3c5c4de66993"
	request, err := http.NewRequest(http.MethodGet, pageUrl, nil)
	if err != nil {
		fmt.Printf("could not create request: %s\n", err)
		os.Exit(1)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("error making HTTP request: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading HTTP response: %s\n", err)
		os.Exit(1)
	}

	responseBody := fmt.Sprintf("%s", responseBodyBytes)
	csrfToken := parseCsrfToken(responseBody)
	viewId, viewSession, viewStatic := parseLiveViewData(responseBody)

	websocketUrl := fmt.Sprintf("ws://fanatics.localhost:4000/live/websocket?_csrf_token=%s&timezone=America/New_York&vsn=2.0.0", csrfToken)
	wsConfig, err := websocket.NewConfig(websocketUrl, pageUrl)
	wsConfig.Header.Set("Cookie", response.Header.Get("Set-Cookie"))
	ws, err := websocket.DialConfig(wsConfig)

	if err != nil {
		fmt.Printf("error making WS request: %s\n", err)
		os.Exit(1)
	}

	joinCommand := fmt.Sprintf(`["4","4","lv:%s","phx_join",{"url":"%s","params":{"_csrf_token":"%s","timezone":"America/New_York","_mounts":0},"session":"%s","static":"%s"}]`, viewId, pageUrl, csrfToken, viewSession, viewStatic)
	if _, err := ws.Write([]byte(joinCommand)); err != nil {
		fmt.Printf("error sending WS command: %s\n", err)
		os.Exit(1)
	}

	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		fmt.Printf("error received WS response: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}

func parseCsrfToken(text string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(text))

	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken:
			return ""
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" && getAttribute(token, "name") == "csrf-token" {
				return getAttribute(token, "content")
			}
		}
	}
}

func parseLiveViewData(text string) (string, string, string) {
	tokenizer := html.NewTokenizer(strings.NewReader(text))

	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken:
			return "", "", ""
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "div" && getAttribute(token, "data-phx-main") == "true" {
				return getAttribute(token, "id"),
					getAttribute(token, "data-phx-session"),
					getAttribute(token, "data-phx-static")
			}
		}
	}
}

func getAttribute( token html.Token, name string) string {
	for _, attribute := range token.Attr {
		if attribute.Key == name {
			return attribute.Val
		}
	}

	return ""
}
