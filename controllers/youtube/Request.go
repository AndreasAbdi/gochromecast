package youtube

const youtubeBaseURL = "https://www.youtube.com/"
const bindURL = youtubeBaseURL + "api/lounge/bc/bind"
const loungeTokenURL = youtubeBaseURL + "api/lounge/pairing/get_lounge_token_batch"

//GetLoungeIDRequest is a Request body for a request for a lounge id to attach to.
type GetLoungeIDRequest struct {
	ScreenIDs string `json:"screen_ids"`
}
