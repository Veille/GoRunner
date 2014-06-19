package cardData

import "encoding/json"

type RunnerCard struct {
	LastModified 	string 	`json:"last-modified"`
	Code 			string	`json:"code"`
	Title			string	`json:"title"`
	Type			string	`json:"type"`
	Type_Code		string	`json:"type_code"`
	Subtype			string	`json:"subtype"`
	Subtype_Code	string	`json:"subtype_code"`
	Text			string	`json:"text"`
	BaseLink		int		`json:"baselink"`
	Faction			string	`json:"faction"`
	Faction_Code	string	`json:"faction_code"`
	Faction_Letter	string	`json:"faction_letter"`
	Flavor			string	`json:"flavor"`
	Illustrator		string	`json:"illustrator"`
	InfluenceLimit	int		`json:"influencelimit"`
	MinimumDeckSize	int		`json:"minimumdecksize"`
	Number			int		`json:"number"`
	Quantity		int		`json:"quantity"`
	SetName			string	`json:"setname"`
	Set_Code		string	`json:"set_code"`
	Side			string	`json:"side"`
	Side_Code		string	`json:"side_code"`
	Uniqueness		bool	`json:"uniqueness"`
	CycleNumber		int		`json:"cyclenumber"`
	URL				string	`json:"url"`
	ImageSrc		string	`json:"imagesrc"`
	LargeImageSrc	string	`json:"largeimagesrc"`
}

const (
	_ = iota
	core
	wla	
	ta
	ca
	asis
	hs
	fp
	cac
	om
	st
	mt
	tc
	fal
	dt
	hap
	up
	tsb
	fc
	uao
	atr
)
