package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/bits"
	"os"

	"image/png"
)

var Gray = color.Gray16{0x8888}

func main() {
	widths := []int{6, 8, 10, 12, 14, 16, 20, 24, 28, 32}
	fontFile, err := os.Open("0T5UIC1.HZK")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("")
	for size, width := range widths {
		height := 2 * width
		var numBytes int
		if width%8 > 0 {
			numBytes = ((width / 8) + 1)
		} else {
			numBytes = (width / 8)
		}
		picWidth := (width + 1) * 16
		picHeight := (height + 1) * 8
		charImage := image.NewRGBA(image.Rect(0, 0, picWidth, picHeight))
		charBytes := make([]byte, numBytes)
		for imgRow := 0; imgRow < 8; imgRow++ {
			for imgColumn := 0; imgColumn < 16; imgColumn++ {
				if size == 9 && imgRow == 7 && imgColumn == 15 {
					continue
				}
				for i := 0; i <= width; i++ {
					charImage.Set((imgColumn*(width+1))+i, imgRow*(height+1), Gray)
				}
				for vert := 1; vert <= height; vert++ {
					_, err := fontFile.Read(charBytes)
					if err != nil {
						log.Fatal(err)
					}
					bits := bitsToBits(charBytes)
					charImage.Set((imgColumn * (width + 1)), imgRow*(height+1)+vert, Gray)
					for horz := 1; horz <= width; horz++ {
						if bits[horz-1] == 1 {
							charImage.Set((imgColumn*(width+1))+horz, imgRow*(height+1)+vert, color.Black)
						} else {
							charImage.Set((imgColumn*(width+1))+horz, imgRow*(height+1)+vert, color.White)
						}
					}
				}
			}
		}
		for i := 0; i < picWidth; i++ {
			charImage.Set(i, picHeight-1, Gray)
		}
		for j := 0; j < picHeight; j++ {
			charImage.Set(picWidth-1, j, Gray)
		}
		imgFile, err := os.Create(fmt.Sprintf("0x%02d_%dx%d_0-127.png", size, width, height))
		if err != nil {
			log.Fatal(err)
		}
		png.Encode(imgFile, charImage)
		imgFile.Close()
	}
	// not sure what this chunk is
	mid := make([]byte, 5888)
	_, err = fontFile.Read(mid)
	if err != nil {
		log.Fatal(err)
	}
	for uni := 161; uni < 255; uni++ {
		eofError := false
		width := 16
		height := 16
		var numBytes int
		if width%8 > 0 {
			numBytes = ((width / 8) + 1)
		} else {
			numBytes = (width / 8)
		}
		picWidth := (width + 1) * 16
		picHeight := (height + 1) * 6
		charImage := image.NewRGBA(image.Rect(0, 0, picWidth, picHeight))
		charBytes := make([]byte, numBytes)
		for imgRow := 0; imgRow < 6; imgRow++ {
			for imgColumn := 0; imgColumn < 16; imgColumn++ {
				for i := 0; i <= width; i++ {
					charImage.Set((imgColumn*(width+1))+i, imgRow*(height+1), Gray)
				}
				if (imgRow == 0 && imgColumn == 0) || (imgRow == 5 && imgColumn == 15) {
					//unused squares, fill in
					for x := 1; x <= 16; x++ {
						for y := 1; y <= 16; y++ {
							charImage.Set((imgColumn*(width+1))+x, imgRow*(height+1)+y, Gray)
						}
					}
				} else {
					for vert := 1; vert <= height; vert++ {
						_, err := fontFile.Read(charBytes)
						if err != nil {
							fmt.Println("EOF")
							eofError = true
							break
						}
						bits := bitsToBits(charBytes)
						charImage.Set((imgColumn * (width + 1)), imgRow*(height+1)+vert, Gray)
						for horz := 1; horz <= width; horz++ {
							if bits[horz-1] == 1 {
								charImage.Set((imgColumn*(width+1))+horz, imgRow*(height+1)+vert, color.Black)
							} else {
								charImage.Set((imgColumn*(width+1))+horz, imgRow*(height+1)+vert, color.White)
							}
						}
					}
				}
				if eofError {
					break
				}
			}

			if eofError {
				break
			}
		}
		for i := 0; i < picWidth; i++ {
			charImage.Set(i, picHeight-1, Gray)
		}
		for j := 0; j < picHeight; j++ {
			charImage.Set(picWidth-1, j, Gray)
		}
		imgFile, err := os.Create(fmt.Sprintf("%dx%d_%X.png", width, height, uni))
		if err != nil {
			log.Fatal(err)
		}
		png.Encode(imgFile, charImage)
		imgFile.Close()

		if eofError {
			break
		}
	}
}

func bitsToBits(data []byte) (st []int) {
	st = make([]int, len(data)*8) // Performance x 2 as no append occurs.
	for i, d := range data {
		for j := 0; j < 8; j++ {
			if bits.LeadingZeros8(d) == 0 {
				// No leading 0 means that it is a 1
				st[i*8+j] = 1
			} else {
				st[i*8+j] = 0
			}
			d = d << 1
		}
	}
	return
}
