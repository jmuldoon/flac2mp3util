package convert

import ("os")

type Converter interface {
	GetFiles() error
}

type Conversion struct {
	Files *[]FileInfo
}

type MetaTag struct {
	Artist      string
	Title       string
	Album       string
	Date        string
	TrackNumber int
	Genre       string
	Performer   string
	Composer    string
	Lyrics      string
	AlbumArtist string
	DiscNumber  int
	TotalDiscs  int
	TotalTracks int
	Comment     string
}

func (c *Conversion) GetFiles() error {
	c.Files := os.ReadDir
}
