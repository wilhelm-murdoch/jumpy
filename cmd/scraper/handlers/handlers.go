package handlers

import (
	"errors"
	"fmt"
	"os"

	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/library"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/logger"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/models"
	"github.com/wilhelm-murdoch/jumps.care/cmd/scraper/processor"
)

func HandleMovieDetails(job processor.Job[models.Movie]) (processor.Job[models.Movie], error) {
	path := fmt.Sprintf("dist/movies/%s.json", job.Body.Id)

	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		logger.Warning(fmt.Sprintf("file for %s already exists; skipping ...", path))
		return job, nil
	}

	movie, err := library.GetMovieDetailsFromUrl(job.Body)
	if err != nil {
		return job, err
	}

	if err = movie.Save(path, movie); err != nil {
		return job, err
	}

	return job, nil
}

func HandleMoviePoster(job processor.Job[models.Movie]) (processor.Job[models.Movie], error) {
	// for _, r := range movie.Reviews {
	// 	if strings.ToLower(r.Name) == "imdb" {
	// 		pattern := regexp.MustCompile(`title/tt(\d+)/?$`)
	// 		matches := pattern.FindStringSubmatch(r.Url)

	// 		if len(matches) == 2 {
	// 			document, err = GetDocumentFromUrl("https://www.movieposterdb.com/-i" + matches[1])
	// 			if err != nil {
	// 				return movie, err
	// 			}

	// 			document.Find("div.container > figure > a > img.img-responsive").Each(func(i int, s *goquery.Selection) {
	// 				fmt.Println(r.Name)
	// 				src, ok := s.Attr("src")
	// 				if ok {
	// 					fmt.Println(src)
	// 				}
	// 			})
	// 			time.Sleep(time.Second)
	// 		}
	// 	}
	// }
	return job, nil
}

func HandleMovieSubtitles(job processor.Job[models.Movie]) (processor.Job[models.Movie], error) {
	// var output []string
	// var prefix string = "Minor "

	// subtitle := "jump scare ahead!"
	// for i, j := range m.JumpScares {
	// 	if spoilers {
	// 		subtitle = "- " + j.Spoiler
	// 	}

	// 	if j.Major {
	// 		prefix = "Major "
	// 	}

	// 	output = append(
	// 		output,
	// 		fmt.Sprint(i+1),                // The index of the subtitle
	// 		j.TimeStart+" --> "+j.TimeStop, // The start and stop timestamps of the subtitle separated by " --> "
	// 		prefix+subtitle+"\n",           // The subtitle text prefixed with "Major" or "Minor"
	// 	)
	// }

	// file, err := os.Create(path)
	// if err != nil {
	// 	return err
	// }

	// _, err = file.WriteString(strings.Join(output, "\n"))
	// if err != nil {
	// 	return err
	// }

	// defer file.Close()

	return job, nil
}

func HandleMovieDetailsAgain(job processor.Job[models.Movie]) (processor.Job[models.Movie], error) {
	job.Body.Title += "again"
	return job, nil
}
