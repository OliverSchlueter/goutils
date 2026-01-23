package sitemap

import "encoding/xml"

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

// Url represents a single URL entry in the sitemap.
// It includes the location, last modification date, change frequency, and priority.
// All fields except Loc are optional.
type Url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type UrlProvider func() []Url
