package gowl

import (
	"net/http"

	"github.com/lokhman/gowl/httputil"
	"github.com/lokhman/gowl/types"
)

func NegotiateRequestListener(offers []string, useFirstOffer bool) func(event EventInterface) {
	return func(event EventInterface) {
		if len(offers) == 0 {
			return
		}

		ev := event.(*RequestEvent)
		request := ev.Request()

		offer := httputil.NegotiateAcceptHeader(request.Header, "Accept", offers)
		if offer == "" && !useFirstOffer {
			link := make(httputil.HeaderValues, len(offers))
			for i, offer := range offers {
				link[i] = httputil.HeaderValue{
					Value:  "<" + request.URL.String() + ">",
					Params: types.StringMap{"type": offer},
				}
			}

			response := ErrorResponse(http.StatusNotAcceptable, "")
			response.Header().Add("Link", link.String())
			ev.SetResponse(response)
			return
		} else if useFirstOffer {
			offer = offers[0]
		}

		request.params.Set(":accept", offer)
	}
}
