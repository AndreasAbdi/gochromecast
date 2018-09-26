package youtube

//ScreenTokenData is the internal data structure for a screen/loungetoken pairing.
type ScreenTokenData struct {
	ScreenID    string `json:"screenId"`
	LoungeToken string `json:"loungeToken"`
	Expiration  int64  `json:"expiration"`
}

//LoungeTokenResponse is all the screen/loungetoken pairings requested.
type LoungeTokenResponse struct {
	Screens []*ScreenTokenData `json:"screens"`
}
