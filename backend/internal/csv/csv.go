package csv

import (
	"encoding/csv"
	"fmt"
	"os"

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
