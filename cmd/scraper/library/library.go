package library

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/wilhelm-murdoch/jumpy/cmd/scraper/logger"
	"github.com/wilhelm-murdoch/jumpy/cmd/scraper/models"

	"github.com/PuerkitoBio/goquery"
)

func GetDocumentFromUrl(url string) (*goquery.Document, error) {
	response, err := http.Get(url)
	if err != nil {
		logger.Error(err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		logger.Error("status %d returned for: %s", response.StatusCode, url)
	}

	return goquery.NewDocumentFromReader(response.Body)
}

func GetMovieListFromUrl(url string) (*models.Feed, error) {
	feed := models.NewFeed()

	document, err := GetDocumentFromUrl(url)
	if err != nil {
		return feed, err
	}

	document.Find("div.entry-content > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		node := s.Find("td:nth-child(1) a")

		url, ok := node.Attr("href")
		if ok {
			title := node.Text()
			release, _ := strconv.Atoi(s.Find("td:nth-child(3)").Text())
			movie := models.NewMovie(title, release, url)

			feed.Add(movie)
		}
	})

	return feed, nil
}

func GetMovieDetailsFromUrl(movie models.Movie) (models.Movie, error) {
	document, err := GetDocumentFromUrl(movie.SourceUrl)
	if err != nil {
		return movie, err
	}

	document.Find("div.inside-article").Each(func(i int, s *goquery.Selection) {
		s.Find("div.entry-content > h3 ~ p").Each(func(i int, s *goquery.Selection) {
			pattern := regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}) – (.*).$`)
			matches := pattern.FindStringSubmatch(s.Text())

			if len(matches) == 3 {
				var major bool = false
				if s.Has("strong").Length() >= 1 {
					major = true
				}

				movie.AddJumpScare(matches[1], matches[2], major)
			}
		})

		s.Find("div.video-info-grid-column > p").Each(func(i int, s *goquery.Selection) {
			switch strings.Fields(s.Text())[0] {
			case "Reviews:":
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					href, ok := s.Attr("href")
					if ok {
						movie.AddReview(s.Text(), href)
					}
				})
			case "Director:":
				for _, d := range strings.Split(strings.Replace(s.Text(), "Director:", "", 1), ",") {
					movie.AddDirector(d)
				}
			case "Rating:":
				movie.SetContentRating(s.Text())
			case "Runtime:":
				movie.SetRuntimeFromPattern(regexp.MustCompile(`^Runtime:\s+(\d+)\sminutes$`), s.Text())
			case "Tags:":
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					movie.AddTag(s.Text())
				})
			}
		})
	})

	return movie, nil
}
