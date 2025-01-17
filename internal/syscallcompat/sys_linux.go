// Package syscallcompat wraps Linux-specific syscalls.
package syscallcompat

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	_FALLOC_FL_KEEP_SIZE = 0x01

	// O_DIRECT means oncached I/O on Linux. No direct equivalent on MacOS and defined
	// to zero there.
	O_DIRECT = syscall.O_DIRECT

	// O_PATH is only defined on Linux
	O_PATH = unix.O_PATH

	// Only defined on Linux
	RENAME_NOREPLACE = unix.RENAME_NOREPLACE
	RENAME_WHITEOUT  = unix.RENAME_WHITEOUT
	RENAME_EXCHANGE  = unix.RENAME_EXCHANGE
)

// EnospcPrealloc preallocates ciphertext space without changing the file
// size. This guarantees that we don't run out of space while writing a
// ciphertext block (that would corrupt the block).
func EnospcPrealloc(fd int, off int64, len int64) (err error) {
	for {
		err = syscall.Fallocate(fd, _FALLOC_FL_KEEP_SIZE, off, len)
		if err == syscall.EINTR {
			// fallocate, like many syscalls, can return EINTR. This is not an
			// error and just signifies that the operation was interrupted by a
			// signal and we should try again.
			continue
		}
		if err == syscall.EOPNOTSUPP {
			// ZFS and ext3 do not support fallocate. Warn but continue anyway.
			// https://github.com/rfjakob/gocryptfs/issues/22
			return nil
		}
		return err
	}
}

// Getdents syscall with "." and ".." filtered out.
func Getdents(fd int) ([]DirEntry, error) {
	entries, _, err := getdents(fd)
	return entries, err
}
