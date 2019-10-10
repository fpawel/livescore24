package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)
func parse24Score(d *goquery.Document) error {
	//str,_ := d.Html()
	//ioutil.WriteFile("x.html", []byte(str), 0644)
	d.Find(".daymatches tbody").Each(func(_ int, selection *goquery.Selection) {
		selection.Find("tr th.champheader a").Each(func(i int, selection *goquery.Selection) {
			log.Info( selection.Text()  )
		})
		selection.Find("tr.odd,tr.even").Each(func(i int, selection *goquery.Selection) {
			var time string
			selection.Find("td.time").Each(func(i int, selection *goquery.Selection) {
				time = strings.TrimSpace(selection.Nodes[0].FirstChild.Data)
			})
			//time := strings.TrimSpace(selection.Find("td.time").Text())
			var team [2]string
			selection.Find("td.team a").Each(func(i int, selection *goquery.Selection) {
				team[i] = selTxt( selection )
			})

			score := selTxt( selection.Find("td.left span a") )
			selection = selection.Find("td.left span").Contents()
			score2 := ""
			if len(selection.Nodes) > 2{
				score2 = strings.TrimSpace(selection.Nodes[2].Data)
				score2 = regexp.MustCompile(`\t+`).ReplaceAllString(score2, " ")
			}
			log.Info(fmt.Sprintf("%s : %q - %q %q %q", time, team[0], team[1], score, score2))

		})


	})
	return nil
}

func selTxt(selection *goquery.Selection) string{
	return strings.TrimSpace( selection.Text() )
}

func fetch( url string, parse func(*goquery.Document) error ) error {
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
	if err := parse(doc); err != nil {
		return log.Err(err)
	}
	return nil
}
