package response

type ServerStatusCode int

const (
	StatusOK          ServerStatusCode = 200
	StatusBadRequest  ServerStatusCode = 400
	StatusServerError ServerStatusCode = 500
)

var statusText = map[ServerStatusCode]string{
	StatusOK:          "Ok",
	StatusBadRequest:  "Bad Request",
	StatusServerError: "Internal Server Error",
}

func (rs responseState) String() string {
	switch rs {
	case respWritingStatusLine:
		return "respWritingStatusLine"
	case respWritingHeaders:
		return "respWritingHeaders"
	case respWritingBody:
		return "respWritingBody"
	case respWritingTrailer:
		return "respWritingTrailer"
	default:
		return "unknown responseState"
	}
}
