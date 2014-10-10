package ca

import (
    "math"
    //"fmt"
)

func fitness(predicted, real []byte) (c3, cL, cH, cS, q3, qL, qH, qS float64) {
    /* Funcao que calcula o fitness (energia) para ser utilizada na otimizacao das regras do automato.
    INPUT:
    vetores contendo a estrutura secundaria predita e real
    OUTPUT:
    c3 = media da correlacao cL, cH e cS
    cL = correlacao da predicao de regioes de loops
    cH = correlacao da predicao de regioes de helice
    cS = correlacao da predicao de regioes de fitas
    ATENCAO:
    Aparentemente proteinas menores sao favorecidas neste calculo!!! */

    //as correlacoes somente devem ser calculadas se estiverem presentes na estrutura secundaria da proteina
    //a principio, considerasse que nenhuma delas esteja presente
    calculatecL := false
    calculatecH := false
    calculatecS := false

    //numero de true positives (tp), true negatives (tn), false positives (fp), false negatives (fn)
    cl_tp := 0
    cl_tn := 0
    cl_fp := 0
    cl_fn := 0

    ch_tp := 0
    ch_tn := 0
    ch_fp := 0
    ch_fn := 0

    cs_tp := 0
    cs_tn := 0
    cs_fp := 0
    cs_fn := 0

    //Q3, qL, qH e qS
    q3 = 0.0

    qL = 0.0
    npL := 0
    nrL := 0

    qH = 0.0
    npH := 0
    nrH := 0

    qS = 0.0
    npS := 0
    nrS := 0

    //percorre os vetores de estrutura secundaria fazendo a comparacao entre ambos
    //a posicao 0 e a ultima são nulas (#) e por isso a comparacao e feita do segundo ate o penultimo
    for i := 1; i < len(predicted) - 1; i++ {

        //se no vetor da estrutura secundaria real houver um loop na posicao i
        if real[i] == 1 {

            //possivel calcular cL
            calculatecL = true

            nrL += 1

            if predicted[i] == 1 {
                cl_tp += 1
                ch_tn += 1
                cs_tn += 1

                npL += 1

            } else if predicted[i] == 2 {
                cl_fn += 1
                ch_fp += 1
            } else if predicted[i] == 3 {
                cl_fn += 1
                cs_fp += 1
            } else {
                cl_fn += 1
            }

        //se no vetor da estrutura secundaria real houver um helice na posicao i
        } else if real[i] == 2 {

            //possivel calcular cH
            calculatecH = true

            nrH += 1

            if predicted[i] == 2 {
                ch_tp += 1
                cl_tn += 1
                cs_tn += 1

                npH += 1

            } else if predicted[i] == 1 {
                ch_fn += 1
                cl_fp += 1
            } else if predicted[i] == 3 {
                ch_fn += 1
                cs_fp += 1
            } else {
                ch_fn += 1
            }

        //se no vetor da estrutura secundaria real houver uma fita na posicao i
        } else if real[i] == 3 {

            //possivel calcular cS
            calculatecS = true

            nrS += 1

            if predicted[i] == 3 {
                cs_tp += 1
                cl_tn += 1
                ch_tn += 1

                npS += 1

            } else if predicted[i] == 1 {
                cs_fn += 1
                cl_fp += 1
            } else if predicted[i] == 2 {
                cs_fn += 1
                ch_fp += 1
            } else {
                cs_fn += 1
            }
        }
    }

    //A more rigorous measure (introduced to this field by [Matthews, 1975]) involves calculating the correlation coefficient for each target class
    //Calculo da correlacao cL, cH e cS. O calculo e feito em duas etapas e caso o divisor (denc*) seja Zero, a correlacao e considerada -1. Isto esta
    //correto?
    dencL := math.Sqrt(float64((cl_tp+cl_fp)*(cl_tp+ch_fn)*(cl_tn+cl_fp)*(cl_tn+cl_fn)))
    if dencL == 0.0 {
        cL = -1.0
    } else {
        cL = float64(cl_tp*cl_tn - cl_fp*cl_fn)/dencL
    }

    dencH := math.Sqrt(float64((ch_tp+ch_fp)*(ch_tp+ch_fn)*(ch_tn+ch_fp)*(ch_tn+ch_fn)))
    if dencH == 0.0 {
        cH = -1.0
    } else {
        cH = float64(ch_tp*ch_tn - ch_fp*ch_fn)/dencH
    }

    dencS := math.Sqrt(float64((cs_tp+cs_fp)*(cs_tp+cs_fn)*(cs_tn+cs_fp)*(cs_tn+cs_fn)))
    if dencS == 0.0 {
        cS = -1.0
    } else {
        cS = float64(cs_tp*cs_tn - cs_fp*cs_fn)/dencS
    }

    //calculo da c3 (media das correlacoes), apenas sao considerados para a media a correlacao de
    //estruturas secundarias presentes na estrutura secundaria real da proteina
    if calculatecL && calculatecH && calculatecS {
        c3 = (cL + cH + cS)/3.0
    } else if calculatecL && calculatecH && !calculatecS {
        c3 = (cL + cH)/2.0
    } else if calculatecL && !calculatecH && calculatecS {
        c3 = (cL + cS)/2.0
    } else if !calculatecL && calculatecH && calculatecS {
        c3 = (cH + cS)/2.0
    } else if calculatecL && !calculatecH && !calculatecS {
        c3 = cL
    } else if !calculatecL && calculatecH && !calculatecS {
        c3 = cH
    } else if !calculatecL && !calculatecH && calculatecS {
        c3 = cS
    }

    //Q3, ql, qh e qs
    q3 = float64((npL + npH + npS))/float64((nrL + nrH + nrS))
    //if float64((nrL + nrH + nrS) < 1 { q3 = 0.0}

    qL = float64(npL)/float64(nrL)
    qH = float64(npH)/float64(nrH)
    qS = float64(npS)/float64(nrS)

    return
}

