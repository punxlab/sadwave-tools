package photos

import (
	"encoding/csv"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

type EventPhoto struct {
	EventUrl *url.URL
	PhotoUrl *url.URL
}

func FromCSV(path string) ([]*EventPhoto, error) {
	lines, err := readCSV(path)
	if err != nil {
		return nil, err
	}

	if err := validateCSV(lines); err != nil {
		return nil, err
	}

	res, err := convertToModel(lines)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "open file %s", path)
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", path)
	}

	return lines, nil
}

func validateCSV(lines [][]string) error {
	if len(lines) == 0 {
		return errors.New("empty csv")
	}

	for i, l := range lines {
		lineNumber := i + 1

		if len(l) != 2 {
			return errors.Errorf("incorrect format: line %d should contains 2 values", lineNumber)
		}

		if strings.TrimSpace(l[0]) == "" {
			return errors.Errorf("invalid value: the 1st value of line %d is required", lineNumber)
		}

		if !govalidator.IsURL(l[0]) {
			return errors.Errorf("invalid value: the 1st value of line %d should be an URL", lineNumber)
		}

		if strings.TrimSpace(l[1]) != "" {
			if !govalidator.IsURL(l[1]) {
				return errors.Errorf("invalid value: the 2d value of line %d can be empty or an URL", lineNumber)
			}

			if err := validateImageByUrl(l[1]); err != nil {
				return errors.Wrapf(err, "invalid value: the 2d value of line %d contains invalid URL", lineNumber)
			}
		}
	}

	return nil
}

func validateImageByUrl(url string) error {
	img, err := imageByURL(url)
	if err != nil {
		return err
	}

	b := img.Bounds()

	if b.Dy() <= 0 {
		return errors.Errorf("url %s provides an image with 0 height")
	}

	if b.Dx() <= 0 {
		return errors.Errorf("url %s provides an image with 0 width")
	}

	return nil
}

func imageByURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "url %s incorrect or not available", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("url %s is not available: status %s", url, resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "looks like provided url %s is not an image", url)
	}

	return img, nil
}

func convertToModel(lines [][]string) ([]*EventPhoto, error) {
	res := make([]*EventPhoto, 0, len(lines))
	for _, l := range lines {
		eventUrl, err := url.Parse(l[0])
		if err != nil {
			return nil, err
		}

		photoUrl, err := url.Parse(l[0])
		if err != nil {
			return nil, err
		}

		res = append(res, &EventPhoto{
			EventUrl: eventUrl,
			PhotoUrl: photoUrl,
		})
	}

	return res, nil
}
