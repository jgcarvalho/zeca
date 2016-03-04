package cass

import (
	"io/ioutil"
	"strings"

	"github.com/jgcarvalho/zeca/proteindb"
)

//estrutura para armazenar informacao de um automato celular
type CellAuto struct {
	id     string
	seq    []byte
	cell   []byte
	trueSS []byte
}

func CreateOne(fn string) *CellAuto {
	/* Funcao que cria um automato celular de acordo com o arquivo passado
	INPUT:
	Nome do arquivo
	OUTPUT:
	Automato celular com id (nome=pdb id), sequencia(celulas) e estrutura real */

	c := new(CellAuto)
	c.id = fn
	c.seq, c.cell, c.trueSS = loadFile(fn)

	/*DEBUG
	fmt.Println(c.cell)
	fmt.Println(c.trueSS)
	println(lines[0][0:5]) */
	return c
}

func CreateN(fns []string) []CellAuto {
	/* Funcao que cria N automatos celulares de acordo com um vetor contendo o nome dos arquivos
	INPUT:
	Slice com nome dos arquivos das proteinas
	OUTPUT:
	Slice de automatos celulares com id (nome=pdb id), sequencia(celulas) e estrutura real */

	//cria uma slice de automatos celulares com dimensao igual ao numero de arquivos de proteinas
	cas := make([]CellAuto, len(fns))

	//inicializa os automatos
	for i := 0; i < len(fns); i++ {
		cas[i].id = fns[i]
		cas[i].seq, cas[i].cell, cas[i].trueSS = loadFile(fns[i])
	}
	return cas
}

func CreateFromProteins(p []proteindb.Protein) []CellAuto {
	cas := make([]CellAuto, len(p))
	for i := 0; i < len(p); i++ {
		cas[i].id = p[i].Pdb_id
		cas[i].seq = encode(p[i].Chains[0].Seq_pdb)
		cas[i].cell = encode(p[i].Chains[0].Seq_pdb)
		cas[i].trueSS = encode(p[i].Chains[0].Ss3_cons_all)
	}
	return cas
}

func loadFile(fn string) (seq []byte, cell []byte, trueSS []byte) {
	content, err := ioutil.ReadFile("/home/jgcarvalho/sscago/data/" + fn)
	if err != nil {
		println("erro na leitura do arquivo", fn, err)
	}

	//considerando haver duas linhas, a primeira a seq e a segunda a ss
	lines := strings.Split(string(content), "\n")
	seq = encode(lines[0])
	cell = encode(lines[0])
	trueSS = encode(lines[1])
	return
}

func encode(s string) []byte {
	cell := make([]byte, len(s))
	for i, v := range s {
		switch v {
		case '#':
			cell[i] = 0
		case '-':
			cell[i] = 1
		case '*':
			cell[i] = 2
		case '|':
			cell[i] = 3
		case 'A':
			cell[i] = 4
		case 'C':
			cell[i] = 5
		case 'D':
			cell[i] = 6
		case 'E':
			cell[i] = 7
		case 'F':
			cell[i] = 8
		case 'G':
			cell[i] = 9
		case 'H':
			cell[i] = 10
		case 'I':
			cell[i] = 11
		case 'K':
			cell[i] = 12
		case 'L':
			cell[i] = 13
		case 'M':
			cell[i] = 14
		case 'N':
			cell[i] = 15
		case 'P':
			cell[i] = 16
		case 'Q':
			cell[i] = 17
		case 'R':
			cell[i] = 18
		case 'S':
			cell[i] = 19
		case 'T':
			cell[i] = 20
		case 'V':
			cell[i] = 21
		case 'W':
			cell[i] = 22
		case 'Y':
			cell[i] = 23
		case '?':
			cell[i] = 24
		}
	}
	return cell
}

func decode(cell []byte) string {
	s := make([]byte, len(cell))
	for i, v := range cell {
		switch v {
		case 0:
			s[i] = '#'
		case 1:
			s[i] = '-'
		case 2:
			s[i] = '*'
		case 3:
			s[i] = '|'
		case 4:
			s[i] = 'A'
		case 5:
			s[i] = 'C'
		case 6:
			s[i] = 'D'
		case 7:
			s[i] = 'E'
		case 8:
			s[i] = 'F'
		case 9:
			s[i] = 'G'
		case 10:
			s[i] = 'H'
		case 11:
			s[i] = 'I'
		case 12:
			s[i] = 'K'
		case 13:
			s[i] = 'L'
		case 14:
			s[i] = 'M'
		case 15:
			s[i] = 'N'
		case 16:
			s[i] = 'P'
		case 17:
			s[i] = 'Q'
		case 18:
			s[i] = 'R'
		case 19:
			s[i] = 'S'
		case 20:
			s[i] = 'T'
		case 21:
			s[i] = 'V'
		case 22:
			s[i] = 'W'
		case 23:
			s[i] = 'Y'
		case 24:
			s[i] = '?'
		}
	}
	return string(s)
}
