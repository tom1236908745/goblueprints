package main

import (
	"bufio"
	"fmt"
	"github.com/tom1236908745/thesaurus"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHugh{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("\""+word+"\"の検索に失敗しました。", err)
		}
		if len(syns) == 0 {
			log.Fatalln("\"" + word + "\"に類義語が見つかりませんでした。")
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
