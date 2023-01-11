package setlistfm

import (
	"context"
	"fmt"

	"github.com/jmg-duarte/setlistfm"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableSetlistFMSetlist(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "setlistfm_setlist",
		Description: "Setlist.fm setlist",
		List: &plugin.ListConfig{
			Hydrate: listSetlistFMSetlist,
		},
		//Get: &plugin.GetConfig{
		//Hydrate:    getFlyApp,
		//KeyColumns: plugin.SingleColumn("name"),
		//},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The setlist identifier.",
				Type:        proto.ColumnType_STRING,
			},
			// whole artist object??  id?
			// {
			// 	Name:        "artist",
			// 	Description: "The artist.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			{
				Name:        "artist_name",
				Description: "The artist name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Artist.Name"),
			},

			{
				Name:        "event_date",
				Description: "The setlist event date.",
				//to do - transform Type:        proto.ColumnType_TIMESTAMP,
				Type: proto.ColumnType_STRING,
			},

			{
				Name:        "set_name",
				Description: "The setlist set name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "set_number",
				Description: "The setlist set number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "set_encore",
				Description: "The encore number.",
				Type:        proto.ColumnType_INT,
			},

			{
				Name:        "song_name",
				Description: "The song name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "song_number",
				Description: "The song number in the set.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "song_info",
				Description: "The song info.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "original_artist_name",
				Description: "The original artist.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SongCover.Name"),
			},
			{
				Name:        "guest_artist_name",
				Description: "Guest artist on the song.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SongWith.Name"),
			},

			{
				Name:        "tour_name",
				Description: "The setlist tour name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Tour.Name"),
			},

			// whole venue object?? venue id?
			// {
			// 	Name:        "venue",
			// 	Description: "The setlist venue.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			{
				Name:        "venue_name",
				Description: "The setlist venue name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Venue.Name"),
			},

			{
				Name:        "venue_city",
				Description: "The setlist venue city.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Venue.City.Name"),
			},
			{
				Name:        "venue_country",
				Description: "The setlist venue country.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Venue.City.Country.Name"),
			},
			{
				Name:        "venue_state",
				Description: "The setlist venue state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Venue.City.State"),
			},

			{
				Name:        "setlist_info",
				Description: "The setlist info.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Info"),
			},

			{
				Name:        "url",
				Description: "The setlist url.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id",
				Description: "The setlist version identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated",
				Description: "The setlist last updated date.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// {
			// 	Name:        "raw",
			// 	Description: "The raw response",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromValue(),
			// },
		},
	}
}

//// LIST FUNCTION

func listSetlistFMSetlist(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	setlistfmConfig := GetConfig(d.Connection)
	if setlistfmConfig.Token == nil {
		return nil, fmt.Errorf("token must be passed")
	}
	client := setlistfm.NewClient(*setlistfmConfig.Token)

	artist := "6faa7ca7-0d99-4a5e-bfa6-1fd5037520c6" //grateful dead
	//artist := "65f4f0c5-ef9e-490c-aee3-909e7ae6b2ab" #metallica

	for i := 0; true; i++ {
		conn, err := client.ArtistSetlistsByMBID(ctx, artist, i)

		if err != nil {
			plugin.Logger(ctx).Error("setlistfm.listSetlistFMSetlist", "connection_error", err)
			return nil, err
		}

		rows := buildSetlistRows(conn.Setlists)
		for _, row := range rows {
			d.StreamLeafListItem(ctx, row)
		}

		if i > conn.Total/conn.ItemsPerPage {
			break
		}
	}

	return nil, nil
}