func fitnessSimple(predicted, real []byte) (c3, q3 float64) {
    /* Funcao que calcula o fitness (energia) para ser utilizada na otimizacao das regras do automato.
    INPUT:
    vetores contendo a estrutura secundaria predita e real
    OUTPUT:
    c3 = media da correlacao cL, cH e cS
    cL = correlacao da predicao de regioes de loops
    cH = correlacao da predicao de regioes de helice
    cS = correlacao da predicao de regioes de fitas
    ATENCAO:
    Aparentemente proteinas menores sao favorecidas neste calculo!!! */

    //as correlacoes somente devem ser calculadas se estiverem presentes na estrutura secundaria da proteina
    //a principio, considerasse que nenhuma delas esteja presente
    calculatecL := false
    calculatecH := false
    calculatecS := false

    cL := 0.0
    cH := 0.0
    cS := 0.0

    //numero de true positives (tp), true negatives (tn), false positives (fp), false negatives (fn)
    cl_tp := 0
    cl_tn := 0
    cl_fp := 0
    cl_fn := 0

    ch_tp := 0
    ch_tn := 0
    ch_fp := 0
    ch_fn := 0

    cs_tp := 0
    cs_tn := 0
    cs_fp := 0
    cs_fn := 0

    //Q3, qL, qH e qS
    q3 = 0.0

    //qL := 0.0
    npL := 0
    nrL := 0

    //qH := 0.0
    npH := 0
    nrH := 0

    //qS := 0.0
    npS := 0
    nrS := 0

    //percorre os vetores de estrutura secundaria fazendo a comparacao entre ambos
    //a posicao 0 e a ultima são nulas (#) e por isso a comparacao e feita do segundo ate o penultimo
    for i := 1; i < len(predicted) - 1; i++ {

        //se no vetor da estrutura secundaria real houver um loop na posicao i
        if real[i] == 1 {

            //possivel calcular cL
            calculatecL = true

            nrL += 1

            if predicted[i] == 1 {
                cl_tp += 1
                ch_tn += 1
                cs_tn += 1

                npL += 1

            } else if predicted[i] == 2 {
                cl_fn += 1
                ch_fp += 1
            } else if predicted[i] == 3 {
                cl_fn += 1
                cs_fp += 1
            } else {
                cl_fn += 1
            }

        //se no vetor da estrutura secundaria real houver um helice na posicao i
        } else if real[i] == 2 {

            //possivel calcular cH
            calculatecH = true

            nrH += 1

            if predicted[i] == 2 {
                ch_tp += 1
                cl_tn += 1
                cs_tn += 1

                npH += 1

            } else if predicted[i] == 1 {
                ch_fn += 1
                cl_fp += 1
            } else if predicted[i] == 3 {
                ch_fn += 1
                cs_fp += 1
            } else {
                ch_fn += 1
            }

        //se no vetor da estrutura secundaria real houver uma fita na posicao i
        } else if real[i] == 3 {

            //possivel calcular cS
            calculatecS = true

            nrS += 1

            if predicted[i] == 3 {
                cs_tp += 1
                cl_tn += 1
                ch_tn += 1

                npS += 1

            } else if predicted[i] == 1 {
                cs_fn += 1
                cl_fp += 1
            } else if predicted[i] == 2 {
                cs_fn += 1
                ch_fp += 1
            } else {
                cs_fn += 1
            }
        }
    }

    //A more rigorous measure (introduced to this field by [Matthews, 1975]) involves calculating the correlation coefficient for each target class
    //Calculo da correlacao cL, cH e cS. O calculo e feito em duas etapas e caso o divisor (denc*) seja Zero, a correlacao e considerada -1. Isto esta
    //correto?
    dencL := math.Sqrt(float64((cl_tp+cl_fp)*(cl_tp+ch_fn)*(cl_tn+cl_fp)*(cl_tn+cl_fn)))
    if dencL == 0.0 {
        cL = -1.0
    } else {
        cL = float64(cl_tp*cl_tn - cl_fp*cl_fn)/dencL
    }

    dencH := math.Sqrt(float64((ch_tp+ch_fp)*(ch_tp+ch_fn)*(ch_tn+ch_fp)*(ch_tn+ch_fn)))
    if dencH == 0.0 {
        cH = -1.0
    } else {
        cH = float64(ch_tp*ch_tn - ch_fp*ch_fn)/dencH
    }

    dencS := math.Sqrt(float64((cs_tp+cs_fp)*(cs_tp+cs_fn)*(cs_tn+cs_fp)*(cs_tn+cs_fn)))
    if dencS == 0.0 {
        cS = -1.0
    } else {
        cS = float64(cs_tp*cs_tn - cs_fp*cs_fn)/dencS
    }

    //calculo da c3 (media das correlacoes), apenas sao considerados para a media a correlacao de
    //estruturas secundarias presentes na estrutura secundaria real da proteina
    if calculatecL && calculatecH && calculatecS {
        c3 = (cL + cH + cS)/3.0
    } else if calculatecL && calculatecH && !calculatecS {
        c3 = (cL + cH)/2.0
    } else if calculatecL && !calculatecH && calculatecS {
        c3 = (cL + cS)/2.0
    } else if !calculatecL && calculatecH && calculatecS {
        c3 = (cH + cS)/2.0
    } else if calculatecL && !calculatecH && !calculatecS {
        c3 = cL
    } else if !calculatecL && calculatecH && !calculatecS {
        c3 = cH
    } else if !calculatecL && !calculatecH && calculatecS {
        c3 = cS
    }

    //Q3, ql, qh e qs
   // fmt.Println("(npL + npH + npS))/float64((nrL + nrH + nrS)", npL, npH, npS,float64(nrL + nrH + nrS))
    q3 = float64((npL + npH + npS))/float64((nrL + nrH + nrS))
    //if float64((nrL + nrH + nrS) < 1 { q3 = 0.0}

    //qL = float64(npL)/float64(nrL)
    //qH = float64(npH)/float64(nrH)
    //qS = float64(npS)/float64(nrS)

    return
}
