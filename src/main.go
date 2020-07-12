package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"std/fmt"
	"std/strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type jobResult struct {
	id string
	title string
	location string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main () {
	var totalJobs []jobResult
	pages := getPages()

	for i := 0 ; i < pages ; i++{
		jobs := accessLink(i)
		totalJobs = append(totalJobs,jobs...) //jobs의 모든 contents
	}
	writeJobs(totalJobs) // [{id title location} {} {}]의 형태
	fmt.Println("Extracted", len(totalJobs))
}

func writeJobs(jobs []jobResult) {
	file, err := os.Create("jobs.csv") // create file using os package
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title","Location","Salary","Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _,job := range jobs{
		jobSlice := []string{"https://kr.indeed.com/%EC%B1%84%EC%9A%A9%EB%B3%B4%EA%B8%B0?jk="+job.id+" ", job.title, job.location, job.salary, job.summary}
		jErr := w.Write(jobSlice)
		checkErr(jErr)
	}
}

func accessLink (page int) []jobResult {
	var jobs []jobResult   //[] slice인데 들어가는 내용은 {id title location}
	URL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Reqeusting", URL)
	res, err := http.Get(URL)
	checkErr(err)
	checkCode(res.StatusCode)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, s *goquery.Selection){
	job := extractJob(s) // {id title location}
	jobs = append(jobs, job) // [{} {} {}]의 형태
	})
	return jobs
}

func extractJob(s *goquery.Selection) jobResult{
	id, _ := s.Attr("data-jk")
	title := cleanString(s.Find(".title>a").Text())
	location := cleanString(s.Find(".location").Text())
	salary := cleanString(s.Find(".salaryText").Text())
	summary:= cleanString(s.Find(".summary").Text())
	return jobResult{
		id: id, 
		title: title, 
		location: location, 
		salary: salary, 
		summary: summary}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str))," ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res.StatusCode)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr (err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode (code int) {
	if code != 200 {
		log.Fatalln("Request failed with status code", code)
	}
}