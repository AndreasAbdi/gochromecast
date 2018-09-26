package generic

import "github.com/imroc/req"

//RequestComponents is an abstraction over the components to be send for a post request.
type RequestComponents struct {
	URL    string
	Body   interface{}
	Header req.Header
	Params req.Param
}

//Post the request components.
func (reqComponents *RequestComponents) Post() (*req.Resp, error) {
	return req.Post(
		reqComponents.URL,
		reqComponents.Header,
		reqComponents.Params,
		reqComponents.Body,
	)
}
