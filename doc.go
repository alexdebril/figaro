/*
Package figaro is a library designed to speed-up Go applications writing.

It features HTTP response building through a simple set of helpers. Intended use-case is the simple generation of JSON Payload:

		type Message struct {
			Date    time.Time `json:"date"`
			Message string    `json:"message"`
		}

		func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			message := &Message{
				Date:    time.Now(),
				Message: "this is the new message",
			}
			resp := response.NewResponse(response.WithJsonPayload(message))
			resp.Write(w)
		}

*/

package figaro
