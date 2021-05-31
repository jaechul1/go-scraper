package main

import (
	"os"
	"strings"

	"github.com/jaechul1/go-scraper/scraper"
	"github.com/labstack/echo"
)

func hello(c echo.Context) error {
	return c.File("home.html")
  }

func handleScrape(c echo.Context) error {
	term := strings.ToLower(strings.TrimSpace(c.FormValue("term")))
	scraper.Scrape(term)
	defer os.Remove(term + ".csv")
	return c.Attachment(term + ".csv", term + ".csv")
}

func main() {
	e := echo.New()
	e.GET("/", hello)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}