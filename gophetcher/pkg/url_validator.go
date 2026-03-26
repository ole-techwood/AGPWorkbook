package pkg

import (
	"fmt"
	"net/url"
)

type iUrlValidator interface {
	isURLEmpty(candidate string) error
	isURLMalformed(candidate string) (*url.URL, error)
	isURLSchemeValid(candidate *url.URL) error
}

type URLValidator struct {
	iUrlValidator iUrlValidator
}

func (uv *URLValidator) Validate(candidate string) (*url.URL, error) {
	err := uv.iUrlValidator.isURLEmpty(candidate)
	if err != nil {
		return nil, err
	}

	url, err := uv.iUrlValidator.isURLMalformed(candidate)
	if err != nil {
		return nil, err
	}

	err = uv.iUrlValidator.isURLSchemeValid(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

// WebURLValidator

type WebURLValidator struct {
	URLValidator
}

func NewWebURLValidator() *WebURLValidator {
	wuv := &WebURLValidator{}
	wuv.iUrlValidator = wuv
	return wuv
}

func (wuv *WebURLValidator) isURLEmpty(candidate string) error {
	if candidate == "" {
		return fmt.Errorf("error: URL is empty")
	}
	return nil
}

func (wuv *WebURLValidator) isURLMalformed(candidate string) (*url.URL, error) {
	parsed, err := url.ParseRequestURI(candidate)
	if err != nil {
		return nil, fmt.Errorf("error: malformed URL: %s", candidate)
	}
	return parsed, nil
}

func (wuv *WebURLValidator) isURLSchemeValid(candidate *url.URL) error {
	if candidate.Scheme != "http" && candidate.Scheme != "https" {
		return fmt.Errorf("error: unsupported URL scheme: %s", candidate.Scheme)
	}
	return nil
}

// FileURLValidator

type FileURLValidator struct {
	WebURLValidator
}

func NewFileURLValidator() *FileURLValidator {
	fuv := &FileURLValidator{}
	fuv.iUrlValidator = fuv
	return fuv
}

func (fuv *FileURLValidator) isURLMalformed(candidate string) (*url.URL, error) {
	parsed, err := url.Parse(candidate)
	if err != nil {
		return nil, fmt.Errorf("error: malformed URL: %s", candidate)
	}
	return parsed, nil
}

func (fuv *FileURLValidator) isURLSchemeValid(candidate *url.URL) error {
	if candidate.Scheme != "file" {
		return fmt.Errorf("error: unsupported URL scheme: %s", candidate.Scheme)
	}
	if candidate.Path == "" {
		return fmt.Errorf("error: file URL must include a path: %s", candidate.String())
	}
	return nil
}
