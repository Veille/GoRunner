package main

import "fmt"
import "net/http"
import "io"
import "io/ioutil"
import "encoding/json"
import "errors"
import "os"
import "strconv"

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

const CARDPATH string = "/dev/Go/src/github.com/veille/GoRunner/cards"

// Used to generate a request to a URL and return the content supplied
func getContent(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		//fmt.Println("Problem with GET")
		err := errors.New("!! Problem with GET request")
		return nil, err
	}
	// Defer the closing of the body
	defer resp.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("\nResponse has been extracted...\n")
	// We're done. Return the bytes
	return body, nil
}

// Used to extract card data from a supplied URL
func getCardData(cardId string) (*RunnerCard, error) {
	// Fetch the JSON content for the given card
	url := fmt.Sprintf("http://netrunnerdb.com/api/card/%s", cardId)
	content, err := getContent(url)
	// static URL for testing
	//content, err := getContent("http://netrunnerdb.com/api/cards")
	if err != nil {
		fmt.Printf("Problem with URL: %s\n", url)
		return nil, err
	}
	//fmt.Println("Successful getting content...")
	// Slice the jason data out of the array
	content = content[1: len(content) - 1]	
	// Fill the card with the JSON data
	newCards := &RunnerCard{}
	err = json.Unmarshal(content, &newCards)
	if err != nil {
		fmt.Println("\n--------------------------------------")
		fmt.Printf("An error has occured while unmarshalling the JSON data: \n%s\n", err)
		fmt.Println(string(content))
		fmt.Println("\n--------------------------------------\n")
		return nil, err
	}
	//fmt.Println("JSON data has been unmarshalled...\n")
	return newCards, err
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func createCardFile(card *RunnerCard) error {
	cardCode := card.Code[(len(card.Code) - 3): (len(card.Code))]

	// path = /set/side/faction/
	path := fmt.Sprintf("%s/%s/", CARDPATH, card.Set_Code)
	pathExists, err := exists(path)
	if (!pathExists) {
		if err != nil {
			return err
		}
		os.MkdirAll(path, 0700)
	}
	path = fmt.Sprintf("%s/%s/%s.card", CARDPATH, card.Set_Code, cardCode)
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	
	output, err := json.Marshal(card)
	check(err)
	
	_, err = file.Write(output)
	check(err)
	
	path = fmt.Sprintf("%s/%s/images", CARDPATH, card.Set_Code)
	pathExists, err = exists(path)
	if (!pathExists) {
		if err != nil {
			return err
		}
		os.MkdirAll(path, 0700)
	}
	path = fmt.Sprintf("%s/%s/images/%s.png", CARDPATH, card.Set_Code, cardCode)
	pathExists, err = exists(path)
	if (!pathExists) {
		fmt.Println("New image created\n")
		out, err := os.Create(path)
		check(err)
		defer out.Close()
		resp, err := http.Get(fmt.Sprintf("http://netrunnerdb.com%s", card.LargeImageSrc))
		check(err)
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		check(err)
	} else {
		fmt.Println("Image already created")
	}
	return nil
}

func getSetData(startCode, endCode int) {
	currentCard := startCode
	//fmt.Printf("Start: %d  End: %d\n", startCode, endCode)
	for ;currentCard <= endCode; {
		var cardString string
		if currentCard < 10000 {
			cardString = fmt.Sprintf("0%d", currentCard)
		} else {
			cardString = fmt.Sprintf("%d", currentCard)
		}
		if (!alreadyRetrieved(cardString)) {
			card, err := getCardData(cardString)
			if err != nil {
				fmt.Print(err)
				fmt.Printf("\nCurrent card: 0%d\n", currentCard)
				break;
			}
			fmt.Printf("%s successfully retrieved...\n", card.Title)
			err = createCardFile(card)
			check(err)
		} else {
			fmt.Printf("Card %s has already been retrieved\n", cardString)
		}
		currentCard++
	}
}

func getSetName(setCode string) (string, error) {
	intSetCode, err := strconv.Atoi(setCode)
	check(err)
	switch intSetCode {
		case core:
			return "core", nil
		case wla:
			return "wla", nil
		case ta:
			return "ta", nil
		case ca:
			return "ca", nil
		case asis:
			return "asis", nil
		case hs:
			return "hs", nil
		case fp:
			return "fp", nil
		case cac:
			return "cac", nil
		case om:
			return "om", nil
		case st:
			return "st", nil
		case mt:
			return "mt", nil
		case tc:
			return "tc", nil
		case fal:
			return "fal", nil
		case dt:
			return "dt", nil
		case hap:
			return "hap", nil
		case up:
			return "up", nil
		case tsb:
			return "tsb", nil
		case fc:
			return "fc", nil
		case uao:
			return "uao", nil
		case atr:
			return "atr", nil
	}
	err = errors.New("!! Set name not found")
	return "", err
}

func alreadyRetrieved(code string) bool {
	setCode := code[0: len(code) - 3]
	cardCode := code[len(code) - 3: len(code)]
	setName, err := getSetName(setCode)
	check(err)
	
	path := fmt.Sprintf("%s/%s/%s.card", CARDPATH, setName, cardCode)
	pathExists, err := exists(path)
	check(err)
	
	return pathExists
}

func main() {	
	// Core Set 				-- 01001 -> 01113
	// Genesis Cycle
	//	|- What Lies Ahead 		-- 02001 -> 02020
	//	|- Trace Amounts		-- 02021 -> 02040
	//	|- Cyber Exodus			-- 02041 -> 02060
	//	|- A Study In Static	-- 02061 -> 02080
	//	|- Humanity's Shadow	-- 02081 -> 02100
	//	|- Future Proof			-- 02101 -> 02120
	//	|----------------------------------------
	// Creation and Control		-- 03001 -> 03055
	// Spin Cycle
	//	|- Opening Moves		-- 04001 -> 04020
	//	|- Second Thoughts		-- 04021 -> 04040
	//	|- Mala Tempora			-- 04041 -> 04060
	//	|- True Colors			-- 04061 -> 04080
	//	|- Fear and Loathing	-- 04081 -> 04100
	//	|- Double Time			-- 04101 -> 04120
	//	|----------------------------------------
	// Honor and Profit			-- 05001 -> 05055
	// Lunar Cycle
	//	|- Upstalk				-- 06001 -> 06020
	//	|- The Spaces Between	-- 06021 -> 06040
	//	|- First Contact		-- 06041 -> 06060
	//	|- Up and Over			-- 06061 -> 06080
	//	|- All That Remains		-- 06081 -> 06100
	//	|-						-- 06101 -> 06120
	
	//getSetData(1001, 1003)		// TEST
	getSetData(1001, 1113)	// Core
	//getSetData(2001, 2120)	// Genesis
	//getSetData(3001, 3055)	// Creation & Control
	//getSetData(4001, 4120)	// Spin
	//getSetData(5001, 5055)	// Honor & Profit
}
