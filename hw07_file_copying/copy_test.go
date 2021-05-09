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
		expectedError    error
	}{
		{
			testCaseID:       "Test 1: copy all ",
			limit:            0,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit0.txt",
			expectedError:    nil,
		},
		{
			testCaseID:       "Test 2: copy all ",
			limit:            10,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit10.txt",
			expectedError:    nil,
		},
		{
			testCaseID:       "Test 3: copy all ",
			limit:            1000,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit1000.txt",
			expectedError:    nil,
		},
		{
			testCaseID:       "Test 1: copy all ",
			limit:            10000,
			offset:           0,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset0_limit10000.txt",
			expectedError:    nil,
		},
		{
			testCaseID:       "Test 1: copy all ",
			limit:            1000,
			offset:           100,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset100_limit1000.txt",
			expectedError:    nil,
		},
		{
			testCaseID:       "Test 1: copy all ",
			limit:            1000,
			offset:           6000,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset6000_limit1000.txt",
			expectedError:    nil,
		},
		{
			testCaseID:       "Test 1: copy all ",
			limit:            1000,
			offset:           6000,
			originFilePath:   "testdata/input.txt",
			comparedFilePath: "testdata/out_offset6000_limit1000.txt",
			expectedError:    nil,
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

func TestCopyingProcess(t *testing.T) {
	testdata := []byte("1234567890ABCDE")
	TestCases := []struct {
		testCaseID     string // Name of testcase
		limit          int64  // Number of bytes to copy
		offset         int64
		source         []byte // Offset in the source file
		expectedResult string // Expected result data
		expectedError  error  // Expected error
	}{
		{
			testCaseID:     "Test 1: copy all ",
			limit:          0,
			offset:         0,
			source:         testdata,
			expectedResult: "1234567890ABCDE",
			expectedError:  nil,
		},
		{
			testCaseID:     "Test 2: copy half from beginning",
			limit:          5,
			offset:         0,
			source:         testdata,
			expectedResult: "12345",
			expectedError:  nil,
		},
		{
			testCaseID:     "Test 3: copy half from middle",
			limit:          5,
			offset:         5,
			source:         testdata,
			expectedResult: "67890",
			expectedError:  nil,
		},
	}

	for _, test := range TestCases {
		reader := bytes.NewReader(test.source)
		writer := bytes.NewBuffer([]byte{})
		err := copyingProcess(reader, writer, int64(len(test.source)), test.offset, test.limit)
		if err != nil {
			t.Error("Error during test execution for test case - ", test.testCaseID, err)
		}
		if result := writer.String(); result != test.expectedResult {
			t.Error("Copied data does not equals the expected result", test.testCaseID, result, test.expectedResult)
		}
	}
}
