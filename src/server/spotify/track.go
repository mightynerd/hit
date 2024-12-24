package spotify

import (
	"fmt"
)

func (spotify *Spotify) Play(trackId string) error {
	fmt.Println("playing track", trackId)

	resp, err := spotify.newRequest().SetBody(fmt.Sprintf(`{"uris":["%s"]}`, trackId)).Put("/me/player/play")

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not play")
	}

	if resp.StatusCode() != 204 {
		fmt.Println("could not play", resp.StatusCode(), string(resp.Body()))
		return fmt.Errorf("could not play")
	}

	return nil
}
