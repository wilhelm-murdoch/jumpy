package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wilhelm-murdoch/go-batch"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/library"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/logger"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/models"
)

func HandleMovieDetails(j *batch.Job[models.Movie]) (*batch.Job[models.Movie], error) {
	document, err := library.GetDocumentFromUrl(j.Body.SourceUrl)
	if err != nil {
		return nil, err
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

				j.Body.AddJumpScare(matches[1], matches[2], major)
			}
		})

		s.Find("div.entry-content > p:nth-child(4)").Each(func(i int, s *goquery.Selection) {
			switch strings.Fields(s.Text())[0] {
			case "Synopsis:":
				j.Body.SetSynopsis(s.Text())
			}
		})

		s.Find("div.video-info-grid-column > p").Each(func(i int, s *goquery.Selection) {
			switch strings.Fields(s.Text())[0] {
			case "Reviews:":
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					href, ok := s.Attr("href")
					if ok {
						j.Body.AddReview(s.Text(), href)
					}
				})
			case "Director:":
				for _, d := range strings.Split(strings.Replace(s.Text(), "Director:", "", 1), ",") {
					j.Body.AddDirector(d)
				}
			case "Rating:":
				j.Body.SetContentRating(s.Text())
			case "Runtime:":
				j.Body.SetRuntimeFromPattern(regexp.MustCompile(`^Runtime:\s+(\d+)\sminutes$`), s.Text())
			case "Tags:":
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					j.Body.AddTag(s.Text())
				})
			}
		})
	})

	return j, nil
}

func HandleMoviePoster(j *batch.Job[models.Movie]) (*batch.Job[models.Movie], error) {
	path := fmt.Sprintf("dist/images/covers/%s.jpg", j.Body.Id)

	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		logger.Warning(fmt.Sprintf("file for %s already exists; skipping ...", path))
		return j, nil
	}

	for _, r := range j.Body.Reviews {
		if strings.ToLower(r.Name) == "imdb" {
			pattern := regexp.MustCompile(`(\d+)`)
			matches := pattern.FindStringSubmatch(r.Url)

			if len(matches) == 2 {
				document, _ := library.GetDocumentFromUrl("https://www.movieposterdb.com/-i" + matches[1])

				document.Find("body > div.section > div.container > div.row.mt-5.mb-4 > div.col-md-4.ml-auto > figure > a > img").Each(func(i int, s *goquery.Selection) {
					src, ok := s.Attr("src")
					if ok {
						img, _ := os.Create(path)
						defer img.Close()

						response, _ := http.Get(src)
						defer response.Body.Close()

						_, err := io.Copy(img, response.Body)
						if err != nil {
							logger.Warning("could not save %s locally; skipping ...", path)
						} else {
							logger.Info("%s written to disk.", path)
						}
					}
				})
				time.Sleep(time.Second)
			}
		}
	}
	return j, nil
}

// func HandleMovieSubtitles(j batch.Job[models.Movie]) (batch.Job[models.Movie], error) {
// 	// var output []string
// 	// var prefix string = "Minor "

// 	// subtitle := "jump scare ahead!"
// 	// for i, j := range m.JumpScares {
// 	// 	if spoilers {
// 	// 		subtitle = "- " + j.Spoiler
// 	// 	}

// 	// 	if j.Major {
// 	// 		prefix = "Major "
// 	// 	}

// 	// 	output = append(
// 	// 		output,
// 	// 		fmt.Sprint(i+1),                // The index of the subtitle
// 	// 		j.TimeStart+" --> "+j.TimeStop, // The start and stop timestamps of the subtitle separated by " --> "
// 	// 		prefix+subtitle+"\n",           // The subtitle text prefixed with "Major" or "Minor"
// 	// 	)
// 	// }

// 	// file, err := os.Create(path)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// _, err = file.WriteString(strings.Join(output, "\n"))
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// defer file.Close()

// 	return j, nil
// }
