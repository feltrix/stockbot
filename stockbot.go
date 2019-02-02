package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Analise struct {
	StockName       string
	ShortTime       string
	MediumTime      string
	LongTime        string
	BuyPrices       []string
	SellPrices      []string
	BuyDescription  string
	SellDescription string
}

func main() {

	fmt.Println("Hello, master.")

	// print
	analise := Analise{
		StockName: "CIEL3",
	}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: app.tororadar.com.br
		colly.AllowedDomains("app.tororadar.com.br"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnHTML("div.graph-panel", func(e *colly.HTMLElement) {

		analizeTimeLabel := e.ChildText("div:nth-child(2) h5:first-child")

		tendency := parseGraphTendency(e)

		switch analizeTimeLabel {
		case "Curto prazo":
			analise.ShortTime = tendency
		case "Médio prazo":
			analise.MediumTime = tendency
		case "Longo prazo":
			analise.LongTime = tendency
		}

	})

	c.OnHTML("div.gray-box", func(e *colly.HTMLElement) {
		log.Println("h4: " + e.ChildText("h4"))

		isBuyDescription := strings.Contains(e.ChildText("h4"), "Avaliar compras")
		isSellDescritpion := strings.Contains(e.ChildText("h4"), "Avaliar vendas")
		var description string
		if isBuyDescription || isSellDescritpion {
			description = e.ChildText("p")
			prices := parsePrices(description)
			if isBuyDescription {
				analise.BuyPrices = prices
				analise.BuyDescription = description

			} else {
				analise.SellPrices = prices
				analise.SellDescription = description
			}
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	URL := fmt.Sprintf("https://app.tororadar.com.br/analise/%s/", analise.StockName)

	//The method visit executes a HEAD http resquest method and its returns a http code 405 for this site
	//c.Visit(fmt.Sprintf("https://app.tororadar.com.br/analise/%s/",analise.StockName))
	c.Request("GET", URL, nil, nil, nil)

	fmt.Println(analise)

}

func parseGraphTendency(e *colly.HTMLElement) string {

	var value string
	switch {
	case containsClass(e, "graph-up"):
		value = "up"
	case containsClass(e, "graph-down"):
		value = "down"
	case containsClass(e, "graph-normal"):
		value = "stay"
	}

	return value
}

func containsClass(e *colly.HTMLElement, className string) bool {
	divClass := e.Attr("class")
	return strings.Contains(divClass, className)
}

func parsePrices(descriptionText string) []string {
	re := regexp.MustCompile("\\d*\\,\\d{1,2}")
	prices := re.FindAllString(descriptionText, -1)

	return prices
}
