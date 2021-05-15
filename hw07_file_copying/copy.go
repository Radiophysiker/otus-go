package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile          = errors.New("unsupported file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrRequiredFromPathOrToPath = errors.New("required fromPath and toPath")
	ErrNegativeLimit            = errors.New("limit cannot be negative")
	ErrNegativeOffset           = errors.New("offset cannot be negative")
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
	var totalNumberOfBytesToCopy int64
	/* вычислим суммарное количество байтов для копирования
	данное значение нужно для прогресса и для задания лимита для копирования */
	if limit == 0 || limit > sourceFileSize-offset {
		totalNumberOfBytesToCopy = sourceFileSize - offset
	} else {
		totalNumberOfBytesToCopy = limit
	}

	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	bar := pb.Full.Start64(totalNumberOfBytesToCopy)
	bar.Set(pb.Bytes, true)
	reader := io.LimitReader(rs, totalNumberOfBytesToCopy)
	barReader := bar.NewProxyReader(reader)
	_, err := io.Copy(w, barReader)
	if err != nil {
		return err
	}
	bar.Finish()
	return nil
}
