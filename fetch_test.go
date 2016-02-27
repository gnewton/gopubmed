package gopubmed

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	pmg := PubmedGetter{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 500,
			DisableKeepAlives:     false,
			DisableCompression:    false,
		},
	}

	pmids := []string{"24000000", "24000001", "24000002", "24000003"}
	articles, err := pmg.GetArticles(pmids)
	if err != nil {
		log.Fatal(err)
	}
	if len(articles) != len(pmids) {
		log.Fatal(err)
	}
}

func TestFetch_BadUrl(t *testing.T) {
	pmg := PubmedGetter{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 500,
			DisableKeepAlives:     false,
			DisableCompression:    false,
		},
		BaseUrl: "foobar",
	}

	pmids := []string{"24000000", "24000001", "24000002", "24000003"}
	_, err := pmg.GetArticles(pmids)
	if err == nil {
		log.Fatal(err)
	}

}

func TestFetch_BadPmids(t *testing.T) {
	pmg := PubmedGetter{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 500,
			DisableKeepAlives:     false,
			DisableCompression:    false,
		},
	}

	// This will not return an error from entrez, just empty xml
	pmids := []string{"a", "b", "c", "e"}
	articles, err := pmg.GetArticles(pmids)
	if err != nil {
		log.Fatal(err)
	}
	if len(articles) != 0 {
		log.Fatal(err)
	}

}

func TestFetch_GoodAndBadPmids(t *testing.T) {
	pmg := PubmedGetter{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 500,
			DisableKeepAlives:     false,
			DisableCompression:    false,
		},
	}

	// This will not return an error from entrez, ignores pmids it cannot find, even illegal ones
	pmids := []string{"a", "b", "c", "e", "24000000", "24000001", "24000002", "24000003"}
	articles, err := pmg.GetArticles(pmids)
	if err != nil {
		log.Fatal(err)
	}
	if len(articles) != 4 {
		log.Fatal(err)
	}

}

func TestFetch_Transport(t *testing.T) {
	pmg := PubmedGetter{
		Transport: nil,
	}

	pmids := []string{"24000000", "24000001", "24000002", "24000003"}
	_, err := pmg.GetArticles(pmids)
	if err == nil {
		log.Fatal(err)
	}

}
