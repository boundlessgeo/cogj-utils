package util

import (
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
	"log"
	"github.com/paulmach/orb"
	"os"
	"bufio"
)


func ToCoj(infile string, outfile string, tiles int)(bool, error) {


	//read geojson in
	rawbytes, err := ioutil.ReadFile(infile)
	if err != nil {
		log.Fatalf("Input file does not exist: %v", infile)
	}

	//unmarshall features
	gj, err := geojson.UnmarshalFeatureCollection(rawbytes)
	if err != nil {
		log.Fatal("Error reading GeoJson from file %s", infile)
	}

	b := gj.Features[0].Geometry.Bound()
	//calculate bbox
	for _, feat:= range gj.Features{
	 	b = b.Extend(feat.Geometry.Bound().Max)
		b = b.Extend(feat.Geometry.Bound().Min)
	}
	log.Printf("Feature BBOX is %v", b)

	bounds := SplitBounds(tiles,b)

	coj := NewCoj(bounds)
	coj.AddFeatures(gj.Features)

	if outfile == "" {
		outfile = infile+".coj"
	}
	f, err := os.Create(outfile)
	writer := bufio.NewWriter(f)
	coj.Write(writer)
	writer.Flush()
	return true, nil

}

func DebugBounds(bounds []orb.Bound){

	fc := geojson.FeatureCollection{}
	for _,b := range bounds{

		p := b.ToPolygon()
		f := geojson.Feature{Geometry:p}
		fc.Append(&f)
	}
	raw, _ := fc.MarshalJSON()
	log.Printf(string(raw))

}

//splits a bounding box into equal sized pieces
func SplitBounds(pieces int, bound orb.Bound) ([]orb.Bound){

	dx := (bound.Left() - bound.Right()) / float64(pieces)
	dy := (bound.Top() - bound.Bottom()) / float64(pieces)

	bounds := make([]orb.Bound,(pieces*pieces))

	cnt := 0
	for x := 0; x < pieces; x++{
		for y := 0; y < pieces; y++{

			temp := orb.Bound{
				orb.Point{bound.Right()+(dx*float64(x)), bound.Bottom()+(dy * float64(y))},
				orb.Point{bound.Right()+(dx*float64(x+1)), bound.Bottom()+(dy*float64(y+1))},
			}
			//todo something is fucked up with this -- shouldn't need to clean the bbox
			bounds[cnt] = temp.ToPolygon().Bound()
			cnt++

		}
	}

	return bounds
}


