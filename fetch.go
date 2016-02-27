package gopubmed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
)

const BASE_ENTREZ_URL_FETCH_PUBMED = "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&rettype=xml&id="

type PubmedGetter struct {
	Transport *http.Transport
	BaseUrl   string
}

func (pmg *PubmedGetter) GetArticles(pmids []string) ([]PubmedArticle, error) {
	articles, _, err := pmg.GetArticlesAndRaw(pmids)
	return articles, err

}

func (pmg *PubmedGetter) GetArticlesAndRaw(pmids []string) ([]PubmedArticle, []byte, error) {
	if len(pmids) == 0 {
		return nil, nil, errors.New("Error: Empty list of pmids")
	}
	if pmg.BaseUrl == "" {
		pmg.BaseUrl = BASE_ENTREZ_URL_FETCH_PUBMED
	}

	body, err := getPubmedArticlesRaw(pmids, pmg.Transport, pmg.BaseUrl)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	v := ArticleSet{}
	err = xml.Unmarshal(body, &v)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	if v.ArticleList == nil {
		var pma []PubmedArticle
		return pma, body, nil
	}
	return v.ArticleList, body, nil
}

func (pmg *PubmedGetter) GetArticlesRaw(pmids []string) ([]byte, error) {
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
	client := &http.Client{Transport: transport}
	req, err := http.NewRequest("GET", url, nil)
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