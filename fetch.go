package gopubmed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/gnewton/pubmedstruct"
	"log"
	"net/http"
)

const BASE_ENTREZ_URL_FETCH_PUBMED = "eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&rettype=xml&id="

type Fetcher struct {
	Transport *http.Transport
	BaseUrl   string
	Ssl       bool
}

func (pmg *Fetcher) GetArticles(pmids []string) ([]*pubmedstruct.PubmedArticle, error) {
	articles, _, err := pmg.GetArticlesAndRaw(pmids)
	return articles, err

}

func (pmg *Fetcher) GetArticlesAndRaw(pmids []string) ([]*pubmedstruct.PubmedArticle, []byte, error) {
	if len(pmids) == 0 {
		return nil, nil, errors.New("Error: Empty list of pmids")
	}
	if pmg.BaseUrl == "" {
		if pmg.Ssl {
			pmg.BaseUrl = "https://" + BASE_ENTREZ_URL_FETCH_PUBMED
		} else {
			pmg.BaseUrl = "http://" + BASE_ENTREZ_URL_FETCH_PUBMED
		}
	}

	body, err := getPubmedArticlesRaw(pmids, pmg.Transport, pmg.BaseUrl)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	v := pubmedstruct.PubmedArticleSet{}
	err = xml.Unmarshal(body, &v)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	if v.PubmedArticle == nil {
		var pma []*pubmedstruct.PubmedArticle
		return pma, body, nil
	}
	return v.PubmedArticle, body, nil
}

func (pmg *Fetcher) GetArticlesRaw(pmids []string) ([]byte, error) {
	_, body, err := pmg.GetArticlesAndRaw(pmids)
	return body, err
}

func getPubmedArticlesRaw(pmids []string, transport *http.Transport, baseUrl string) ([]byte, error) {
	if len(pmids) == 0 {
		return nil, errors.New("Error: Empty list of pmids")
	}

	if baseUrl == "" {
		return nil, errors.New("Error: Pubmed entrez URL is empty")
	}

	if transport == nil {
		return nil, errors.New("Error: Transport cannot be nil")
	}

	url := makeUrl(baseUrl, pmids)

	if Debug {
		log.Println("Getting: ", url)
	}

	client := &http.Client{Transport: transport}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux i686; rv:10.0) Gecko/20100101 Firefox/10.0")
	req.Close = true
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error opening url:", url, "   error=", err)
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

var Debug = false

func makeUrl(baseUrl string, pmids []string) string {
	url := baseUrl
	for i := 0; i < len(pmids); i++ {
		if pmids[i] != "" {
			if i != 0 {
				url += ","
			}
			url += pmids[i]
		}
	}
	return url
}
