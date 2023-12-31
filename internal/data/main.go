package data

// docker run -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -p 5432:5432 postgres

//go:generate xo schema "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -o ./ --single=schema.xo.go --src templates
//go:generate xo schema "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -o pg --single=schema.xo.go --src=pg/templates --go-context=both
