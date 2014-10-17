package plot

import (
	"fmt"
	"image/color"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
)

func Histogram(popFitness []float64, selFitness []float64, gen int) {
	popValues := make(plotter.Values, len(popFitness))
	selValues := make(plotter.Values, len(selFitness))
	for i := range popValues {
		popValues[i] = popFitness[i]
		fmt.Println("plot", i, popValues[i], popFitness[i])
	}
	for j := range selValues {
		selValues[j] = selFitness[j]
	}
	fmt.Println("PLOT", len(popFitness), len(popValues))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = fmt.Sprint("Generation ", gen)
	p.X.Label.Text = "Fitness"
	p.Y.Label.Text = "Number of individuals"

	hPop, err := plotter.NewHist(popValues, 10)
	if err != nil {
		panic(err)
	}
	hPop.Normalize(1)
	gray := color.RGBA{0, 0, 0, 64}
	hPop.FillColor = gray
	hSel, err := plotter.NewHist(selValues, 10)
	if err != nil {
		panic(err)
	}
	hSel.Normalize(1)
	blue := color.RGBA{0, 61, 245, 64}
	hSel.FillColor = blue

	p.X.Min, p.X.Max = 0.0, 1.0
	p.Add(hPop)
	p.Add(hSel)
	if err := p.Save(8, 4, fmt.Sprintf("hist_gen_%v.svg", gen)); err != nil {
		panic(err)
	}
}
