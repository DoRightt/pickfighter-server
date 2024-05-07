package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"fightbettr.com/fb-server/pkg/model"
	"fightbettr.com/scraper/internal/scraperutil"
	data "fightbettr.com/scraper/pkg"
	"fightbettr.com/scraper/pkg/logger"
	"gopkg.in/yaml.v3"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

var config Config
var gc *colly.Collector
var detailsCollector *colly.Collector
var collection = model.FightersCollection{}
var wg sync.WaitGroup
var l *zap.SugaredLogger

// main function responsible for initializing the web scraping process.
// It sets up the logger, creates collector instances, defines URL and limits for the main collector,
// and specifies callback functions for HTML elements. It initiates the web scraping process by visiting
// the initial URL and waits for the wait group to finish before printing "DONE" to the console and saving
// the collected data to a JSON file using the saveToJSON function.
func main() {
	var useProxy bool
	var startPage int
	var toAdd bool

	flag.BoolVar(&toAdd, "add", false, "Add fighters")
	flag.BoolVar(&useProxy, "proxy", false, "Use proxy")
	flag.IntVar(&startPage, "start", 0, "Start page")
	flag.Parse()

	logFlag := scraperutil.GetLoggerFlag(toAdd)
	if err := logger.Initialize(logFlag); err != nil {
		fmt.Println("Error while initializing logger: ", err)
		return
	}

	l = logger.Get()
	gc = colly.NewCollector()
	detailsCollector = gc.Clone()

	url := "https://www.ufc.com/athletes/all"

	if startPage > 0 {
		url = fmt.Sprintf("%s?page=%d", url, startPage)
	}

	if useProxy {
		f, err := os.Open("./configs/proxy.yaml")
		if err != nil {
			l.Fatal("Failed to open configuration", zap.Error(err))
		}

		if err := yaml.NewDecoder(f).Decode(&config); err != nil {
			l.Fatal("Failed to parse configuration", zap.Error(err))
		}
	}

	gc.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 3 * time.Second,
	})

	gc.OnRequest(func(r *colly.Request) {
		if useProxy {
			proxy := getProxy()
			proxyUrl := fmt.Sprintf("socks5h://%s:%s@%s", config.Login, config.Password, proxy)

			gc.SetProxy(proxyUrl)
			l.Infow(proxy, "type", "proxy address")
		}

		r.Headers.Set("User-Agent", "Mozilla/5.0")
	})

	gc.OnHTML("div[class*='flipcard__action'] a[href]", parseAthletesListing)
	gc.OnHTML("li.pager__item a[href]", moveNextPage)
	detailsCollector.OnHTML("div[class='hero-profile-wrap']", getData)

	err := gc.Visit(url)
	if err != nil {
		log.Fatalf("Error while request: %v", err)
	}

	wg.Wait()

	fmt.Println("DONE")
	l.Infow("DONE", "type", "result")

	saveToJSON(collection, toAdd)
}

// parseAthletesListing is a callback function used with colly that extracts athlete URLs from a given colly.HTMLElement 'e'.
// It increments the wait group and defers its decrement for synchronization. The function extracts the athlete's URL,
// converts it to an absolute URL, prints it to the console, and then uses a detailsCollector to visit the athlete's detailed page.
// This function is typically used during web scraping to collect athlete URLs for subsequent detailed data extraction.
func parseAthletesListing(e *colly.HTMLElement) {
	wg.Add(1)
	defer wg.Done()

	athleteURL := e.Attr("href")
	athleteURL = e.Request.AbsoluteURL(athleteURL)

	fmt.Println("Athlete link:", athleteURL)
	l.Infow(athleteURL, "type", "athlete link")

	detailsCollector.Visit(athleteURL)

	// Multithread scrapping
	// go func() {
	// 	defer wg.Done()
	// 	athleteURL := e.Attr("href")
	// 	athleteURL = e.Request.AbsoluteURL(athleteURL)

	// 	fmt.Println("Athlete link:", athleteURL)

	// 	detailsCollector.Visit(athleteURL)
	// }()
}

// getData is a callback function used with colly that extracts fighter data from a given colly.HTMLElement 'e'.
// It initializes a Fighter model, sets basic information such as name, nickname, URL, and image, and then calls
// SetDivision and SetStatistic functions to update division and general statistics. Finally, it calls parseData
// to extract additional details about the fighter and appends the resulting Fighter instance to a FightersCollection.
// This function is typically used during web scraping to gather comprehensive information about a fighter.
func getData(e *colly.HTMLElement) {
	wg.Add(1)
	defer wg.Done()

	fighterEl := e.DOM.Parent()

	profileEl := fighterEl.Find("div.hero-profile-wrap")
	statString := profileEl.Find("p.hero-profile__division-body").Text()

	fighter := model.Fighter{
		Name:       profileEl.Find("h1.hero-profile__name").Text(),
		NickName:   profileEl.Find("p.hero-profile__nickname").Text(),
		FighterUrl: e.Request.URL.String(),
		ImageUrl:   profileEl.Find(".hero-profile__image-wrap img").AttrOr("src", ""),
	}

	data.SetDivision(&fighter, profileEl.Find("p.hero-profile__division-title").Text())
	data.SetStatistic(&fighter, statString)

	parseData(&fighter, fighterEl)

	collection.Fighters = append(collection.Fighters, fighter)
}

