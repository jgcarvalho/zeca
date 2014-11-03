package proteindb

import (
	// "github.com/jgcarvalho/zeca/ca"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2"
	//"labix.org/v2/mgo/bson"
	"fmt"
)

type Protein struct {
	Pdb_id string  "pdb_id"
	Chains []Chain "chains_data"
}

type Chain struct {
	Chain                        string "chain"
	Seq_pdb                      string "seq_pdb"
	Seq_viz                      string "seq_viz"
	Ss_dssp                      string "ss_dssp"
	Ss_stride                    string "ss_stride"
	Ss_kaksi3                    string "ss_kaksi3"
	Ss_pross                     string "ss_pross"
	Ss3_dssp                     string "ss3_dssp"
	Ss3_stride                   string "ss3_stride"
	Ss3_kaksi3                   string "ss3_kaksi3"
	Ss3_pross                    string "ss3_pross"
	Ss3_cons_dssp_stride         string "ss3_cons_dssp_stride"
	Ss3_cons_dssp_pross          string "ss3_cons_dssp_pross"
	Ss3_cons_stride_kaksi3       string "ss3_cons_stride_kaksi3"
	Ss3_cons_stride_pross        string "ss3_cons_stride_pross"
	Ss3_cons_kaksi3_pross        string "ss3_cons_kaksi3_pross"
	Ss3_cons_dssp_stride_kaksi3  string "ss3_cons_dssp_stride_kaksi3"
	Ss3_cons_dssp_stride_pross   string "ss3_cons_dssp_stride_pross"
	Ss3_cons_dssp_kaksi3_pross   string "ss3_cons_dssp_kaksi3_pross"
	Ss3_cons_stride_kaksi3_pross string "ss3_cons_stride_kaksi3_pross"
	Ss3_cons_all                 string "ss3_cons_all"
}

type ProtdbConfig struct {
	Ip         string `toml:"db-ip"`
	Name       string `toml:"db-name"`
	Collection string `toml:"collection-name"`
	Init       string `toml:"init"`
	Target     string `toml:"target"`
}

func loadProteinsFromMongo(ip string, db string, collection string) []Protein {
	session, err := mgo.Dial(ip)
	if err != nil {
		fmt.Println("Can't connect to the database at", ip)
		panic(err)
	}
	defer session.Close()
	c := session.DB(db).C(collection)

	var result []Protein
	err = c.Find(nil).Iter().All(&result)
	if err != nil {
		panic(err)
	}
	// for i:=0; i<len(result); i++ {
	// 	fmt.Println("ConsAll:", result[i].Chains[0].ss3_cons_all)
	// }
	return result
}

func (c *Chain) getField(field string) string {
	r := reflect.ValueOf(c)
	s := reflect.Indirect(r).FieldByName(field)
	return s.String()
}

func GetProteins(db ProtdbConfig) (id, start, end string, e error) {
	proteins := loadProteinsFromMongo(db.Ip, db.Name, db.Collection)
	id = "all"
	start = "#"
	end = "#"
	for i := 0; i < len(proteins); i++ {
		start += proteins[i].Chains[0].getField(strings.Title(db.Init)) + "#"
		end += proteins[i].Chains[0].getField(strings.Title(db.Target)) + "#"
	}
	if len(start) != len(end) {
		e = fmt.Errorf("Error: Number of CA start cells is different from end cells")
	}
	return id, start, end, e
}
