package sitemap

import (
	"encoding/xml"
	"io"

	jtutils "github.com/johnietre/utils/go"
)

type SitemapIndex struct {
	XMLName  jtutils.Unit   `xml:"sitemapindex"`
	Sitemaps []SitemapEntry `xml:"sitemap"`
}

type SitemapEntry struct {
	XMLName jtutils.Unit `xml:"sitemap"`
	Loc     string       `xml:"loc"`
	LastMod string       `xml:"lastmod,omitempty"`
}

type Sitemap struct {
	XMLName jtutils.Unit `xml:"urlset"`
	Urls    []UrlEntry   `xml:"url"`
}

type UrlEntry struct {
	XMLName    jtutils.Unit `xml:"url"`
	Loc        string       `xml:"loc"`
	LastMod    string       `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq   `xml:"changefreq,omitempty"`
	Priority   float32      `xml:"priority,omitempty"`
}

type ChangeFreq string

const (
	ChangeFreqAlways  ChangeFreq = "always"
	ChangeFreqHourly  ChangeFreq = "hourly"
	ChangeFreqWeekly  ChangeFreq = "weekly"
	ChangeFreqMonthly ChangeFreq = "monthly"
	ChangeFreqYearly  ChangeFreq = "yearly"
	ChangeFreqNever   ChangeFreq = "never"
)

type Encoder struct {
	*xml.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		xml.NewEncoder(w),
	}
}

func (e *Encoder) WriteHeader() error {
	return e.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0" encoding="UTF-8"`),
	})
}

func (e *Encoder) EncodeWithHeader(v any) error {
	err := e.WriteHeader()
	if err != nil {
		return err
	}
	return e.Encode(v)
}

type Test = io.Writer