// parseData a unifying function for parsing data from different blocks of information
func parseData(f *model.Fighter, fighterEl *goquery.Selection) {
	parseBioFields(f, fighterEl)
	parseMainStats(f, fighterEl)
	parseSpecialStats(f, fighterEl)
	parseWinMethodStats(f, fighterEl)
}

// parseBioFields parses fighter's data from biography block and sets values to model.Fighter
func parseBioFields(f *model.Fighter, fighterEl *goquery.Selection) {
	fields := fighterEl.Find("div.c-bio__info-details")
	fields.Find("div.c-bio__info-details .c-bio__field").Each(func(index int, bioField *goquery.Selection) {
		fieldLabel := bioField.Find(".c-bio__label").Text()
		fieldValue := strings.TrimSpace(bioField.Find(".c-bio__text").Text())

		switch fieldLabel {
		case "Age":
			v, err := strconv.Atoi(fieldValue)
			if err != nil {
				l.Errorf("Age conversion error: %s", err)
			} else {
				f.Age = int8(v)
			}
		case "Status":
			f.Status = fieldValue
		case "Hometown":
			f.Hometown = fieldValue
		case "Trains at":
			f.TrainsAt = fieldValue
		case "Fighting style":
			f.FightingStyle = fieldValue
		case "Height":
			v, err := strconv.ParseFloat(fieldValue, 32)
			if err != nil {
				l.Errorf("Height conversion error: %s", err)
			} else {
				f.Height = float32(v)
			}
		case "Weight":
			v, err := strconv.ParseFloat(fieldValue, 32)
			if err != nil {
				l.Errorf("Weight conversion error:", err)
			} else {
				f.Weight = float32(v)
			}
		case "Octagon Debut":
			f.OctagonDebut = fieldValue
			f.DebutTimestamp = scraperutil.GetDebutTimestamp(fieldValue)
		case "Reach":
			v, err := strconv.ParseFloat(fieldValue, 32)
			if err != nil {
				l.Error("Reach conversion error:", err)
			} else {
				f.Reach = float32(v)
			}
		case "Leg reach":
			v, err := strconv.ParseFloat(fieldValue, 32)
			if err != nil {
				l.Error("Leg Reach conversion error:", err)
			} else {
				f.LegReach = float32(v)
			}
		}
	})
}

// parseMainStats parses stats from main fighter block and sets values to fighter.Stats
func parseMainStats(f *model.Fighter, fighterEl *goquery.Selection) {
	reg := regexp.MustCompile("[^0-9]+")
	fields := fighterEl.Find("div.stats-records-inner-wrap")
	fields.Find("div.c-stat-compare__group").Each(func(index int, bioField *goquery.Selection) {
		fieldLabel := bioField.Find(".c-stat-compare__label").Text()
		fieldValue := strings.TrimSpace(bioField.Find(".c-stat-compare__number").Text())

		switch fieldLabel {
		case "Sig. Str. Landed":
			if fieldValue != "" {
				v, err := strconv.ParseFloat(fieldValue, 32)
				if err != nil {
					l.Error("Sig. Str. Landed conversion error:", err)
				} else {
					f.Stats.SigStrLanded = float32(v)
				}
			}
		case "Sig. Str. Absorbed":
			if fieldValue != "" {
				v, err := strconv.ParseFloat(fieldValue, 32)
				if err != nil {
					l.Error("Sig. Str. Absorbed conversion error:", err)
				} else {
					f.Stats.SigStrAbs = float32(v)
				}
			}
		case "Sig. Str. Defense":
			numericString := reg.ReplaceAllString(fieldValue, "")
			if numericString != "" {
				v, err := strconv.Atoi(numericString)
				if err != nil {
					l.Error("Sig. Str. Defense conversion error:", err)
				} else {
					f.Stats.SigStrDefense = int8(v)
				}
			}
		case "Takedown Defense":
			numericString := reg.ReplaceAllString(fieldValue, "")
			v, err := strconv.Atoi(numericString)
			if err != nil {
				if fieldValue != "" {
					l.Error("Takedown Defense conversion error:", err)
				}
			} else {
				f.Stats.TakedownDefense = int8(v)
			}
		case "Takedown avg":
			if fieldValue != "" {
				v, err := strconv.ParseFloat(fieldValue, 32)
				if err != nil {
					l.Error("Takedown avg conversion error:", err)
				} else {
					f.Stats.TakedownAvg = float32(v)
				}
			}
		case "Submission avg":
			if fieldValue != "" {
				v, err := strconv.ParseFloat(fieldValue, 32)
				if err != nil {
					l.Error("Submission avg conversion error:", err)
				} else {
					f.Stats.SubmissionAvg = float32(v)
				}
			}
		case "Knockdown Avg":
			if fieldValue != "" {
				v, err := strconv.ParseFloat(fieldValue, 32)
				if err != nil {
					l.Error("Knockdown Avg conversion error:", err)
				} else {
					f.Stats.KnockdownAvg = float32(v)
				}
			}
		case "Average fight time":
			f.Stats.AvgFightTime = fieldValue
		}
	})
}

