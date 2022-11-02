package bmp

import (
	"errors"
	"fmt"

	"github.com/digvijaysingh13/imgpro/util"
)

// marker index is slice of len 2
// 0th index tells from which index reading is to be started
// 1st index tells the length of reading
type marker [2]uint

// contains bmp image parsing information after parsing
type Bmp struct {
	// whole image data in byte form
	Data *[]byte
	// first 14 byte of header
	header marker
	// index from where image bitmap of starts
	// this informaton is parsed from first 14 byte of header.
	dataOffset uint

	// next four bytes after header i.e. [14, 18), tells the info of DIB (Device Independent Bitmap) header info
	// DIB header has usually 40 bytes which will capture the all informations
	// after 40 bytes color pallets information start till dataOffset
	dibHeader marker
	// horizonatal width of bitmap in pixel
	width uint
	// vertical width of bitmap in pixel
	height uint
	// number of planes (=1)
	planes uint
	// Bitmap per pixel use to store palate information, possible values as follows:
	// 1 = monochrome palette. NumColors = 1
	// 4 = 4bit palletized. NumColors = 16
	// 8 = 8bit palletized. NumColors = 256
	// 16 = 16bit RGB. NumColors = 65536
	// 24 = 24bit RGB. NumColors = 16M
	bpp uint
	// type of compressions
	// 0 = BI_RGB   no compression
	// 1 = BI_RLE8 8bit RLE encoding
	// 2 = BI_RLE4 4bit RLE encoding
	compression uint
	// (compressed) size of image
	// set mImageSize = 0 if mCompression = 0
	imageSize uint
	// horizontal resolution: pixel per meter
	xPixelPerMtr uint
	// vertical resolution: pixel per meter
	yPixelPerMtr uint
	// number of actual color used
	colorUsed uint
	// number of important color (0=all)
	importantColor uint

	// all the data after 40 bytes in DIB header belongs to color pallets
	// color pallet range, this is important if mBpp is less than or equal to 8
	// hence this block is semi optional
	// when mBpp is 16, 24, 32 color value is calculate with individual value of blue, green and red
	colorPallet marker

	// color byte len, this reader is supported to modify if bitmap per pixel is 24 or 32
	colorLen uint
}

func NewBmp(data *[]byte) Bmp {
	// header important data parsing
	header := marker{0, 14}
	offset := (*data)[10:14]
	dataOffset := util.BytesToUnsignInt(&offset)
	// dib header important data parsing
	dibSizeBites := (*data)[14:18]
	dibSize := util.BytesToUnsignInt(&dibSizeBites)
	// [18,22) 4 bytes for width (horizontal pixel length)
	byts := (*data)[18:22]
	width := util.BytesToUnsignInt(&byts)
	// [22, 26) 4 bytes for height (vertical pixel length)
	byts = (*data)[22:26]
	height := util.BytesToUnsignInt(&byts)
	// [26, 28) 2 bytes for plane
	byts = (*data)[26:28]
	plane := util.BytesToUnsignInt(&byts)
	// [28, 30) 2 bytes for bpp (bitmap per pixel)
	byts = (*data)[28:30]
	bpp := util.BytesToUnsignInt(&byts)
	// [30, 34) 4 bytes for compression
	byts = (*data)[30:34]
	compression := util.BytesToUnsignInt(&byts)
	// [34, 38) 4 bytes for image size
	byts = (*data)[34:38]
	imageSize := util.BytesToUnsignInt(&byts)
	// [38, 42) 4 bytes for horizontal pixel per meter
	byts = (*data)[38:42]
	xPpm := util.BytesToUnsignInt(&byts)
	// [42, 46) 4 bytes for vertical pixel per meter
	byts = (*data)[42:46]
	yPpm := util.BytesToUnsignInt(&byts)
	// [46, 50) 4 bytes for color used
	byts = (*data)[46:50]
	colorUsed := util.BytesToUnsignInt(&byts)
	// [50, 54) 4 bytes for important color
	byts = (*data)[50:54]
	impColor := util.BytesToUnsignInt(&byts)

	// after parsing all data finally create Bmp and return
	return Bmp{
		Data:           data,
		header:         header,
		dataOffset:     dataOffset,
		dibHeader:      marker{14, dibSize},
		width:          width,
		height:         height,
		planes:         plane,
		bpp:            bpp,
		compression:    compression,
		imageSize:      imageSize,
		xPixelPerMtr:   xPpm,
		yPixelPerMtr:   yPpm,
		colorUsed:      colorUsed,
		importantColor: impColor,
		colorPallet:    marker{54, dibSize - 54},
		colorLen:       bpp / 8,
	}
}

