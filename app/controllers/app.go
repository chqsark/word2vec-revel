package controllers

import (
	"fmt"
	"github.com/revel/revel"
	wordmodel "github.com/thinxer/go-word2vec"
	"log"
	"strings"
	"wordapp/app"
)

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

type App struct {
	*revel.Controller
}

type Pair struct {
	Word string
	Sim  float32
}

func (c App) Index() revel.Result {
	log.Printf("visitor from %s\n", strings.Split(c.Request.RemoteAddr, ":")[0])
	return c.Render()
}

func mostSimilar(query string, n int, m *wordmodel.Model) ([]Pair, error) {
	vec := wordmodel.Vector(make([]float32, m.Layer1Size))
	if wordId, ok := m.Vocab[query]; !ok {
		return nil, fmt.Errorf("word not found: %s", query)
	} else {
		vec.Add(1, m.Vector(wordId))
	}
	vec.Normalize()

	r := make([]Pair, n)
	for w, i := range m.Vocab {
		sim := vec.Dot(m.Vector(i))
		this := Pair{w, sim}
		for j := 0; j < n; j++ {
			if this.Sim > r[j].Sim {
				this, r[j] = r[j], this
			}
		}
	}
	return r[1:], nil
}

func handleQueryMultiple(query string) []wordmodel.Pair {
	model := app.Model
	log.Println(query)
	keywords := strings.Split(query, ";")
	if len(keywords) < 2 {
		log.Println("query parsing error, or input not valid")
	}
	positive_words := strings.Split(keywords[0], ",")
	negative_words := strings.Split(keywords[1], ",")
	results, err := model.MostSimilar(positive_words, negative_words, 50)
	if err != nil {
		log.Println(err)
	}
	return results
}

func (c App) HandleQuery(query string) revel.Result {
	c.Validation.MaxSize(query, 2*KB).Message("word too long")
	log.Println(query)
	if strings.Contains(query, ";") {
		results := handleQueryMultiple(query)
		type data struct {
			Results []wordmodel.Pair
		}
		d := data{Results: results}
		return c.Render(d)
	} else {
		results := handleQuerySingle(query)
		type data struct {
			Results []Pair
		}
		d := data{Results: results}
		return c.Render(d)
	}
	return c.Render()
}

func handleQuerySingle(query string) []Pair {
	model := app.Model
	log.Println(query)
	results, err := mostSimilar(query, 50, model)
	if err != nil {
		log.Println(err)
	}
	return results
}