// parseSpecialStats parses stats from special fighter block and sets values to fighter.Stats.
func parseSpecialStats(f *model.Fighter, fighterEl *goquery.Selection) {
	fields := fighterEl.Find("div.stats-records-inner-wrap")

	fields.Find("div.c-overlap__inner .c-overlap__stats").Each(func(index int, bioField *goquery.Selection) {
		fieldLabel := bioField.Find("dt.c-overlap__stats-text").Text()
		fieldValue := strings.TrimSpace(bioField.Find("dd.c-overlap__stats-value").Text())

		switch fieldLabel {
		case "Sig. Strikes Landed":
			v, err := strconv.Atoi(fieldValue)
			if err != nil {
				l.Error("Total Sig. Strikes Landed conversion error:", err)
			} else {
				f.Stats.TotalSigStrLanded = v
			}
		case "Sig. Strikes Attempted":
			v, err := strconv.Atoi(fieldValue)
			if err != nil {
				l.Error("Total Sig. Strikes Attempted conversion error:", err)
			} else {
				f.Stats.TotalSigStrAttempted = v
			}
		case "Takedowns Landed":
			v, err := strconv.Atoi(fieldValue)
			if err != nil {
				if fieldValue != "" {
					l.Error("Total Takedowns Landed conversion error:", err)
				}
			} else {
				f.Stats.TotalTkdLanded = v
			}
		case "Takedowns Attempted":
			v, err := strconv.Atoi(fieldValue)
			if err != nil {
				l.Error("Total Takedowns Attempted conversion error:", err)
			} else {
				f.Stats.TotalTkdAttempted = v
			}
		}
	})

	if f.Stats.TotalTkdAttempted != 0 {
		f.Stats.TkdAccuracy = int(float64(f.Stats.TotalTkdLanded) / float64(f.Stats.TotalTkdAttempted) * 100)
	}

	if f.Stats.TotalSigStrAttempted != 0 {
		f.Stats.StrAccuracy = int(float64(f.Stats.TotalSigStrLanded) / float64(f.Stats.TotalSigStrAttempted) * 100)
	}
}

// parseWinMethodStats parses stats from the win methods block and sets values to fighter.Stats
func parseWinMethodStats(f *model.Fighter, el *goquery.Selection) {
	fields := el.Find("div.stats-records-inner-wrap")

	fields.Find("div.stats-records:last-of-type div.stats-records-inner .c-stat-3bar__group").Each(func(index int, bioField *goquery.Selection) {
		fieldLabel := strings.TrimSpace(bioField.Find("div.c-stat-3bar__label").Text())
		fieldValue := strings.TrimSpace(bioField.Find("div.c-stat-3bar__value").Text())

		switch fieldLabel {
		case "KO/TKO":
			v, err := strconv.Atoi(strings.Split(fieldValue, " ")[0])
			if err != nil {
				l.Error("KO/TKO data conversion error:", err)
			} else {
				f.Stats.WinByKO = v
			}
		case "DEC":
			v, err := strconv.Atoi(strings.Split(fieldValue, " ")[0])
			if err != nil {
				l.Error("DEC data conversion error:", err)
			} else {
				f.Stats.WinByDec = v
			}

		case "SUB":
			v, err := strconv.Atoi(strings.Split(fieldValue, " ")[0])
			if err != nil {
				l.Error("SUB data conversion error:", err)
			} else {
				f.Stats.WinBySub = v
			}
		}
	})
}

// moveNextPage navigates to the next page using colly.
// It extracts the "href" attribute from the provided HTML element,
// converts it to an absolute URL, prints it, and visits the next page.
func moveNextPage(e *colly.HTMLElement) {
	wg.Add(1)
	defer wg.Done()

	nextUrl := e.Attr("href")
	nextUrl = e.Request.AbsoluteURL(nextUrl)

	fmt.Println("Next page:", nextUrl)
	l.Infow(nextUrl, "type", "next page")

	e.Request.Visit(nextUrl)
}

// saveToJSON marshals the provided FightersCollection to JSON format
// and writes it to the file. If there is an error during
// marshalling, file opening, or file writing, it logs the error. This function
// is used to store fighter data in JSON format.
func saveToJSON(c model.FightersCollection, toAdd bool) {
	if toAdd {
		scraperutil.AddToExistedCollection(c)
	} else {
		scraperutil.CreateNewCollection(c)
	}
}

func getProxy() string {
	if len(config.Proxys) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(config.Proxys))

	return config.Proxys[idx]
}
