package proteindb

import (
	"fmt"
	"testing"
)

func TestLoadProteinsFromMongo(t *testing.T) {
	proteins := LoadProteinsFromMongo("proteindb_dev", "protein")
	fmt.Println(len(proteins))

}