// prints the header
// prints the DIB header
// prints the color pallets
func (b *Bmp) PrintHeader() {
	println("=========== Header ==========")
	first := b.header[0]
	second := b.header[0] + b.header[1]
	header := (*b.Data)[first:second]
	// 1st 2 bytes of header field are used to identify BMP and DIB file is 0x42 0x4D in hexadecimal,
	// same as BM in ASCII. It can follow possible values.* BM – Windows 3.1x, 95, NT, … etc.
	// * BA – OS/2 struct bitmap array * CI – OS/2 struct color icon * CP – OS/2 const color pointer
	// * IC – OS/2 struct icon * PT – OS/2 pointer
	identifier := string(header[0:2])
	// next 4 bytes size of bmp file in bytes.
	byt := header[2:6]
	size := util.BytesToUnsignInt(&byt)
	// next 2 byte is reserved for application that created the image.
	reserveFist := int(header[6:8][0])
	// next 2 byte is reserved for application that created the image.
	reserveSecond := int(header[8:10][0])
	// next 4 byte is offset, i.e. starting address of bytes where the bitmap image data (pixel array) can be found.
	// offset is already been parsed as dataOffset
	fmt.Printf("identified: %s\nsize: %d\nreserve_1: %d\nreserve_2:%d\ndata_offset: %d\n",
		identifier, size, reserveFist, reserveSecond, b.dataOffset)
	println("========== DIB (Device Independent Bitmap) Header ==========")
	fmt.Printf("width: %d\nheight: %d\nplane: %d\nbitmap per pixel: %d\ncompression: %d\ncompressed image size: %d\nhorizontal pixel per meter: %d\nvertical pixel per meter: %d\ncolor used: %d\nimportant color: %d\n",
		b.width, b.height, b.planes, b.bpp, b.compression, b.imageSize, b.xPixelPerMtr, b.yPixelPerMtr, b.colorUsed, b.importantColor)
	println("========== Color Pallets ==========")
	start := b.colorPallet[0]
	end := b.colorPallet[0] + b.colorPallet[1]
	for start < end {
		blue := (*b.Data)[start]
		green := (*b.Data)[start+1]
		red := (*b.Data)[start+2]
		fmt.Printf("red: %d, green: %d, blue: %d\n\n", red, green, blue)
		start += 3
	}
}

func (b *Bmp) Width() uint {
	return b.width
}

func (b *Bmp) Height() uint {
	return b.height
}

/**
 * if color width is 3 then array of 3 len RGB order is (blue green red)
 * if color width is 4 then array of 4 len RGBA order is
 * do not work for other case
 *
 * top left bpp data shown at the bottom left of image.
 * */
func (b *Bmp) pixelDataIndex(x, y uint) (uint, error) {
	if x > b.width || y > b.height {
		return 0, errors.New(fmt.Sprintf("x: %d, y: %d should not be greater than width: %d, height: %d respectively.", x, y, b.width, b.height))
	}
	row := b.height - y - 1
	dataStartIndex := b.dataOffset + (((row * b.width) + x) * b.colorLen)
	return dataStartIndex, nil
}

// it iterates through each pixels
func (b *Bmp) EachPixel(callback func(*[]byte)) error {
	for y := uint(0); y < b.height; y++ {
		for x := uint(0); x < b.width; x++ {
			index, error := b.pixelDataIndex(x, y)
			if error != nil {
				return error
			}
			slice := (*b.Data)[index : index+b.colorLen]
			callback(&slice)
		}
	}
	return nil
}

func (b *Bmp) Copy() *Bmp {
	original := (*b.Data)
	data := make([]byte, len(original))
	copy(data, original)
	bmp := NewBmp(&data)
	return &bmp
}

// removes the color
func MakeRed(bmp *Bmp) {
	// 0th index is blue, 1st is green, 2nd is red
	bmp.EachPixel(func(b *[]byte) {
		byt := *b
		byt[0] = 0
		byt[1] = 0
	})
}

func MakeGreen(bmp *Bmp) {
	// 0th index is blue, 1st is green, 2nd is red
	bmp.EachPixel(func(b *[]byte) {
		byt := *b
		byt[0] = 0
		byt[2] = 0
	})
}

func MakeBlue(bmp *Bmp) {
	// 0th index is blue, 1st is green, 2nd is red
	bmp.EachPixel(func(b *[]byte) {
		byt := *b
		byt[1] = 0
		byt[2] = 0
	})
}

func MakeGray(bmp *Bmp) {
	// 0th index is blue, 1st is green, 2nd is red
	bmp.EachPixel(func(b *[]byte) {
		byt := *b
		color := util.LuminousGrayscale(byt[2], byt[1], byt[0])
		byt[0] = color
		byt[1] = color
		byt[2] = color
	})
}
