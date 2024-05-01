package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/dns"
	"github.com/cloudflare/cloudflare-go/v2/zones"
)

type CloudflareTXTRecordUpdater struct {
	client *cloudflare.Client

	zoneID string
}

func NewCloudflareTXTRecordUpdater(client *cloudflare.Client, zoneID string) *CloudflareTXTRecordUpdater {
	return &CloudflareTXTRecordUpdater{
		client: client,
		zoneID: zoneID,
	}
}

func (u *CloudflareTXTRecordUpdater) create(ctx context.Context, name, content string) error {
	_, err := u.client.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(u.zoneID),
		Record: dns.TXTRecordParam{
			Name:    cloudflare.F(name),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
			Content: cloudflare.F(content),
		},
	})
	if err != nil {
		return fmt.Errorf("create TXT record: %w", err)
	}

	return nil
}

func (u *CloudflareTXTRecordUpdater) update(ctx context.Context, recordID, name, content string) error {
	_, err := u.client.DNS.Records.Update(ctx, recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(u.zoneID),
		Record: dns.TXTRecordParam{
			Name:    cloudflare.F(name),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
			Content: cloudflare.F(content),
		},
	})
	if err != nil {
		return fmt.Errorf("update TXT record: %w", err)
	}

	return nil
}

var errRecordNotFound = errors.New("record not found")

func (u *CloudflareTXTRecordUpdater) find(ctx context.Context, name string) (string, error) {
	zone, err := u.client.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cloudflare.F(u.zoneID),
	})
	if err != nil {
		return "", fmt.Errorf("get zone: %w", err)
	}

	// The "name" responded from Cloudflare API is the fully-qualified domain name.
	// Therefore, we need to join the name with the zone name to get the fully-qualified domain name.
	fullQualifiedName := name + "." + zone.Name
	if name == "@" {
		fullQualifiedName = zone.Name
	}

	res, err := u.client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(u.zoneID),
		Type:   cloudflare.F(dns.RecordListParamsTypeTXT),
		Name:   cloudflare.F(fullQualifiedName),
	})
	if err != nil {
		var cfErr *cloudflare.Error

		if !errors.As(err, &cfErr) {
			return "", fmt.Errorf("request error: %w", err)
		}

		if cfErr.StatusCode != 404 {
			return "", fmt.Errorf("API error: %w", err)
		}

		return "", errRecordNotFound
	}

	if len(res.Result) < 1 {
		return "", errRecordNotFound
	}

	return res.Result[0].ID, nil
}

func (u *CloudflareTXTRecordUpdater) UpdateOrCreate(ctx context.Context, name, value string) error {
	// check if there has been a TXT record with the same name
	recordID, err := u.find(ctx, name)
	if err != nil {
		if !errors.Is(err, errRecordNotFound) {
			return err
		}

		// create a new TXT record
		return u.create(ctx, name, value)
	}

	// update the existing TXT record
	return u.update(ctx, recordID, name, value)
}
