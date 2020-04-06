package wordCount

import (
	"fmt"
	"goFullStack/fileReader"
	"goFullStack/sortMap"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type wordcount struct {
	fileName      string
	err           error
	wordCountList []string
}

type Reader interface {
	ReadFile(filename string) ([]byte, error)
}

func ProcessDir(dirPath string) {
	wordCountChannel := make(chan wordcount, 4)
	processDirForWordcount(dirPath, wordCountChannel, fileReader.NewClient())
}

func processDirForWordcount(dirPath string, wordCountChannel chan wordcount, ioClient Reader) {
	var wg sync.WaitGroup
	err := filepath.Walk(dirPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if isFileHidden(fileInfo.Name()) || fileInfo.IsDir() {
				return nil
			} else {
				wg.Add(1)
				go getWordCountInFile(path, wordCountChannel, ioClient, &wg)
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	wg.Wait()
	close(wordCountChannel)
	readWordCountFromChannel(wordCountChannel)
}

func getWordCountInFile(filePath string, wordCountChannel chan wordcount, ioClient Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	fileData, err := ioClient.ReadFile(filePath)
	if err != nil {
		w := wordcount{fileName: filePath, err: err, wordCountList: nil}
		wordCountChannel <- w
		return
	}
	wordCountChannel <- processFileDataForWordcount(string(fileData), filePath)
}

func processFileDataForWordcount(fileData string, filePath string) wordcount {
	wordList := strings.Fields(fileData)
	wordCountMap := make(map[string]int)
	for _, word := range wordList {
		if count, ok := wordCountMap[word]; ok {
			wordCountMap[word] = count + 1
		} else {
			wordCountMap[word] = 1
		}
	}
	wordCountList := convertWordcountMapToList(wordCountMap)
	return wordcount{fileName: filePath, err: nil, wordCountList: wordCountList}
}

func convertWordcountMapToList(wordCountMap map[string]int) []string {
	wordCountList := []string{}
	pairList := sortMap.SortMapByValue(wordCountMap)
	for _, pair := range pairList {
		v := pair.Key + ":" + strconv.Itoa(pair.Value)
		wordCountList = append(wordCountList, v)
	}
	return wordCountList
}

func readWordCountFromChannel(wordCountChannel chan wordcount) {
	for wordCountObj := range wordCountChannel {
		if wordCountObj.err != nil {
			fmt.Println(wordCountObj.fileName + " : " + wordCountObj.err.Error())
		} else {
			fmt.Println(wordCountObj.fileName + " : " + strings.Join(wordCountObj.wordCountList, ", ") + "\n")
		}
	}
}

func isFileHidden(fileName string) bool {
	return fileName[0:1] == "."
}
