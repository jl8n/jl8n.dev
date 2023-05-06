package structs

import "time"

type NowPlaying struct {
	Artist    string
	Album     string
	Track     string
	AlbumMBID string
	Coverart  string
}

type LBPlayingNow struct {
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

// https://musicbrainz.org/ws/2/recording/?query=artist:{artist}%20AND%20recording:{track}&fmt=json
type LBTrackInfo struct {
	Created    time.Time `json:"created"`
	Count      int       `json:"count"`
	Offset     int       `json:"offset"`
	Recordings []struct {
		ID           string `json:"id"`
		Score        int    `json:"score"`
		Title        string `json:"title"`
		Length       int    `json:"length"`
		Video        any    `json:"video"`
		ArtistCredit []struct {
			Name   string `json:"name"`
			Artist struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				SortName       string `json:"sort-name"`
				Disambiguation string `json:"disambiguation"`
			} `json:"artist"`
		} `json:"artist-credit"`
		FirstReleaseDate string `json:"first-release-date"`
		Releases         []struct {
			ID           string `json:"id"`
			StatusID     string `json:"status-id"`
			Count        int    `json:"count"`
			Title        string `json:"title"`
			Status       string `json:"status"`
			ReleaseGroup struct {
				ID            string `json:"id"`
				TypeID        string `json:"type-id"`
				PrimaryTypeID string `json:"primary-type-id"`
				Title         string `json:"title"`
				PrimaryType   string `json:"primary-type"`
			} `json:"release-group"`
			Date          string `json:"date"`
			Country       string `json:"country"`
			ReleaseEvents []struct {
				Date string `json:"date"`
				Area struct {
					ID            string   `json:"id"`
					Name          string   `json:"name"`
					SortName      string   `json:"sort-name"`
					Iso31661Codes []string `json:"iso-3166-1-codes"`
				} `json:"area"`
			} `json:"release-events"`
			TrackCount int `json:"track-count"`
			Media      []struct {
				Position int    `json:"position"`
				Format   string `json:"format"`
				Track    []struct {
					ID     string `json:"id"`
					Number string `json:"number"`
					Title  string `json:"title"`
					Length int    `json:"length"`
				} `json:"track"`
				TrackCount  int `json:"track-count"`
				TrackOffset int `json:"track-offset"`
			} `json:"media"`
		} `json:"releases"`
	} `json:"recordings"`
}

type LastfmNowPlaying struct {
	Recenttracks struct {
		Track []struct {
			Artist struct {
				URL   string `json:"url"`
				Name  string `json:"name"`
				Image []struct {
					Size string `json:"size"`
					Text string `json:"#text"`
				} `json:"image"`
				Mbid string `json:"mbid"`
			} `json:"artist"`
			Mbid  string `json:"mbid"`
			Name  string `json:"name"`
			Image []struct {
				Size string `json:"size"`
				Text string `json:"#text"`
			} `json:"image"`
			Streamable string `json:"streamable"`
			Album      struct {
				Mbid string `json:"mbid"`
				Text string `json:"#text"`
			} `json:"album"`
			URL  string `json:"url"`
			Attr struct {
				Nowplaying string `json:"nowplaying"`
			} `json:"@attr,omitempty"`
			Loved string `json:"loved"`
			Date  struct {
				Uts  string `json:"uts"`
				Text string `json:"#text"`
			} `json:"date,omitempty"`
		} `json:"track"`
		Attr struct {
			User       string `json:"user"`
			TotalPages string `json:"totalPages"`
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			Total      string `json:"total"`
		} `json:"@attr"`
	} `json:"recenttracks"`
}

type MBAlbumInfo struct {
	Created  time.Time `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []struct {
		ID                 string `json:"id"`
		Score              int    `json:"score"`
		StatusID           string `json:"status-id"`
		PackagingID        string `json:"packaging-id"`
		Count              int    `json:"count"`
		Title              string `json:"title"`
		Status             string `json:"status"`
		Packaging          string `json:"packaging"`
		TextRepresentation struct {
			Language string `json:"language"`
			Script   string `json:"script"`
		} `json:"text-representation"`
		ArtistCredit []struct {
			Name   string `json:"name"`
			Artist struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				SortName string `json:"sort-name"`
			} `json:"artist"`
		} `json:"artist-credit"`
		ReleaseGroup struct {
			ID            string `json:"id"`
			TypeID        string `json:"type-id"`
			PrimaryTypeID string `json:"primary-type-id"`
			Title         string `json:"title"`
			PrimaryType   string `json:"primary-type"`
		} `json:"release-group"`
		Date          string `json:"date"`
		Country       string `json:"country"`
		ReleaseEvents []struct {
			Date string `json:"date"`
			Area struct {
				ID            string   `json:"id"`
				Name          string   `json:"name"`
				SortName      string   `json:"sort-name"`
				Iso31661Codes []string `json:"iso-3166-1-codes"`
			} `json:"area"`
		} `json:"release-events"`
		Barcode   string `json:"barcode"`
		Asin      string `json:"asin"`
		LabelInfo []struct {
			Label struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"label"`
		} `json:"label-info"`
		TrackCount int `json:"track-count"`
		Media      []struct {
			Format     string `json:"format"`
			DiscCount  int    `json:"disc-count"`
			TrackCount int    `json:"track-count"`
		} `json:"media"`
		Tags []struct {
			Count int    `json:"count"`
			Name  string `json:"name"`
		} `json:"tags"`
	} `json:"releases"`
}

type CoverartArchive struct {
	Images []struct {
		Approved   bool   `json:"approved"`
		Back       bool   `json:"back"`
		Comment    string `json:"comment"`
		Edit       int    `json:"edit"`
		Front      bool   `json:"front"`
		ID         int64  `json:"id"`
		Image      string `json:"image"`
		Thumbnails struct {
			Num250  string `json:"250"`
			Num500  string `json:"500"`
			Num1200 string `json:"1200"`
			Large   string `json:"large"`
			Small   string `json:"small"`
		} `json:"thumbnails"`
		Types []string `json:"types"`
	} `json:"images"`
	Release string `json:"release"`
}
