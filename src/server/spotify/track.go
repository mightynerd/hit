package spotify

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (config *Spotify) Play(trackId string) error {
	fmt.Println("playing track", trackId)
	client := http.Client{}
	request, err := config.newAPIRequest("PUT",
		"/me/player/play",
		strings.NewReader(fmt.Sprintf(`{"uris":["%s"]}`, trackId)))

	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		fmt.Println("could not play", resp.StatusCode, string(body))
		return fmt.Errorf("could not play")
	}

	return nil
}
