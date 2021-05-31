package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	ccsv "github.com/tsak/concurrent-csv-writer"
)

type jobInfo struct {
	link string
	title string
	company string
	location string
	salary string
	summary string
}

// Scrape kr.indeed.com by a term
func Scrape(term string) {
	baseURL := "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []jobInfo
	scrapeC := make(chan []jobInfo)
	numberPages := 5
	
	for i := 0; i < numberPages; i++ {
		go getPage(baseURL, i, scrapeC)
	}
	
	for i := 0; i < numberPages; i++ {
		pageJobs := <- scrapeC
		jobs = append(jobs, pageJobs...)
	}
	writeJobs(term, jobs) 
	fmt.Println("Extracted done")
}

func getPage(baseURL string, page int, scrapeC chan<- []jobInfo) {
	fmt.Println("Getting the jobs in page:", page)
	var jobs []jobInfo
	c := make(chan jobInfo)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobCards := doc.Find(".jobsearch-SerpJobCard")
	jobCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < jobCards.Length(); i++ {
		job := <- c
		jobs = append(jobs, job)
	}
	scrapeC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- jobInfo) {
	id, _ := card.Attr("data-jk")
	title := strings.TrimSpace(card.Find(".title>a").Text())
	company := strings.TrimSpace(card.Find(".company").Text())
	location := strings.TrimSpace(card.Find(".location").Text())
	salary := strings.TrimSpace(card.Find(".salary").Text())
	summary := strings.TrimSpace(card.Find(".summary").Text())
	c <- jobInfo{
		link: "https://kr.indeed.com/viewjob?jk=" + id, 
		title: title, 
		company: company, 
		location: location, 
		salary: salary, 
		summary: summary}
}

func writeJobs(term string, jobs []jobInfo) {
	fmt.Println("Writing jobs in csv")
	w, err := ccsv.NewCsvWriter(term + ".csv")
    checkErr(err)
    defer w.Close()

	header := []string{"Link", "Title", "Company", "Location", "Salary", "Summary"}
	headerWriteErr := w.Write(header)
	checkErr(headerWriteErr)

	done := make(chan bool)
	for _, job := range jobs {
		go func(job jobInfo) {
			jobSlice := []string{job.link, job.title, job.company, job.location, job.salary, job.summary}		
			jobWriteErr := w.Write(jobSlice)
			checkErr(jobWriteErr)
			done <- true
		}(job)
	}

	for i := 0; i < len(jobs); i++ {
		<- done
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}