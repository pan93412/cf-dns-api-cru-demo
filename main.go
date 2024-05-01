package main

import (
	"context"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v2"
	cfoption "github.com/cloudflare/cloudflare-go/v2/option"
)

func main() {
	client := cloudflare.NewClient(
		cfoption.WithAPIToken(os.Getenv("CF_API_TOKEN")),
	)

	updater := NewCloudflareTXTRecordUpdater(client, os.Getenv("CF_ZONE_ID"))
	if err := updater.UpdateOrCreate(context.Background(), "cfdns-test", time.Now().String()); err != nil {
		panic(err)
	}
}
