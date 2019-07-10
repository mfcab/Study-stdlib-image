package main

import (
	"github.com/mfcab/canvas"
	"fmt"
	"math"
)

func main(){
	ctx:=canvas.New(400,400)
	ctx.BeginPath()
	ctx.MoveTo(0,0)
	ctx.LineTo(40,0)
	ctx.LineTo(0,40)
	ctx.LineTo(40,40)
	ctx.LineTo(40,80)
	ctx.LineTo(0,40)
	ctx.Arc(50,50,50,math.Pi/2,math.Pi)
	ctx.Full()
	err:=ctx.Draw("123.jpg")
	fmt.Println(err)
}

