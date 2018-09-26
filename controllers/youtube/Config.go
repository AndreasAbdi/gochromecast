package youtube

const youtubeBaseURL = "https://www.youtube.com/"
const bindURL = youtubeBaseURL + "api/lounge/bc/bind"
const loungeTokenURL = youtubeBaseURL + "api/lounge/pairing/get_lounge_token_batch"

var defaultHeaders = map[string]string{
	"Origin":       youtubeBaseURL,
	"Content-Type": "application/x-www-form-urlencoded"}

const loungeIDHeader = "X-YouTube-LoungeId-Token"
const requestIDKey = "RID"
const sessionIDKey = "SID"
const versionKey = "VER"
const cVersionKey = "CVER"

const bindVersion = "8"
const bindCVersion = "1"
