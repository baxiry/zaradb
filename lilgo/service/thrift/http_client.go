package thrift

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var DefaultHttpClient *http.Client = http.DefaultClient

type THttpClient struct {
	client             *http.Client
	response           *http.Response
	url                *url.URL
	requestBuffer      *bytes.Buffer
	header             http.Header
	nsecConnectTimeout int64
	nsecReadTimeout    int64
}

type THttpClientTransportFactory struct {
	options THttpClientOptions
	url     string
}

func (p *THttpClientTransportFactory) GetTransport(trans TTransport) (TTransport, error) {
	if trans != nil {
		t, ok := trans.(*THttpClient)
		if ok && t.url != nil {
			return NewTHttpClientWithOptions(t.url.String(), p.options)
		}
	}
	return NewTHttpClientWithOptions(p.url, p.options)
}

type THttpClientOptions struct {
	Client *http.Client
}

func NewTHttpClientTransportFactory(url string) *THttpClientTransportFactory {
	return NewTHttpClientTransportFactoryWithOptions(url, THttpClientOptions{})
}

func NewTHttpClientTransportFactoryWithOptions(url string, options THttpClientOptions) *THttpClientTransportFactory {
	return &THttpClientTransportFactory{url: url, options: options}
}

func NewTHttpClientWithOptions(urlstr string, options THttpClientOptions) (TTransport, error) {
	parsedURL, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 0, 2048)
	client := options.Client
	if client == nil {
		client = DefaultHttpClient
	}
	httpHeader := map[string][]string{"Content-Type": {"application/x-thrift"}}
	return &THttpClient{client: client, url: parsedURL, requestBuffer: bytes.NewBuffer(buf), header: httpHeader}, nil
}

func NewTHttpClient(urlstr string) (TTransport, error) {
	return NewTHttpClientWithOptions(urlstr, THttpClientOptions{})
}

func (p *THttpClient) SetHeader(key string, value string) {
	p.header.Add(key, value)
}

func (p *THttpClient) GetHeader(key string) string {
	return p.header.Get(key)
}

func (p *THttpClient) DelHeader(key string) {
	p.header.Del(key)
}

func (p *THttpClient) Open() error {
	return nil
}

func (p *THttpClient) IsOpen() bool {
	return p.response != nil || p.requestBuffer != nil
}

func (p *THttpClient) closeResponse() error {
	var err error
	if p.response != nil && p.response.Body != nil {
		io.Copy(ioutil.Discard, p.response.Body)

		err = p.response.Body.Close()
	}

	p.response = nil
	return err
}

func (p *THttpClient) Close() error {
	if p.requestBuffer != nil {
		p.requestBuffer.Reset()
		p.requestBuffer = nil
	}
	return p.closeResponse()
}

func (p *THttpClient) Read(buf []byte) (int, error) {
	if p.response == nil {
		return 0, NewTTransportException(NOT_OPEN, "Response buffer is empty, no request.")
	}
	n, err := p.response.Body.Read(buf)
	if n > 0 && (err == nil || err == io.EOF) {
		return n, nil
	}
	return n, NewTTransportExceptionFromError(err)
}

func (p *THttpClient) ReadByte() (c byte, err error) {
	return readByte(p.response.Body)
}

func (p *THttpClient) Write(buf []byte) (int, error) {
	return p.requestBuffer.Write(buf)
}

func (p *THttpClient) WriteByte(c byte) error {
	return p.requestBuffer.WriteByte(c)
}

func (p *THttpClient) WriteString(s string) (n int, err error) {
	return p.requestBuffer.WriteString(s)
}

func (p *THttpClient) Flush() error {
	p.closeResponse()
	req, err := http.NewRequest("POST", p.url.String(), p.requestBuffer)
	if err != nil {
		return NewTTransportExceptionFromError(err)
	}
	req.Header = p.header
	response, err := p.client.Do(req)
	if err != nil {
		return NewTTransportExceptionFromError(err)
	}
	if response.StatusCode != http.StatusOK {
		p.response = response
		p.closeResponse()

		return NewTTransportException(UNKNOWN_TRANSPORT_EXCEPTION, "HTTP Response code: "+strconv.Itoa(response.StatusCode))
	}
	p.response = response
	return nil
}

func (p *THttpClient) RemainingBytes() (num_bytes uint64) {
	len := p.response.ContentLength
	if len >= 0 {
		return uint64(len)
	}

	const maxSize = ^uint64(0)
	return maxSize
}

func NewTHttpPostClientTransportFactory(url string) *THttpClientTransportFactory {
	return NewTHttpClientTransportFactoryWithOptions(url, THttpClientOptions{})
}

func NewTHttpPostClientTransportFactoryWithOptions(url string, options THttpClientOptions) *THttpClientTransportFactory {
	return NewTHttpClientTransportFactoryWithOptions(url, options)
}

func NewTHttpPostClientWithOptions(urlstr string, options THttpClientOptions) (TTransport, error) {
	return NewTHttpClientWithOptions(urlstr, options)
}

func NewTHttpPostClient(urlstr string) (TTransport, error) {
	return NewTHttpClientWithOptions(urlstr, THttpClientOptions{})
}