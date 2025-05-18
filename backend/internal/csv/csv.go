package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/glup3/trendingrepos/internal/api"
)

func ToCSV(filename string, repos []api.Repo) error {
	f, err := os.Create(fmt.Sprintf("%s.csv", filename))
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	for i, repo := range repos {
		line := repo.CSVRecord()
		if i == 0 {
			line = repo.CSVHeader()
		}
		if err := w.Write(line); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

func ReadCsvFile(filePath string) ([]api.Repo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	repos := make([]api.Repo, len(records))
	for i, r := range records {
		if i == 0 {
			continue
		}
		stars, err := strconv.Atoi(r[2])
		if err != nil {
			return nil, err
		}
		repos[i] = api.Repo{
			Id:              r[0],
			NameWithOwner:   r[1],
			Stars:           stars,
			PrimaryLanguage: r[3],
			Description:     r[4],
		}
	}

	return repos, err
}
