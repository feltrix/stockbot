package main

import (
	"strings"
	"fmt"
	"github.com/gocolly/colly"
)

type Analise struct {
	StockName  string
	ShortTime  string
	MediumTime string
	LongTime   string 
}

func main() {

	fmt.Println("Hello, master.")
	
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
		
		analizeLabel := e.ChildText("div:nth-child(2) h5:first-child")
		divClass := e.Attr("class")
		var value string 
		switch {
			case strings.Contains(divClass, "graph-up"):
				value = "up"
			case strings.Contains(divClass, "graph-down"):
				value = "down"
			case strings.Contains(divClass, "graph-normal"):
				value = "stay"
		}

		switch analizeLabel {
			case "Curto prazo":
				analise.ShortTime = value
			case "Médio prazo":
				analise.MediumTime = value
			case "Longo prazo":	
				analise.LongTime = value
		}
		
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	URL := fmt.Sprintf("https://app.tororadar.com.br/analise/%s/",analise.StockName)

	//The method visit executes a HEAD http resquest method and its returns a http code 405 for this site
	//c.Visit(fmt.Sprintf("https://app.tororadar.com.br/analise/%s/",analise.StockName))
	c.Request("GET", URL, nil, nil, nil)

	fmt.Println(analise)
	
}