func buildSetlistRows(setlistRaw []setlistfm.Setlist) []SetlistRow {
	rows := []SetlistRow{}

	for _, setlist := range setlistRaw {
		row := SetlistRow{
			ID:          setlist.ID,
			URL:         setlist.URL,
			VersionID:   setlist.VersionID,
			EventDate:   setlist.EventDate,
			LastUpdated: setlist.LastUpdated,
			SetlistInfo: setlist.Info,

			Artist: setlist.Artist,
			Venue:  setlist.Venue,
			Tour:   setlist.Tour,
		}

		// sets
		for setNumber, set := range setlist.Sets.Set {
			row.SetEncore = set.Encore
			row.SetName = set.Name
			row.SetNumber = setNumber + 1

			//songs
			for songNumber, song := range set.Song {
				row.SongCover = song.Cover
				row.SongInfo = song.Info
				row.SongName = song.Name
				row.SongWith = song.With
				row.SongTape = song.Tape
				row.SongNumber = songNumber + 1
				rows = append(rows, row)

			}
		}

	}
	return rows
}

/*
	type Setlist struct {
		Artist      Artist `json:"artist"`
		Venue       Venue  `json:"venue"`
		Tour        Tour   `json:"tour"`
		Sets        Sets   `json:"sets"`
		Info        string `json:"info"`
		URL         string `json:"url"`
		ID          string `json:"id"`
		VersionID   string `json:"versionId"`
		EventDate   string `json:"eventDate"`
		LastUpdated string `json:"lastUpdated"`
	}

// Artist - This class represents an artist.
// An artist is a musician or a group of musicians.
// Each artist has a definite Musicbrainz Identifier (MBID)
// with which the artist can be uniquely identified.

	type Artist struct {
		MBID           string `json:"mbid"`
		TMID           int    `json:"tmid"`
		Name           string `json:"name"`
		SortName       string `json:"sortName"`
		Disambiguation string `json:"disambiguation"`
		URL            string `json:"url"`
	}

// Venue - Venues are places where concerts take place.
// They usually consist of a venue name and a city -
// but there are also some venues that do not have a city attached yet.
// In such a case, the city simply isn't set and the city and country
// may (but do not have to) be in the name.

	type Venue struct {
		City City   `json:"city"`
		URL  string `json:"url"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}

// City -  	This class represents a city where Venues are located.
// Most of the original city data was taken from Geonames.org.

	type City struct {
		ID        string      `json:"id"`
		Name      string      `json:"name"`
		StateCode string      `json:"stateCode"`
		State     string      `json:"state"`
		Coords    Coordinates `json:"coords"`
		Country   Country     `json:"country"`
	}

// Tour - The tour a setlist was a part of.

	type Tour struct {
		Name string `json:"name"`
	}

// Set - A setlist consists of different (at least one) sets.
// Sets can either be sets as defined in the Guidelines or encores.

	type Set struct {
		Name   string `json:"name"`
		Encore int    `json:"encore"`
		Song   []Song `json:"song"`
	}

// Song - This class represents a song that is part of a Set.

	type Song struct {
		Name  string `json:"name"`
		With  Artist `json:"with,omitempty"`
		Cover Artist `json:"cover,omitempty"`
		Info  string `json:"info"`
		Tape  bool   `json:"tape"`
	}
*/
type SetlistRow struct {
	Artist setlistfm.Artist `json:"artist"`

	Venue setlistfm.Venue `json:"venue"`

	Tour setlistfm.Tour `json:"tour"`

	SetlistInfo string `json:"setlistInfo"`
	URL         string `json:"url"`
	ID          string `json:"id"`
	VersionID   string `json:"versionId"`
	EventDate   string `json:"eventDate"`
	LastUpdated string `json:"lastUpdated"`

	//	Sets        Sets   `json:"sets"`
	SetName   string `json:"setName"`
	SetEncore int    `json:"setEncore"`
	SetNumber int    `json:"setNumber"`

	SongName   string           `json:"songName"`
	SongNumber int              `json:"songNumber"`
	SongWith   setlistfm.Artist `json:"songWith"`
	SongCover  setlistfm.Artist `json:"songCover"`
	SongInfo   string           `json:"songInfo"`
	SongTape   bool             `json:"songTape"`
}
