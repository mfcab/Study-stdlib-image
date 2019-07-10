package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

type Canvas struct{
	Context *image.NRGBA
	CurrentFillColor color.NRGBA
	CurrentStrokeColor color.NRGBA
	CurrentStrokeWidth int
	PathStack  [][]image.Point
	CircleStack []circle

}


type polygon struct {
	plist []image.Point
}

func (p *polygon) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *polygon) Bounds() image.Rectangle {
	maxPoint:=getMaxPoint(p.plist)
	minPoint:=getMinPoint(p.plist)
	return image.Rect(minPoint.X,minPoint.Y,maxPoint.X,maxPoint.Y)
}

func (p *polygon) At(x, y int) color.Color {

	if IsInside(image.Point{x,y},p.plist) {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}


type line struct {
	plist []image.Point
	lineWidth int
}

func (l *line) ColorModel() color.Model {
	return color.AlphaModel
}

func (l *line) Bounds() image.Rectangle {
	maxPoint:=getMaxPoint(l.plist)
	minPoint:=getMinPoint(l.plist)
	return image.Rect(minPoint.X,minPoint.Y,maxPoint.X+1,maxPoint.Y+1)
}

func (l *line) At(x, y int) color.Color {

	if OnLine(image.Point{x,y},l.plist,l.lineWidth) {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

type circle struct{
	p image.Point
	r int
	startAngle float64
	endAngle   float64
}

func (a *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (a *circle ) Bounds() image.Rectangle {
	return image.Rect(a.p.X-a.r, a.p.Y-a.r, a.p.X+a.r, a.p.Y+a.r)
}

func (a *circle ) At(x, y int) color.Color {

	if a.onArc(image.Point{x,y}) {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}


func New(height,width int) Canvas{
	m:=image.NewNRGBA(image.Rect(0,0,height,width))
	draw.Draw(m,m.Bounds(),&image.Uniform{C:color.NRGBA{
		R:255,
		G:255,
		B:255,
		A:255,}},image.ZP,draw.Src)
	return Canvas{
		Context:m,
		CurrentFillColor:color.NRGBA{0,0,0,255},
		CurrentStrokeColor:color.NRGBA{0,0,0,255},
		CurrentStrokeWidth:1,
	}


}

func (c *Canvas) FillStyle(r,g,b,a float64){
	c.CurrentFillColor = color.NRGBA{
		R:uint8(r),
		G:uint8(g),
		B:uint8(b),
		A:uint8(a * 255),
	}

}

func (c *Canvas) FillRect(sh,sw,h,w int){
	d:=c.Context
	draw.Draw(d,image.Rectangle{
		Min:image.Pt(sh,sw),
		Max:image.Pt(sh+h,sw+w),
	},&image.Uniform{C:c.CurrentFillColor},image.ZP,draw.Over)
}

func (c *Canvas) Draw(path string) error{
	fd,err:=os.Create(path)
	if err!=nil{
		return err
	}

	err=jpeg.Encode(fd,c.Context,nil)
	return err
}


func (c *Canvas) StrokeRect(sh,sw,h,w int){
	d:=c.Context
	draw.Draw(d,image.Rectangle{
		Min:image.Pt(sh,sw),
		Max:image.Pt(sh+h,sw+w),
	},&image.Uniform{C:c.CurrentFillColor},image.ZP,draw.Over)
	draw.Draw(d,image.Rectangle{
		Min:image.Pt(sh+c.CurrentStrokeWidth,sw+c.CurrentStrokeWidth),
		Max:image.Pt(sh+h-c.CurrentStrokeWidth,sw+w-c.CurrentStrokeWidth),
	},&image.Uniform{C:color.White},image.ZP,draw.Src)

}


func (c *Canvas) ClearRect(sh,sw,h,w int){
	d:=c.Context
	draw.Draw(d,image.Rectangle{
		Min:image.Pt(sh,sw),
		Max:image.Pt(sh+h,sw+w),
	},&image.Uniform{C:color.White},image.ZP,draw.Src)


}


func (c *Canvas) BeginPath(){
 c.PathStack=make([][]image.Point,0)
}


func (c *Canvas) MoveTo(x,y int){
	c.PathStack=append(c.PathStack,[]image.Point{{x,y}})

}

func (c *Canvas) LineTo(x,y int){
	c.PathStack[len(c.PathStack)-1]=append(c.PathStack[len(c.PathStack)-1],image.Point{x,y})
}

func (c *Canvas)  Arc(x,y,r int,start,end float64){
	c.CircleStack=append(c.CircleStack,circle{image.Point{x,y},r,start,end})
}

func IsInside(p image.Point,plist []image.Point) bool{
	pointNum:=len(plist)
	if pointNum<3{
		return false
	}
	maxPoint:=getMaxPoint(plist)
	MinPoint:=getMinPoint(plist)
	if p.X < MinPoint.X|| p.X > maxPoint.X || p.Y < MinPoint.Y || p.Y > maxPoint.Y {
		return false
	}
	j:=pointNum-1
	var c =false
	for i:=0;i<pointNum;i++{
		if (plist[i].Y>p.Y)!=(plist[j].Y>p.Y){
		if p.X <((plist[j].X-plist[i].X) * (p.Y-plist[i].Y)/(plist[j].Y-plist[i].Y)+plist[i].X){
			c=!c
		}
	}
		j=i
	}
	return c
}


func getMinPoint(plist []image.Point) image.Point{
	var p=image.Point{X:10000,Y:10000}
	for _,v:=range plist{
		if v.X<p.X{
			p.X=v.X
		}
		if v.Y<p.Y{
			p.Y=v.Y
		}

	}
	return p
}


func getMaxPoint(plist []image.Point) image.Point{
	var p=image.Point{X:0,Y:0}
	for _,v:=range plist{
		if v.X>p.X{
			p.X=v.X
		}
		if v.Y>p.Y{
			p.Y=v.Y
		}

	}
	return p

}



func OnLine(p image.Point,plist []image.Point,w int) bool{
	if len(plist)<2{
		return false
	}
	for i:=1;i<len(plist);i++{
		maxPoint:=getMaxPoint(plist[i-1:i+1])
		minPoint:=getMinPoint(plist[i-1:i+1])
		if p.X < minPoint.X|| p.X > maxPoint.X || p.Y < minPoint.Y || p.Y > maxPoint.Y {
			continue
		}
		if float64(plist[i].Y -plist[i-1].Y ) / float64(plist[i].X -plist[i-1].X)==float64(p.Y -plist[i-1].Y) / float64(p.X -plist[i-1].X){
			return true
		}

	}
	return false
}


func (c *Canvas) Full(){
	d:=c.Context
	p:=c.PathStack
	circles:=c.CircleStack
	for _,v:=range p{
		draw.DrawMask(d,(&polygon{v}).Bounds(),&image.Uniform{C:c.CurrentFillColor},image.ZP,&polygon{v},image.Point{v[0].X,v[0].Y},draw.Over)
	}
	for _,circle:=range circles{
		draw.DrawMask(c.Context,(&circle).Bounds(),&image.Uniform{C:c.CurrentFillColor},image.ZP,&circle,circle.Bounds().Min,draw.Over)

	}
	return
}


func (c *Canvas) Stroke(){
	d:=c.Context
	p:=c.PathStack
	arcs:=c.CircleStack
	for _,v:=range p{
		draw.DrawMask(d,(&line{v,1}).Bounds(),&image.Uniform{C:c.CurrentFillColor},image.ZP,&line{v,1},image.Point{v[0].X,v[0].Y},draw.Over)
	}
	for _,arc:=range arcs{
		arc2:=circle{}
		arc2=arc
		arc2.r=arc.r-1
		draw.DrawMask(c.Context,(&arc).Bounds(),&image.Uniform{C:c.CurrentFillColor},image.ZP,&arc,arc.Bounds().Min,draw.Over)
		draw.DrawMask(c.Context,(&arc2).Bounds(),&image.Uniform{C:color.White},image.ZP,&arc2,arc2.Bounds().Min,draw.Over)


	}
	return

}


func (a circle) onArc(p image.Point) bool{
	xx, yy, rr := float64(p.X-a.p.X), float64(a.p.Y-p.Y), float64(a.r)
	if xx*xx+yy*yy <rr*rr {
		if a.startAngle <= a.endAngle {
		if xx > 0 {
			angle := math.Pi - math.Asin(yy/math.Hypot(xx, yy))
			if angle < a.endAngle {
				if angle > a.startAngle {
					return true
				}
			}
			return false
		} else {
			angle := math.Asin(yy / math.Hypot(xx, yy))
			if yy > 0 {
				if angle < a.endAngle {
					if angle > a.startAngle {
						return true
					}
				}
				return false
			} else {
				if (angle + 2*math.Pi) < a.endAngle {
					if angle > a.startAngle {
						return true
					}
				}
				return false
			}
		}
	} else {
			a.startAngle, a.endAngle = a.endAngle, a.startAngle
			if xx > 0 {
				angle := math.Pi - math.Asin(yy/math.Hypot(xx, yy))
				if angle < a.endAngle {
					if angle > a.startAngle {
						return false
					}
				}
				return true
			} else {
				angle := math.Asin(yy / math.Hypot(xx, yy))
				if yy > 0 {
					if angle < a.endAngle {
						if angle > a.startAngle {
							return false
						}
					}
					return true
				} else {
					if (angle + 2*math.Pi) < a.endAngle {
						if angle > a.startAngle {
							return false
						}
					}
					return true
				}

			}
		}
	}
	return false
}


