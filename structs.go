package xtreamcodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Timestamp is a helper struct to convert unix timestamp ints and strings to time.Time.
type Timestamp struct {
	time.Time
	quoted bool
}

// MarshalJSON returns the Unix timestamp as a string.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.quoted {
		return []byte(`"` + strconv.FormatInt(t.Time.Unix(), 10) + `"`), nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// UnmarshalJSON converts the int or string to a Unix timestamp.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	// Timestamps are sometimes quoted, sometimes not, lets just always remove quotes just in case...
	t.quoted = strings.Contains(string(b), `"`)
	ts, err := strconv.Atoi(strings.Replace(string(b), `"`, "", -1))
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(ts), 0)
	return nil
}

// Timezone is a helper struct to convert timezone names into time.Location.
type Timezone struct {
	time.Location
}

// MarshalJSON returns the time.Location name as a string.
func (t Timezone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Location.String() + `"`), nil
}

// UnmarshalJSON converts a JSON string to a time.Location.
func (t *Timezone) UnmarshalJSON(b []byte) error {
	unquoted, unquotedErr := strconv.Unquote(strings.Replace(string(b), `\`, "", 1))
	if unquotedErr != nil {
		return unquotedErr
	}
	tz, err := time.LoadLocation(unquoted)
	if err != nil {
		return err
	}
	t.Location = *tz
	return nil
}

// ConvertibleBoolean is a helper type to allow JSON documents using 0/1 or "true" and "false" be converted to bool.
type ConvertibleBoolean struct {
	bool
	quoted bool
}

// MarshalJSON returns a 0 or 1 depending on bool state.
func (bit ConvertibleBoolean) MarshalJSON() ([]byte, error) {
	var bitSetVar int8
	if bit.bool {
		bitSetVar = 1
	}

	if bit.quoted {
		return json.Marshal(fmt.Sprint(bitSetVar))
	}

	return json.Marshal(bitSetVar)
}

// UnmarshalJSON converts a 0, 1, true or false into a bool
func (bit *ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	bit.quoted = strings.Contains(string(data), `"`)
	// Bools as ints are sometimes quoted, sometimes not, lets just always remove quotes just in case...
	asString := strings.Replace(string(data), `"`, "", -1)
	if asString == "1" || asString == "true" {
		bit.bool = true
	} else if asString == "0" || asString == "false" {
		bit.bool = false
	} else {
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", asString)
	}
	return nil
}

// jsonInt is a int64 which unmarshals from JSON
// as either unquoted or quoted (with any amount
// of internal leading/trailing whitespace).
// Originally found at https://bit.ly/2NkJ0SK and
// https://play.golang.org/p/KNPxDL1yqL
type jsonInt int64

func (f jsonInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(f))
}

func (f *jsonInt) UnmarshalJSON(data []byte) error {
	var v int64

	data = bytes.Trim(data, `" `)

	err := json.Unmarshal(data, &v)
	*f = jsonInt(v)
	return err
}

// ServerInfo describes the state of the Xtream-Codes server.
type ServerInfo struct {
	HTTPSPort    int       `json:"https_port,string"`
	Port         int       `json:"port,string"`
	Process      bool      `json:"process"`
	RTMPPort     int       `json:"rtmp_port,string"`
	Protocol     string    `json:"server_protocol"`
	TimeNow      string    `json:"time_now"`
	TimestampNow Timestamp `json:"timestamp_now,string"`
	Timezone     Timezone  `json:"timezone"`
	URL          string    `json:"url"`
}

// UserInfo is the current state of the user as it relates to the Xtream-Codes server.
type UserInfo struct {
	ActiveConnections    int                `json:"active_cons,string"`
	AllowedOutputFormats []string           `json:"allowed_output_formats"`
	Auth                 ConvertibleBoolean `json:"auth"`
	CreatedAt            Timestamp          `json:"created_at"`
	ExpDate              *Timestamp         `json:"exp_date"`
	IsTrial              ConvertibleBoolean `json:"is_trial,string"`
	MaxConnections       int                `json:"max_connections,string"`
	Message              string             `json:"message"`
	Password             string             `json:"password"`
	Status               string             `json:"status"`
	Username             string             `json:"username"`
}

// AuthenticationResponse is a container for what the server returns after the initial authentication.
type AuthenticationResponse struct {
	ServerInfo ServerInfo `json:"server_info"`
	UserInfo   UserInfo   `json:"user_info"`
}

// Category describes a grouping of Stream.
type Category struct {
	ID     int    `json:"category_id,string"`
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`

	// Set by us, not Xtream.
	Type string `json:"-"`
}

// Stream is a streamble video source.
type Stream struct {
	Added              Timestamp `json:"added"`
	CategoryID         int       `json:"category_id,string"`
	ContainerExtension string    `json:"container_extension"`
	CustomSid          string    `json:"custom_sid"`
	DirectSource       string    `json:"direct_source"`
	EPGChannelID       string    `json:"epg_channel_id"`
	Icon               string    `json:"stream_icon"`
	ID                 int       `json:"stream_id"`
	Name               string    `json:"name"`
	Number             int       `json:"num"`
	Rating             string    `json:"rating"`
	Rating5based       int64     `json:"rating_5based"`
	TVArchive          int       `json:"tv_archive"`
	TVArchiveDuration  jsonInt   `json:"tv_archive_duration"`
	Type               string    `json:"stream_type"`
}

// SeriesInfo contains information about a TV series.
type SeriesInfo struct {
	BackdropPath   JSONStringSlice `json:"backdrop_path"`
	Cast           string          `json:"cast"`
	CategoryID     int             `json:"category_id,string"`
	Cover          string          `json:"cover"`
	Director       string          `json:"director"`
	EpisodeRunTime string          `json:"episode_run_time"`
	Genre          string          `json:"genre"`
	LastModified   Timestamp       `json:"last_modified"`
	Name           string          `json:"name"`
	Num            int             `json:"num"`
	Plot           string          `json:"plot"`
	Rating         int             `json:"rating,string"`
	Rating5        int             `json:"rating_5based"`
	ReleaseDate    string          `json:"releaseDate"`
	SeriesID       int             `json:"series_id"`
	StreamType     string          `json:"stream_type"`
	YoutubeTrailer string          `json:"youtube_trailer"`
}

// VideoOnDemandInfo contains information about a video on demand stream.
type VideoOnDemandInfo struct {
	Info struct {
		Audio          FFMPEGStreamInfo `json:"audio"`
		BackdropPath   []string         `json:"backdrop_path"`
		Bitrate        int              `json:"bitrate"`
		Cast           string           `json:"cast"`
		Director       string           `json:"director"`
		Duration       string           `json:"duration"`
		DurationSecs   int              `json:"duration_secs"`
		Genre          string           `json:"genre"`
		MovieImage     string           `json:"movie_image"`
		Plot           string           `json:"plot"`
		Rating         string           `json:"rating"`
		Releasedate    string           `json:"releasedate"`
		TmdbID         string           `json:"tmdb_id"`
		Video          FFMPEGStreamInfo `json:"video"`
		YoutubeTrailer string           `json:"youtube_trailer"`
	} `json:"info"`
	MovieData struct {
		Added              Timestamp `json:"added"`
		CategoryID         int       `json:"category_id,string"`
		ContainerExtension string    `json:"container_extension"`
		CustomSid          string    `json:"custom_sid"`
		DirectSource       string    `json:"direct_source"`
		Name               string    `json:"name"`
		StreamID           int       `json:"stream_id"`
	} `json:"movie_data"`
}

type epgContainer struct {
	EPGListings []EPGInfo `json:"epg_listings"`
}

// EPGInfo describes electronic programming guide information of a stream.
type EPGInfo struct {
	ChannelID      string             `json:"channel_id"`
	Description    Base64Value        `json:"description"`
	End            string             `json:"end"`
	EPGID          int                `json:"epg_id,string"`
	HasArchive     ConvertibleBoolean `json:"has_archive"`
	ID             int                `json:"id,string"`
	Lang           string             `json:"lang"`
	NowPlaying     ConvertibleBoolean `json:"now_playing"`
	Start          string             `json:"start"`
	StartTimestamp Timestamp          `json:"start_timestamp"`
	StopTimestamp  Timestamp          `json:"stop_timestamp"`
	Title          Base64Value        `json:"title"`
}

type JSONStringSlice struct {
	Slice        []string `json:"-"`
	SingleString bool     `json:"-"`
}

func (b JSONStringSlice) MarshalJSON() ([]byte, error) {
	if !b.SingleString {
		return json.Marshal(b.Slice)
	}
	return json.Marshal(b.Slice[0])
}

func (b *JSONStringSlice) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		data = append([]byte(`[`), data...)
		data = append(data, []byte(`]`)...)
		b.SingleString = true
	}

	return json.Unmarshal(data, &b.Slice)
}