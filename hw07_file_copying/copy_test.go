package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	TestCases := []struct {
		testCaseID       string
		limit            int64
		offset           int64
		originFilePath   string
		comparedFilePath string
	}{
		{
			testCaseID:       "Test 1:",
			limit:            0,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit0.txt",
		},
		{
			testCaseID:       "Test 2:",
			limit:            10,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit10.txt",
		},
		{
			testCaseID:       "Test 3:",
			limit:            1000,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit1000.txt",
		},
		{
			testCaseID:       "Test 4",
			limit:            10000,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit10000.txt",
		},
		{
			testCaseID:       "Test 5",
			limit:            1000,
			offset:           100,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset100_limit1000.txt",
		},
		{
			testCaseID:       "Test 6",
			limit:            1000,
			offset:           6000,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset6000_limit1000.txt",
		},
		{
			testCaseID:       "Test 7",
			limit:            1000,
			offset:           6000,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset6000_limit1000.txt",
		},
	}

	for index, test := range TestCases {
		err := Copy(test.originFilePath, "testdata/temp.txt", test.limit, test.offset)
		if err != nil {
			t.Error(test.testCaseID, index, err)
		}
		result, err := compareTwoFiles("testdata/temp.txt", test.comparedFilePath)
		if err != nil {
			t.Error(test.testCaseID, index, err)
		}
		if result == false {
			t.Error(test.testCaseID, index, "files are not equal")
		}
		os.Remove("testdata/temp.txt")
	}
}

func TestCopyWithNegativeResult(t *testing.T) {
	TestNegativeCases := []struct {
		testCaseID     string
		limit          int64
		offset         int64
		originFilePath string
		expectedError  error
	}{
		{
			testCaseID:     "Test 1: negative limit",
			limit:          -100,
			offset:         0,
			originFilePath: "testdata/input.txt",
			expectedError:  ErrNegativeLimit,
		},
		{
			testCaseID:     "Test 2: offset exceeds file size",
			limit:          100,
			offset:         1000000,
			originFilePath: "testdata/input.txt",
			expectedError:  ErrOffsetExceedsFileSize,
		},
	}

	for index, test := range TestNegativeCases {
		err := Copy(test.originFilePath, "testdata/temp.txt", test.limit, test.offset)
		if err == nil || (err != nil && !errors.Is(err, test.expectedError)) {
			t.Error(test.testCaseID, index, "expected error: ", test.expectedError)
			continue
		}
		os.Remove("testdata/temp.txt")
	}
}

func compareTwoFiles(originFilePath, comparedFilePath string) (bool, error) {
	originFile, err := os.Open(originFilePath)
	if err != nil {
		return false, err
	}
	defer originFile.Close()
	comparedFile, err := os.Open(comparedFilePath)
	if err != nil {
		return false, err
	}
	defer comparedFile.Close()
	if limit == 0 {
		sourceFileInfo, err := originFile.Stat()
		if err != nil {
			return false, err
		}
		limit = sourceFileInfo.Size()
	}
	for {
		bufferSrc := make([]byte, DefaultBlockSize)
		_, errSrc := originFile.Read(bufferSrc)
		bufferTest := make([]byte, DefaultBlockSize)
		_, errTest := comparedFile.Read(bufferTest)
		if errSrc != nil || errTest != nil {
			if errors.Is(errSrc, ErrEOF) && errors.Is(errTest, ErrEOF) {
				return true, nil
			}
			return false, errTest
		}

		if !bytes.Equal(bufferSrc, bufferTest) {
			return false, nil
		}
	}
}
