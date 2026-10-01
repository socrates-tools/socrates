package tty

import "os"

func Attached(stream any) bool {
	file, isFile := stream.(*os.File)
	if !isFile {
		return false
	}

	written, err := file.Stat()
	return err == nil && written.Mode()&os.ModeCharDevice != 0
}
