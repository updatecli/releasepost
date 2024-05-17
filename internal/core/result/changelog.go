package result

import "path/filepath"

const (
	// Created state
	CREATEDSTATE = "created"
	// Modified state
	MODIFIEDSTATE = "modified"
	// UnModified state
	UNMODIFIEDSTATE = "unmodified"
	// UnTracked state
	UNTRACKEDSTATE = "untracked"
)

type Changelog struct {
	Created    []string `json:"created"`
	Modified   []string `json:"modified"`
	UnModified []string `json:"unmodified"`
	UnTracked  []string `json:"untracked"`
}

func (c Changelog) ExitCode() int {
	if len(c.Created) == 0 && len(c.Modified) == 0 {
		return 1
	}
	return 0
}

// UpdateResult update the result with the current changelog information
func (c Changelog) UpdateResult(result *Result) error {

	parser := func(files []string, state string) {
		for _, f := range files {
			dirname := filepath.Dir(f)

			if result.Dir[dirname].IsFileExist(f) {
				continue
			}

			if result.Dir == nil {
				result.Dir = make(map[string]FileResults)
			}

			result.Dir[dirname] = append(
				result.Dir[dirname],
				FileResult{
					Path:  f,
					State: state,
				})
		}
	}

	parser(c.Created, CREATEDSTATE)
	parser(c.Modified, MODIFIEDSTATE)
	parser(c.UnModified, UNMODIFIEDSTATE)

	return nil
}
