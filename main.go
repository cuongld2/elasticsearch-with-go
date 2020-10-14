package main

import (
	"context"
	"elasticsearch-with-go/es"
	_ "elasticsearch-with-go/es"
	"elasticsearch-with-go/model"
	_ "elasticsearch-with-go/model"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func main() {

	ctx := context.Background()
	esclient, err := es.GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	var genes []model.Gene

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("gene-name", "Testi"))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
	}
	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
	/* until this block */

	searchService := esclient.Search().Index("gene").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		var gene model.Gene
		err := json.Unmarshal(hit.Source, &gene)
		if err != nil {
			fmt.Println("[Getting Genes][Unmarshal] Err=", err)
		}

		genes = append(genes, gene)
	}

	if err != nil {
		fmt.Println("Fetching gene fail: ", err)
	} else {
		for _, s := range genes {
			fmt.Printf("Gene found GeneID: %s, GeneName: %s, Description: %s \n", s.GeneId, s.GeneName, s.Description)
		}
	}

}
