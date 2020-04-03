package erzo

import (
	"github.com/camelva/erzo/engine"
	_ "github.com/camelva/erzo/loaders/ffmpeg"
	_ "github.com/camelva/erzo/parsers/soundcloud"
)

type ErrNotURL struct {
	engine.ErrNotURL
}
type ErrUnsupportedService struct {
	engine.ErrUnsupportedService
}
type ErrUnsupportedType struct {
	engine.ErrUnsupportedType
}
type ErrCantFetchInfo struct {
	engine.ErrCantFetchInfo
}
type ErrUnsupportedProtocol struct {
	engine.ErrUnsupportedProtocol
}
type ErrDownloadingError struct {
	engine.ErrDownloadingError
}
type ErrUndefined struct {
	engine.ErrUndefined
}

type options struct {
	output   string
	truncate bool
}

type Option interface {
	apply(*options)
}

type truncateOption bool

func (opt truncateOption) apply(opts *options) {
	opts.truncate = bool(opt)
}
func Truncate(b bool) Option {
	return truncateOption(b)
}

type outputOption string

func (opt outputOption) apply(opts *options) {
	opts.output = string(opt)
}
func Output(s string) Option {
	return outputOption(s)
}

// Get process given url and download song from it.
// @message - url to process
// @options:
// 		Truncate(true|false) - clear output folder before processing
//		Output(string)		 - change output folder
// Return filename or one of the following errors:
// ErrNotURL if there is no urls in your message
// ErrUnsupportedService if url belongs to unsupported service
// ErrUnsupportedType if service supported but certain type - not yet
// ErrCantFetchInfo if fatal error occurred while extracting info from url
// ErrUnsupportedProtocol if there is no downloader for this format
// ErrDownloadingError if fatal error occurred while downloading song
// ErrUndefined any other errors
func Get(message string, opts ...Option) (string, error) {
	options := options{
		output:   "out",
		truncate: false,
	}
	for _, o := range opts {
		o.apply(&options)
	}
	e := engine.New(
		options.output,
		options.truncate,
	)
	r, err := e.Process(message)
	if err != nil {
		var convertedErr error
		switch err.(type) {
		case engine.ErrNotURL:
			convertedErr = ErrNotURL{err.(engine.ErrNotURL)}
		case engine.ErrUnsupportedService:
			convertedErr = ErrUnsupportedService{err.(engine.ErrUnsupportedService)}
		case engine.ErrUnsupportedType:
			convertedErr = ErrUnsupportedType{err.(engine.ErrUnsupportedType)}
		case engine.ErrCantFetchInfo:
			convertedErr = ErrCantFetchInfo{err.(engine.ErrCantFetchInfo)}
		case engine.ErrUnsupportedProtocol:
			convertedErr = ErrUnsupportedProtocol{err.(engine.ErrUnsupportedProtocol)}
		case engine.ErrDownloadingError:
			convertedErr = ErrDownloadingError{err.(engine.ErrDownloadingError)}
		case engine.ErrUndefined:
			convertedErr = ErrUndefined{err.(engine.ErrUndefined)}
		default:
			convertedErr = ErrUndefined{engine.ErrUndefined{}}
		}
		return "", convertedErr
	}
	return r, nil
}
