package gowl

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// ResponseInterface
type ResponseInterface interface {
	StatusCode() int
	Header() http.Header
	Write(w io.Writer) error
}

// Response
type Response struct {
	statusCode int
	header     http.Header
	Content    interface{}
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) SetStatusCode(code int) {
	r.statusCode = code
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) Write(w io.Writer) (err error) {
	switch content := r.Content.(type) {
	case nil:
		// no content
	case []byte:
		_, err = w.Write(content)
	case string:
		_, err = w.Write([]byte(content))
	case io.Reader:
		_, err = io.Copy(w, content)
	case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		_, err = fmt.Fprintf(w, "%d", content)
	case float32, float64:
		_, err = fmt.Fprintf(w, "%f", content)
	case func(io.Writer) error:
		err = content(w)
	default:
		_, err = fmt.Fprint(w, content)
	}
	return err
}

func NewResponse(statusCode int, content interface{}) *Response {
	return &Response{
		statusCode: statusCode,
		header:     make(http.Header),
		Content:    content,
	}
}

// ResponseWriterInterface
type ResponseWriterInterface interface {
	WriteResponse(w http.ResponseWriter) error
}

// emptyResponse
type emptyResponse struct{}

func (r *emptyResponse) StatusCode() int {
	return -1
}

func (r *emptyResponse) Header() http.Header {
	return nil
}

func (r *emptyResponse) Write(w io.Writer) error {
	return r.WriteResponse(w.(http.ResponseWriter))
}

func (r *emptyResponse) WriteResponse(w http.ResponseWriter) error {
	rw, ok := w.(http.Hijacker)
	if !ok {
		return errors.New("gowl: cannot hijack response writer")
	}
	conn, buf, err := rw.Hijack()
	if err != nil {
		return err
	}
	buf.Flush()
	conn.Close()
	return nil
}

func NewEmptyResponse() ResponseInterface {
	return &emptyResponse{}
}

// redirectResponse
type redirectResponse struct {
	request    *Request
	statusCode int
	header     http.Header
	url        string
}

func (r *redirectResponse) StatusCode() int {
	return r.statusCode
}

func (r *redirectResponse) Header() http.Header {
	return r.header
}

func (r *redirectResponse) Write(w io.Writer) error {
	return r.WriteResponse(w.(http.ResponseWriter))
}

func (r *redirectResponse) WriteResponse(w http.ResponseWriter) error {
	http.Redirect(w, r.request.Request, r.url, r.statusCode)
	return nil
}

func NewRedirectResponse(request *Request, statusCode int, url string) ResponseInterface {
	return &redirectResponse{
		request:    request,
		statusCode: statusCode,
		header:     make(http.Header),
		url:        url,
	}
}
