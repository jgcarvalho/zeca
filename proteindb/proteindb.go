package proteindb

import (
	// "bitbucket.org/zgcarvalho/zeca/ca"
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

func LoadProteinsFromMongo(ip string, db string, collection string) []Protein {
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
	// 	fmt.Println("ConsAll:", result[i].Chains[0].Ss3_cons_all)
	// }
	return result
}

// func (p *Protein) CreateCA1D(end string, r rules.Params) ca.CellAuto1D {
// 	switch end {
// 	case "dssp":
// 	case "stride":
// 	case "kaksi":
// 	case "pross":
// 	case "dssp+stride":
// 	case "dssp+pross":
// 	}

// 	c := ca.CreateCA1D(p.Pdb_id, begin, expectedEnd, rule)
// 	return cas
// }
