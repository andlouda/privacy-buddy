package backend

import "os"

func isDevMode() bool {
	return os.Getenv("WAILS_DEV") == "1"
}
