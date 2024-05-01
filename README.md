# Cloudflare API Demo (Updating TXT records)

It updates for a TXT record in the given zone with the name `cfdns-test`
and updates its values to the current timestamp. If there is no such record,
it creates one.

## Usage

```bash
export CF_API_TOKEN="your-bearer-token"
export CF_ZONE_ID="your-zone-id"

go run .
```

## License

AGPL-3.0-or-later
