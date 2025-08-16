env "local" {
  url = "postgres://socialdev:root@localhost:5432/socialdev?sslmode=disable"
  dev = "postgres://socialdev:root@localhost:5432/socialdev_dev?sslmode=disable"
  src = "./internal/db/schema.sql"

  migration {
    dir = "file://internal/db/migrations"
  }
} 
