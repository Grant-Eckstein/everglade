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

func (e *Everglade) Encrypt() error {
	// Repeat for each layer
	for i := 0; i < e.Layers; i++ {
		// For each file, read in, encrypt, and write
		for _, f := range e.Paths {
			// Read in file
			d, err := os.ReadFile(f)
			if err != nil {
				return err
			}

			// Encrypt the file
			ct, err := e.Blind.AES.CBC.Encrypt(d)
			if err != nil {
				return err
			}

			// Write the file
			err = os.WriteFile(f, ct, 600)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func (e *Everglade) Decrypt() error {
	// Repeat for each layer
	for i := 0; i < e.Layers; i++ {
		// For each file, read in, decrypt, and write
		for _, f := range e.Paths {
			// Read in file
			d, err := os.ReadFile(f)
			if err != nil {
				return err
			}

			// Decrypt the file
			pt, err := e.Blind.AES.CBC.Decrypt(d)
			if err != nil {
				return err
			}

			// Write the file
			err = os.WriteFile(f, pt, 666)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}
