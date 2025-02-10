package models

type Link struct {
	ItemID       string
	CurrentLink  string
	PreviousLink string
}

type OutputNotification struct {
	ItemID string
}

type MathomOffer struct {
	ItemName      string
	CurrentPrice  float64
	OriginalPrice float64
	Language      string
	Link          string
}

type BGGData struct {
	BGGName   string
	Rating    float64
	Expansion bool
	Own       int
	Whishlist int
	Sell      int
	Want      int
}
