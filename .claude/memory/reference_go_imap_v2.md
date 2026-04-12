---
name: go-imap/v2 imapclient API reference
description: Correct method signatures and types for github.com/emersion/go-imap/v2/imapclient used in this project
type: reference
---

# go-imap/v2 imapclient API

## Search

```go
// UID-based search (capital UID, not Uid)
func (c *Client) UIDSearch(criteria *imap.SearchCriteria, options *imap.SearchOptions) *SearchCommand

// Extract UIDs from result
searchData, err := c.UIDSearch(criteria, nil).Wait()
uids := searchData.AllUIDs() // returns []imap.UID
```

## Fetch

Não existe `UIDFetch`. Existe só `Fetch`, que aceita `imap.NumSet` (interface):

```go
func (c *Client) Fetch(numSet imap.NumSet, options *imap.FetchOptions) *FetchCommand
```

- Passar `imap.UIDSet` → envia `UID FETCH` (por UID permanente)
- Passar `imap.SeqSet` → envia `FETCH` (por número de sequência, muda se emails deletados)

## UIDSet

```go
// Criar UIDSet a partir de slice de UIDs
uidSet := imap.UIDSetNum(uids...)  // uids é []imap.UID

// Adicionar incrementalmente
var s imap.UIDSet
s.AddNum(uid1, uid2)
s.AddRange(start, stop)
```

## Fluxo completo (buscar emails não lidos por UID)

```go
criteria := &imap.SearchCriteria{
    NotFlag: []imap.Flag{imap.FlagSeen},
}

searchData, err := c.UIDSearch(criteria, nil).Wait()
// searchData.AllUIDs() → []imap.UID

uidSet := imap.UIDSetNum(searchData.AllUIDs()...)

fetchOptions := &imap.FetchOptions{Envelope: true, Flags: true}
messages, err := c.Fetch(uidSet, fetchOptions).Collect()
// messages → []*imapclient.FetchMessageBuffer
```

## FetchMessageBuffer campos úteis

```go
message.UID              // imap.UID (uint32)
message.Envelope.Subject // string
message.Envelope.From    // []imap.Address
message.Envelope.Date    // time.Time
message.Flags            // []imap.Flag
```

## imap.Address

```go
address.Name    // nome display
address.Mailbox // parte antes do @
address.Host    // parte depois do @
```
