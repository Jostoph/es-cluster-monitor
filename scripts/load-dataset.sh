host='localhost:9200'
curl -H 'Content-Type: application/x-ndjson' -XPOST "$host/bank/account/_bulk?pretty" --data-binary @scripts/data/accounts.json