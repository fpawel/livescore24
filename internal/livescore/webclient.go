package livescore

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/powerman/structlog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)
func parseChamps(d *goquery.Document, champs *[]Champ) error {
	d.Find(".daymatches tbody").Each(func(_ int, selection *goquery.Selection) {
		champ := Champ{Name: selTxt(selection.Find("tr th.champheader a"))}
		selection.Find("tr.odd,tr.even").Each(func(i int, selection *goquery.Selection) {
			var game Game
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				switch i {
				case 0:
					game.Time = selTxt(selection.Contents().First())
				case 1:
					game.Team1 = selTxt(selection.Find("a").Contents().First())
				case 2:
					game.Team2 = selTxt(selection.Find("a").Contents().First())
				case 3:
					game.Score = selTxt(selection.Find("span a").Contents().First())
					selection = selection.Find("span").Contents()
					if len(selection.Nodes) > 2{
						game.Score2 = reTabs.ReplaceAllString(strings.TrimSpace(selection.Nodes[2].Data), " ")
					}
				case 4:
					game.Win,_ = strconv.ParseFloat(selTxt(selection.Find("span a span")), 64)
				case 5:
					game.Draw,_ = strconv.ParseFloat(selTxt(selection.Find("span a span")), 64)
				case 6:
					game.Lose,_ = strconv.ParseFloat(selTxt(selection.Find("span a span")), 64)
				}
			})
			champ.Games = append(champ.Games, game)
		})
		*champs = append(*champs, champ)
	})
	return nil
}

var reTabs = regexp.MustCompile(`\t+`)

func selTxt(selection *goquery.Selection) string{
	return strings.TrimSpace( selection.Text() )
}

func FetchChamps( url string, champs *[]Champ) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer log.ErrIfFail(res.Body.Close)
	if res.StatusCode != 200 {
		return log.Err("status code error", "status_code",res.StatusCode, "status", res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return log.Err(err)
	}
	//{
	//	s,_ := doc.Html()
	//	log.Debug(fmt.Sprintf("%d bytes: %s", len(s), url))
	//}
	if err := parseChamps(doc, champs); err != nil {
		return log.Err(err)
	}
	return nil
}

var log = structlog.New()