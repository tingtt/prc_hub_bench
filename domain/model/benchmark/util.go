package benchmark

import "os"

func mkdirIfNotExist(path string) error {
	_, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			// if dir not exists, mkdir
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
			os.Chown(path, 501, 20)
			if err != nil {
				return err
			}
		} else {
			// unhandled errors
			return err
		}
	}

	return nil
}

func writeFile(path string, b []byte) error {
	// Open file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Chown
	err = f.Chown(501, 20)
	if err != nil {
		return err
	}

	// Write
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
