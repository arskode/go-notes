# Simple API built with golang

## API

* `POST /notes` - add a note

```go
type Note struct {
	ID          uint64 `json:"id"`  // read only
	Title       string `json:"title"`
	Description string `json:"description"`
}
```

* `GET /notes/{noteID}` - get a single note

* `GET /notes` - get all notes

* `PUT /notes/{noteID}` - edit note

```go
type Note struct {
	ID          uint64 `json:"id"`  // read only
	Title       string `json:"title"`
	Description string `json:"description"`
}
```

* `DELETE /notes/{noteID}` - delete note


## Development

`docker-compose up`


## Testing

`make test`
