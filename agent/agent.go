package agent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

  "github.com/gorilla/websocket"
	"golang.org/x/net/html"
)

func Execute() {
	requestURL := "http://localhost:4000"
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("could not create request: %s\n", err)
		os.Exit(1)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("error making HTTP request: %s\n", err)
		os.Exit(1)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading HTTP response: %s\n", err)
		os.Exit(1)
	}

	csrfToken := parseCsrfToken(fmt.Sprintf("%s", responseBody))
	fmt.Println(csrfToken)
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
			if token.Data == "meta" && getAttribute("name", token) == "csrf-token" {
				return getAttribute("content", token)
			}
		}
	}
}

func getAttribute(name string, token html.Token) string {
	for _, attribute := range token.Attr {
		if attribute.Key == name {
			return attribute.Val
		}
	}

	return ""
}
