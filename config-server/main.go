package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type spotifyConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RewardID     string `json:"reward_id"`
	Mode         string `json:"mode"`
}

type twitchAPIConfig struct {
	TwitchAPI string `json:"twitch_api"`
}

type config struct {
	DiscordToken             string          `json:"discord_token"`
	AppID                    string          `json:"app_id"`
	AppClientID              string          `json:"app_client_id"`
	AppClientSecret          string          `json:"app_client_secret"`
	TwitchToken              string          `json:"twitch_token"`
	TwitchChannelID          string          `json:"twitch_channel_id"`
	WebhookURL               string          `json:"webhook_url"`
	WebhookSecret            string          `json:"webhook_secret"`
	DiscordUserID            string          `json:"discord_user_id"`
	GameDetectorIdleCategory string          `json:"game_detector_idle_category"`
	ApiURLs                  twitchAPIConfig `json:"api_urls"`
	Spotify                  spotifyConfig   `json:"spotify"`
}

func main() {
	apiKey := os.Getenv("CONFIG_API_KEY")
	if apiKey == "" {
		log.Fatal("CONFIG_API_KEY env var is required")
	}

	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(configFromEnv())
	})

	log.Println("config-server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func configFromEnv() config {
	return config{
		DiscordToken:             os.Getenv("DISCORD_TOKEN"),
		AppID:                    os.Getenv("APP_ID"),
		AppClientID:              os.Getenv("APP_CLIENT_ID"),
		AppClientSecret:          os.Getenv("APP_CLIENT_SECRET"),
		TwitchToken:              os.Getenv("TWITCH_TOKEN"),
		TwitchChannelID:          os.Getenv("TWITCH_CHANNEL_ID"),
		WebhookURL:               os.Getenv("WEBHOOK_URL"),
		WebhookSecret:            os.Getenv("WEBHOOK_SECRET"),
		DiscordUserID:            os.Getenv("DISCORD_USER_ID"),
		GameDetectorIdleCategory: envOrDefault("GAME_DETECTOR_IDLE_CATEGORY", "Just Chatting"),
		ApiURLs: twitchAPIConfig{
			TwitchAPI: envOrDefault("TWITCH_API_URL", "https://api.twitch.tv/helix"),
		},
		Spotify: spotifyConfig{
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RewardID:     os.Getenv("SPOTIFY_REWARD_ID"),
			Mode:         envOrDefault("SPOTIFY_MODE", "disabled"),
		},
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
