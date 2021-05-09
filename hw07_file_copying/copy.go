package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const (
	DefaultBlockSize = 4096
)

var (
	ErrUnsupportedFile          = errors.New("unsupported file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrRequiredFromPathOrToPath = errors.New("required fromPath and toPath")
	ErrNegativeLimit            = errors.New("limit cannot be negative")
	ErrNegativeOffset           = errors.New("offset cannot be negative")
	ErrEOF                      = io.EOF
)

func Copy(fromPath, toPath string, limit, offset int64) error {
	err := validationOfCopyArguments(fromPath, toPath, limit, offset)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	sourceFileSize, err := getSizeFile(srcFile)
	if err != nil {
		return err
	}
	if sourceFileSize == 0 {
		return ErrUnsupportedFile
	}
	if offset > sourceFileSize {
		return ErrOffsetExceedsFileSize
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	err = copyingProcess(srcFile, dstFile, sourceFileSize, offset, limit)
	if err != nil {
		return err
	}
	return nil
}

func getSizeFile(file *os.File) (int64, error) {
	sourceFileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return sourceFileInfo.Size(), nil
}

func validationOfCopyArguments(fromPath, toPath string, limit, offset int64) error {
	if fromPath == "" || toPath == "" {
		return ErrRequiredFromPathOrToPath
	}
	if limit < 0 {
		return ErrNegativeLimit
	}
	if offset < 0 {
		return ErrNegativeOffset
	}
	return nil
}

func copyingProcess(rs io.ReadSeeker, w io.Writer, sourceFileSize, offset, limit int64) error {
	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	var totalAmountOfBytes int64
	switch {
	case limit != 0:
		totalAmountOfBytes = limit
	default:
		totalAmountOfBytes = sourceFileSize - offset
	}

	buffer := make([]byte, DefaultBlockSize)
	lmr := io.LimitReader(rs, totalAmountOfBytes)

	pBar := pb.Start64(totalAmountOfBytes)
	pBar.SetWidth(100)
	pBar.Set(pb.Bytes, true)
	for {
		bytesRead, errRead := lmr.Read(buffer)
		if bytesRead > 0 {
			_, errWrite := w.Write(buffer[:bytesRead])

			if errWrite != nil {
				return errWrite
			}
			pBar.Add(bytesRead)
		}
		if errors.Is(errRead, ErrEOF) {
			pBar.Finish()
			return nil
		}
		if errRead != nil {
			pBar.Finish()
			return errRead
		}
	}
}
