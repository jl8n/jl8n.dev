package structs

type ListenbrainzResponse struct {
	Payload struct {
		Count   int `json:"count"`
		Listens []struct {
			PlayingNow    bool `json:"playing_now"`
			TrackMetadata struct {
				AdditionalInfo struct {
					Duration                int    `json:"duration"`
					MusicServiceName        string `json:"music_service_name"`
					OriginURL               string `json:"origin_url"`
					SubmissionClient        string `json:"submission_client"`
					SubmissionClientVersion string `json:"submission_client_version"`
				} `json:"additional_info"`
				ArtistName  string `json:"artist_name"`
				ReleaseName string `json:"release_name"`
				TrackName   string `json:"track_name"`
			} `json:"track_metadata"`
		} `json:"listens"`
		PlayingNow bool   `json:"playing_now"`
		UserID     string `json:"user_id"`
	} `json:"payload"`
}

type NowPlaying struct {
	Artist string
	Album  string
	Track  string
}
