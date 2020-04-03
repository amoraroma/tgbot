package loaders

import (
	"net/url"

	"github.com/camelva/erzo/parsers"
)

type Loader interface {
	Name() string
	Bin() string
	Compatible(format parsers.Format) bool
	Get(*url.URL, string) error
	UpdateTags(string, []string) error
	AddThumbnail(string, string) error
}
