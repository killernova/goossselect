package goossselect

import (
	"encoding/xml"
	"errors"

	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var InvalidValueErr = errors.New("invalid values for selector")

type Selector struct {
	XMLName            xml.Name `xml:"SelectRequest"`
	Expression         string   `xml:"Expression"`
	CompressionType    string   `xml:"InputSerialization>CompressionType"`
	FileHeaderInfo     string   `xml:"InputSerialization>CSV>FileHeaderInfo"`
	InRecordDelimiter  string   `xml:"InputSerialization>CSV>RecordDelimiter,omitempty"`
	InFieldDelimiter   string   `xml:"InputSerialization>CSV>FieldDelimiter,omitempty"`
	QuoteCharacter     string   `xml:"InputSerialization>CSV>QuoteCharacter,omitempty"`
	CommentCharacter   string   `xml:"InputSerialization>CSV>CommentCharacter,omitempty"`
	Range              string   `xml:"InputSerialization>CSV>Range,omitempty"`
	OutputRawData      bool     `xml:"OutputSerialization>OutputRawData,omitempty"`
	OutRecordDelimiter string   `xml:"OutputSerialization>CSV>RecordDelimiter,omitempty"`
	OutFieldDelimiter  string   `xml:"OutputSerialization>CSV>FieldDelimiter,omitempty"`
	KeepAllColumns     bool     `xml:"OutputSerialization>CSV>KeepAllColumns,omitempty"`
}

type Meta struct {
	XMLName             xml.Name `xml:"CsvMetaRequest"`
	OverwriteIfExisting bool     `xml:"Expression,omitempty"`
	RecordDelimiter     string   `xml:"InputSerialization>CSV>RecordDelimiter,omitempty"`
	FieldDelimiter      string   `xml:"InputSerialization>CSV>FieldDelimiter,omitempty"`
	QuoteCharacter      string   `xml:"InputSerialization>CSV>QuoteCharacter,omitempty"`
	CompressionType     string   `xml:"InputSerialization>CompressionType"`
}

type MetaResponse struct {
	Lines         int64
	Columns       int64
	Splits        int64
	ContentLength int64
}

// Conn oss conn
type Conn struct {
	Config *Config
	Client *http.Client
}

// Config oss configure
type Config struct {
	AccessKeyID           string // accessId
	AccessKeySecret       string // accessKey
	Endpoint              string
	Timeout               time.Duration
	KeepAlive             time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
}

func NewMeta(options ...func(*Meta) error) (*Meta, error) {
	m := &Meta{
		CompressionType: "None",
	}

	for _, f := range options {
		if err := f(m); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func NewSelector(query string, options ...func(*Selector) error) (*Selector, error) {
	exp := toBase64(query)
	s := &Selector{
		CompressionType: "None",
		OutputRawData:   true,
		FileHeaderInfo:  "None",
		Expression:      exp,
	}

	for _, f := range options {
		if err := f(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func NewConfig(options ...func(*Config)) *Config {
	c := &Config{
		Timeout:               10 * time.Second,
		KeepAlive:             10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	}

	for _, f := range options {
		f(c)
	}
	return c
}

func toXML(s interface{}) ([]byte, error) {
	ss, err := xml.Marshal(s)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func getResponse(bucket, objectKey, action string, content []byte, config *Config) (*http.Response, error) {
	now := getGmtIso8601(time.Now().Unix())
	canonicalizedResource := "/" + bucket + "/" + objectKey + "?x-oss-process=csv/" + action
	url := fmt.Sprintf(("https://%s.%s/%s?x-oss-process=csv/" + action), bucket, config.Endpoint, objectKey)
	c := &http.Client{
		Transport: &http.Transport{
			// Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Dial: (&net.Dialer{
				Timeout:   config.Timeout,
				KeepAlive: config.KeepAlive,
			}).Dial,
			TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
			ResponseHeaderTimeout: config.ResponseHeaderTimeout,
			ExpectContinueTimeout: config.ExpectContinueTimeout,
		},
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Date", now)
	req.Header.Set("Content-Length", strconv.Itoa(len(content)))
	//req.Header.Set("Content-MD5", hex.EncodeToString(auth.CalcMD5(content)))

	conn := Conn{
		Config: config,
		Client: c,
	}
	conn.SignHeader(req, canonicalizedResource)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Selector) SelectQuery(bucket, objectKey string, config *Config) ([]byte, error) {
	content, err := toXML(s)
	if err != nil {
		return nil, err
	}
	resp, err := getResponse(bucket, objectKey, "select", content, config)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, errors.New("errors returned from ali oss")
	}
	return body, nil
}

func (s *Selector) SelectQueryToFile(bucket, objectKey string, config *Config, f *os.File) error {
	content, err := toXML(s)
	if err != nil {
		return err
	}
	resp, err := getResponse(bucket, objectKey, "select", content, config)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("errors returned from ali oss")
	}
	return nil
}

func (m *Meta) SelectMeta(bucket, objectKey string, config *Config) (*MetaResponse, error) {
	content, err := toXML(m)
	if err != nil {
		return nil, err
	}
	resp, err := getResponse(bucket, objectKey, "meta", content, config)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	lines, err := strconv.ParseInt(resp.Header["X-Oss-Select-Csv-Rows"][0], 10, 64)
	if err != nil {
		return nil, err
	}
	columns, err := strconv.ParseInt(resp.Header["X-Oss-Select-Csv-Columns"][0], 10, 64)
	if err != nil {
		return nil, err
	}
	splits, err := strconv.ParseInt(resp.Header["X-Oss-Select-Csv-Splits"][0], 10, 64)
	if err != nil {
		return nil, err
	}
	mr := &MetaResponse{
		Lines:         lines,
		Columns:       columns,
		Splits:        splits,
		ContentLength: resp.ContentLength,
	}
	return mr, nil
}
