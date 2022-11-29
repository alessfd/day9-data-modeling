package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	databaseUrl := "postgres://postgres:lmaoxd313@localhost:5433/Projects"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Koneksi ke database gagal: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Koneksi ke database berhasil!")
}
