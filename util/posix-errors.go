package util

import (
	"os"
	"runtime"
	"syscall"
)

// Function not implemented error
func isSysErrNoSys(err error) bool {
	return err == syscall.ENOSYS
}

// Not supported error
func isSysErrOpNotSupported(err error) bool {
	return err == syscall.EOPNOTSUPP
}

// No space left on device error
func isSysErrNoSpace(err error) bool {
	return err == syscall.ENOSPC
}

// Input/output error
func isSysErrIO(err error) bool {
	return err == syscall.EIO
}

// Check if the given error corresponds to EISDIR (is a directory).
func isSysErrIsDir(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		switch pathErr.Err {
		case syscall.EISDIR:
			return true
		}
	}
	return false
}

// Check if the given error corresponds to ENOTDIR (is not a directory).
func isSysErrNotDir(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		switch pathErr.Err {
		case syscall.ENOTDIR:
			return true
		}
	}
	return false
}

// Check if the given error corresponds to the ENAMETOOLONG (name too long).
func isSysErrTooLong(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		switch pathErr.Err {
		case syscall.ENAMETOOLONG:
			return true
		}
	}
	return false
}

// Check if the given error corresponds to ENOTEMPTY for unix
// and ERROR_DIR_NOT_EMPTY for windows (directory not empty).
func isSysErrNotEmpty(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if runtime.GOOS == "windows" {
			if errno, _ok := pathErr.Err.(syscall.Errno); _ok && errno == 0x91 {
				// ERROR_DIR_NOT_EMPTY
				return true
			}
		}
		switch pathErr.Err {
		case syscall.ENOTEMPTY:
			return true
		}
	}
	return false
}

// Check if the given error corresponds to the specific ERROR_PATH_NOT_FOUND for windows
func isSysErrPathNotFound(err error) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	if pathErr, ok := err.(*os.PathError); ok {
		if errno, _ok := pathErr.Err.(syscall.Errno); _ok && errno == 0x03 {
			// ERROR_PATH_NOT_FOUND
			return true
		}
	}
	return false
}

// Check if the given error corresponds to the specific ERROR_INVALID_HANDLE for windows
func isSysErrHandleInvalid(err error) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	// Check if err contains ERROR_INVALID_HANDLE errno
	if errno, ok := err.(syscall.Errno); ok && errno == 0x6 {
		return true
	}
	return false
}

func isSysErrCrossDevice(err error) bool {
	if e, ok := err.(*os.LinkError); ok {
		return e.Err == syscall.EXDEV
	}
	return false
}
