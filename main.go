package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/techdroplabs/dyspatch/service-template/dml"
	"github.com/techdroplabs/dyspatch/service-template/dml/block"
	"github.com/techdroplabs/dyspatch/service-template/dml/editable"
	"github.com/techdroplabs/dyspatch/service-template/dml/types"
)

func main() {
	blk, err := dml.Parse(context.Background(), `<dys-block>
  <dys-row>
    <dys-column>
      <dys-carousel>
        <dys-carousel-image src='https://picsum.photos/id/10/400/400' />
        <dys-carousel-image src='https://picsum.photos/id/160/400/400' />
        <dys-carousel-image src='https://picsum.photos/id/40/400/400' />
      </dys-carousel>
    </dys-column>
  </dys-row>
</dys-block>`)
	if err != nil {
		panic(err)
	}
	rr, err := dml.Render(context.Background(), []*block.Block{blk}, []map[string]*editable.Field{{}}, &types.RenderOptions{Strict: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(rr.AMPOutput.HTML)

	// Get our API key from the environment; configure.
	apiKey := os.Getenv("SPARKPOST_API_KEY")
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}
	var client sp.Client
	err = client.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	// Create a Transmission using an inline Recipient List
	// and inline email Content.
	tx := &sp.Transmission{
		Recipients: []string{"daniel@sendwithus.com"},
		Content: sp.Content{
			AMPHTML: rr.AMPOutput.HTML,
			HTML:    rr.HTMLOutput.HTML,
			From:    "dyspatch@email.dyspatch.io",
			Subject: "AMP Test",
		},
	}
	id, _, err := client.Send(tx)
	if err != nil {
		log.Fatal(err)
	}

	// The second value returned from Send
	// has more info about the HTTP response, in case
	// you'd like to see more than the Transmission id.
	log.Printf("Transmission sent with id [%s]\n", id)
}
