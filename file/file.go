package file

import "os"

func Exists(fname string) bool {
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Read(fname string) (string, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Write(fname string, content string) error {
	return os.WriteFile(fname, []byte(content), 0644)
}
