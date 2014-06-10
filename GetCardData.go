package main

import "fmt"
import "net/http"
import "io"
import "io/ioutil"
import "encoding/json"
import "errors"
import "os"

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
	// path = /set/side/faction/
	path := fmt.Sprintf("/dev/Go/GoRunner/cards/%s/%s/%s/", card.Set_Code, card.Side_Code,
		card.Faction_Code)
	pathExists, err := exists(path)
	if (!pathExists) {
		if err != nil {
			return err
		}
		os.MkdirAll(path, 0700)
	}
	path = fmt.Sprintf("/dev/Go/GoRunner/cards/%s/%s/%s/%d.card", card.Set_Code, card.Side_Code,
		card.Faction_Code, card.Number)
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	
	output, err := json.Marshal(card)
	check(err)
	
	_, err = file.Write(output)
	check(err)
	
	path = fmt.Sprintf("/dev/Go/GoRunner/cards/%s/%s/%s/images", card.Set_Code, card.Side_Code,
		card.Faction_Code)
	pathExists, err = exists(path)
	if (!pathExists) {
		if err != nil {
			return err
		}
		os.MkdirAll(path, 0700)
	}
	path = fmt.Sprintf("/dev/Go/GoRunner/cards/%s/%s/%s/images/%d.png", card.Set_Code, card.Side_Code,
		card.Faction_Code, card.Number)
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

func main() {	
	currentCard := 1001
	for ;currentCard < 1012; {
		cardString := fmt.Sprintf("0%d", currentCard)
		card, err := getCardData(cardString)
		if err != nil {
			fmt.Print(err)
			fmt.Printf("\nCurrent card: 0%d\n", currentCard)
			break;
		}
		fmt.Printf("%s successfully retrieved...\n", card.Title)
		err = createCardFile(card)
		check(err)
		currentCard++
	}
}
