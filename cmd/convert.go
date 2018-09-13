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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/boundlessgeo/coj/util"
)

var in string
var out string
var tiles int

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "converts geojson in coj",
	Long: `Takes a normal geojson file and converts it into COJ `,
	Run: func(cmd *cobra.Command, args []string) {
		util.ToCoj(in, out, tiles)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().IntVar(&tiles,"tiles", 4, "the sqrt number of tiles to break the geojson into, 4 will result in 16 tiles, for example")
	convertCmd.Flags().StringVar(&in,"in", "", "input file")
	convertCmd.Flags().StringVar(&out,"out", "", "output file")
	convertCmd.MarkFlagRequired("in")

}
