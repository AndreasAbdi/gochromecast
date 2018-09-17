package dashcast

//LoadCommand is the command o send to play a url on the dashcast.
type LoadCommand struct {
	URL        string `json:"url"`
	Force      bool   `json:"force"`
	Reload     bool   `json:"reload"`
	ReloadTime int64  `json:"reload_time"`
}
