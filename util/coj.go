// Copyright © 2018 Boundless
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package util

import (
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb"
	"log"
	"io"
	"encoding/json"
	"errors"
	"time"
)

type Coj struct {
	Collections []CojCollection
	Bbox orb.Bound
	featureCount int64
	size int64
}
type CojHeader struct {
	//size of this file in bytes
	Size int64 `json:"size"`
	//number of features contained in the whole file
	Features int64 `json:"features,omitempty"`

	Name string `json:"name,omitempty"`

	Published time.Time `json:"published,omitempty"`

	Version string `json:"version,omitempty"`

	//bbox of the entire file
	Bbox []float64 `json:"bbox"`
	//metadata about the collections
	Collections [] CojHeaderCollection `json:"collections,omitempty"`
}
//metadata about a given collection
type CojHeaderCollection struct {
	//the byte in the size where this tile starts
	Start int64 `json:"start"`
	//the size (length) of the tile
	Size int64 `json:"size"`
	//the bbox of the tile
	Bbox []float64 `json:"bbox"`
	//the number of features it contains
	Features int64 `json:"features"`
	//a mapping to the id of the given tile (internal use)
	tileId int
}

type CojCollection struct {
	Features *geojson.FeatureCollection
	Bound    orb.Bound
	id int
	size int64
	start int64
}

func NewCoj(bounds []orb.Bound) (Coj) {

	coj := Coj{}
	coj.Collections = make([]CojCollection, len(bounds))
	coj.Bbox = bounds[0]
	for i, b := range bounds {
		coj.Collections[i] = CojCollection{Bound: b, Features: &geojson.FeatureCollection{}, id: i}
		coj.Bbox = coj.Bbox.Extend(b.Max)
		coj.Bbox = coj.Bbox.Extend(b.Min)
	}
	return coj

}

func (c *Coj) AddFeatures(features []*geojson.Feature){

	for _, feat := range features {

		tb := feat.Geometry.Bound()

		for _, t := range c.Collections {

			if t.Bound.Contains(tb.Max) || t.Bound.Contains(tb.Min) {
				t.Features = t.Features.Append(feat)
				c.featureCount++
			}

		}
	}
}

func (c *Coj) Write(writer io.Writer)(bool, error){

	tileData := make([]byte,0)

	for i := 0; i < len(c.Collections); i++{
		//don't write bboxes with 0 features
		if len(c.Collections[i].Features.Features) == 0 {
			continue
		}
		data, err := c.Collections[i].Features.MarshalJSON()
		if err != nil{
			log.Printf("Error marshalling features! (%v)", err)
		}
		c.Collections[i].start = int64(len(tileData)) + 10240
		c.Collections[i].size = int64(len(data))-1
		tileData = append(tileData, data...)
	}
	c.size = int64(10240+len(tileData))
	header, err := c.CreateHeader()
	if err != nil{
		return false, errors.New("error marshalling header")
	}
	headerData, err := json.Marshal(header)
	if err != nil {
		return false, errors.New("error marshalling header")
	}

	log.Printf("Header is %v bytes", len(headerData))
	if len(headerData) < 10240{
		log.Printf("Expanding header to use first 10k")
		emptyData := make([]byte,10240-len(headerData))
		for i := 0; i < len(emptyData); i++{
			emptyData[i] = 0x0020
		}
		headerData = append(headerData, emptyData...)
		log.Printf("header is now %v bytes", len(headerData))
	}

	writer.Write(headerData)
	writer.Write(tileData)

	return true, nil



}
func (c *Coj) CreateHeader()(CojHeader,error){

	ch := CojHeader{Features: c.featureCount, Published: time.Now()}
	ch.Size = c.size
	ch.Bbox = toBbox(c.Bbox)
	ch.Collections = make([]CojHeaderCollection,0)

	for _, tile := range c.Collections{
		//don't write empty feature bboxes
		if len(tile.Features.Features) == 0{
			continue
		}
		cht := CojHeaderCollection{}
		cht.Bbox = toBbox(tile.Bound)
		cht.Features = int64(len(tile.Features.Features))
		cht.tileId = tile.id
		cht.Start = tile.start
		cht.Size = tile.size
		ch.Collections = append(ch.Collections,cht)
	}
	return ch, nil
}
func toBbox(bound orb.Bound)([]float64){

	return []float64{bound.Left(),bound.Bottom(), bound.Right(), bound.Top()}
}