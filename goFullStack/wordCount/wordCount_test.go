package wordCount

import (
	"fmt"
	"sync"
	"testing"
)

type mockIoClient struct {
}

func (ioClient *mockIoClient) ReadFile(filename string) ([]byte, error) {
	return []byte("ABC"), nil
}

func TestGetWordCountInFile(t *testing.T) {
	wordCountChannel := make(chan wordcount)
	var wg sync.WaitGroup
	wg.Add(1)
	go getWordCountInFile("test", wordCountChannel, &mockIoClient{}, &wg)
	w := <-wordCountChannel
	if w.wordCountList[0] != "ABC:1" {
		t.Errorf("Error reading from channel")
	}
}

func TestProcessFileDataForWordcount(t *testing.T) {
	actualWordCount := processFileDataForWordcount("mango apple orange apple list ashish", "test1.txt")
	expectedWordCountList := []string{"mango:1", "apple:2", "orange:1", "list:1", "ashish:1"}
	fileName := "test1.txt"
	w := wordcount{fileName: fileName, err: nil, wordCountList: expectedWordCountList}

	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	fileNameCheck := fmt.Sprintf("Test if filename matches")
	t.Run(fileNameCheck, func(t *testing.T) {
		if actualWordCount.fileName != w.fileName {
			t.Errorf("Filename doesn't match")
		}
	})

	mapValueCheck := fmt.Sprintf("Test if map returning the correct value")
	t.Run(mapValueCheck, func(t *testing.T) {
		if !contains(actualWordCount.wordCountList, "mango:1") {
			t.Errorf("Map didn't result in correct value")
		}
	})

	mapValuesCheck := fmt.Sprintf("Test if map returning the correct value")
	t.Run(mapValuesCheck, func(t *testing.T) {
		if !contains(actualWordCount.wordCountList, "apple:2") && !contains(actualWordCount.wordCountList, "ashish:1") {
			t.Errorf("Map didn't result in correct value")
		}
	})

	errValueCheck := fmt.Sprintf("Test if any error is returned")
	t.Run(errValueCheck, func(t *testing.T) {
		if actualWordCount.err != nil {
			t.Errorf("Unable to obtain map. Function returned error")
		}
	})

}

func TestCheckForHiddenFile(t *testing.T) {
	actualResultForHiddenFile := isFileHidden(".DS_Store")
	actualResultForNormalFile := isFileHidden("test.txt")
	hiddenFileCheck := "Test for hidden file"
	t.Run(hiddenFileCheck, func(t *testing.T) {
		if !actualResultForHiddenFile {
			t.Errorf("Failed test case for filename: .DS_Store! Expected false")
		}
	})

	normalFileCheck := "Test for normal file"
	t.Run(normalFileCheck, func(t *testing.T) {
		if actualResultForNormalFile {
			t.Errorf("Failed test case for filename: test.txt! Expected false")
		}
	})
}
