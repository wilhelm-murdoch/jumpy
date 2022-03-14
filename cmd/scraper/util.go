package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"os"
// 	"regexp"
// 	"strconv"
// 	"strings"

// 	"github.com/wilhelm-murdoch/jumpscare-api/feed"
// 	"github.com/wilhelm-murdoch/jumpscare-api/logger"
// 	"github.com/wilhelm-murdoch/jumpscare-api/movie"

// 	"github.com/PuerkitoBio/goquery"
// )

// func GetDocumentFromUrl(url string) (*goquery.Document, error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		logger.Err.Fatal(err)
// 	}

// 	defer response.Body.Close()

// 	if response.StatusCode != 200 {
// 		logger.Err.Fatalf("status %d returned for: %s", response.StatusCode, url)
// 	}

// 	return goquery.NewDocumentFromReader(response.Body)
// }

// func WriteJsonToFile(path string, object interface{}) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}

// 	encoder := json.NewEncoder(file)

// 	err = encoder.Encode(object)
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	return nil
// }

// func GetMovieListFromUrl(url string) (*feed.Feed, error) {
// 	feed := feed.NewFeed()

// 	document, err := GetDocumentFromUrl(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	document.Find("div.entry-content > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
// 		node := s.Find("td:nth-child(1) a")

// 		url, ok := node.Attr("href")
// 		if ok {
// 			title := node.Text()
// 			release, _ := strconv.Atoi(s.Find("td:nth-child(3)").Text())
// 			movie := movie.NewMovie(title, release, url)

// 			feed.Add(movie)
// 		}
// 	})

// 	return feed, nil
// }

// func GetMovieDetailsFromUrl(movie movie.Movie) error {
// 	document, err := GetDocumentFromUrl(movie.SourceUrl)
// 	if err != nil {
// 		return err
// 	}

// 	document.Find("div.inside-article").Each(func(i int, s *goquery.Selection) {
// 		s.Find("div.entry-content > h3 ~ p").Each(func(i int, s *goquery.Selection) {
// 			pattern := regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}) – (.*).$`)
// 			matches := pattern.FindStringSubmatch(s.Text())

// 			if len(matches) == 3 {
// 				var major bool = false
// 				if s.Has("strong").Length() >= 1 {
// 					major = true
// 				}

// 				movie.AddJumpScare(matches[1], matches[2], major)
// 			}
// 		})

// 		s.Find("div.video-info-grid-column > p").Each(func(i int, s *goquery.Selection) {
// 			switch strings.Fields(s.Text())[0] {
// 			case "Reviews:":
// 				s.Find("a").Each(func(i int, s *goquery.Selection) {
// 					href, ok := s.Attr("href")
// 					if ok {
// 						movie.AddReview(s.Text(), href)
// 					}
// 				})
// 			case "Director:":
// 				for _, d := range strings.Split(strings.Replace(s.Text(), "Director:", "", 1), ",") {
// 					movie.AddDirector(d)
// 				}
// 			case "Rating:":
// 				movie.SetContentRating(s.Text())
// 			case "Runtime:":
// 				movie.SetRuntimeFromPattern(regexp.MustCompile(`^Runtime:\s+(\d+)\sminutes$`), s.Text())
// 			case "Tags:":
// 				s.Find("a").Each(func(i int, s *goquery.Selection) {
// 					movie.AddTag(s.Text())
// 				})
// 			}
// 		})
// 	})

// 	return nil
// }
