package db

import (
	// "github.com/jgcarvalho/zeca/ca"
	"encoding/json"
	"log"
	"reflect"
	"strings" // "gopkg.in/mgo.v2"

	"github.com/boltdb/bolt"
	//"labix.org/v2/mgo/bson"
	"fmt"
)

// type Protein struct {
// 	Pdb_id string  "pdb_id"
// 	Chains []Chain "chains_data"
// }

type Config struct {
	Dir    string `toml:"db-dir"`
	Name   string `toml:"db-name"`
	Bucket string `toml:"bucket-name"`
	Init   string `toml:"init"`
	Target string `toml:"target"`
}

type Protein struct {
	ID string
	//original
	Seq    string
	Dssp   string
	Stride string
	Kaksi  string
	Pross  string
	//processed
	Dssp3   string
	Stride3 string
	Kaksi3  string
	Pross3  string
	// consensus 2
	DsspStride3  string
	DsspKaksi3   string
	DsspPross3   string
	StrideKaksi3 string
	StridePross3 string
	KaksiPross3  string
	// consensus 3
	DsspStrideKaksi3  string
	DsspStridePross3  string
	DsspKaksiPross3   string
	StrideKaksiPross3 string
	// consensus 4
	All3 string
}

func loadProteinsFromBoltDB(dirname, dbname, bucket string) []Protein {
	db, err := bolt.Open(dirname+dbname, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	var result []Protein
	var prot Protein
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.ForEach(func(k, v []byte) error {
			err := json.Unmarshal(v, &prot)
			if err != nil {
				fmt.Println("DB error:", err)
			} else {
				result = append(result, prot)
			}

			// fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
	fmt.Println(result)
	return result
}

func (c *Protein) getField(field string) string {
	r := reflect.ValueOf(c)
	s := reflect.Indirect(r).FieldByName(field)
	return s.String()
}

func GetProteins(db Config) (id, start, end string, e error) {
	proteins := loadProteinsFromBoltDB(db.Dir, db.Name, db.Bucket)
	id = "all"
	start = "#"
	end = "#"
	for i := 0; i < len(proteins); i++ {
		start += proteins[i].getField(strings.Title(db.Init)) + "#"
		end += proteins[i].getField(strings.Title(db.Target)) + "#"
		// start += proteins[i].Chains[0].getField(strings.Title(db.Init)) + "#"
		// end += proteins[i].Chains[0].getField(strings.Title(db.Target)) + "#"
	}
	if len(start) != len(end) {
		e = fmt.Errorf("Error: Number of CA start cells is different from end cells")
	}
	return id, start, end, e
}
