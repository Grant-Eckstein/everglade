package everglade

import "os"

func (e *Everglade) Add(path string) error {
	// Verify that file exists
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Add path to list
	e.Paths = append(e.Paths, path)
	return nil
}
