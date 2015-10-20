package db

import (
	"fmt"
	"testing"
)

func TestLoadProteinsFromMongo(t *testing.T) {
	// proteins := LoadProteinsFromMongo("proteindb_dev", "protein")
	// fmt.Println(len(proteins))
	proteins := loadProteinsFromBoltDB("/home/jgcarvalho/sync/data/multissdb/", "chameleonic.db", "proteins")
	fmt.Println("Teste")
	fmt.Println(len(proteins))
}

func TestGetProteins(t *testing.T) {
	db := DBConfig{
		Dir:    "/home/jgcarvalho/sync/data/multissdb/",
		Name:   "chameleonic.db",
		Bucket: "proteins",
		Init:   "Seq",
		Target: "All3",
	}
	id, start, end, e := GetProteins(db)
	fmt.Println(id)
	fmt.Println(start)
	fmt.Println(end)
	fmt.Println(e)
}
