package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
* Output dir is imgpro-output on home dir
**/
func OutputDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	// check and create desktop dir
	outdirPath := filepath.Join(homeDir, "Desktop")
	_, err = os.Stat(outdirPath)
	if os.IsNotExist(err) {
		// case if can not create dir
		if err = os.Mkdir(outdirPath, os.ModePerm); err != nil {
			return "", err
		}

	}
	// check and create imgpro-output dir
	outdirPath = filepath.Join(outdirPath, "imgpro-output")
	if _, err = os.Stat(outdirPath); os.IsNotExist(err) {
		// in case if can not create dir
		if err = os.Mkdir(outdirPath, os.ModePerm); err != nil {
			return "", err
		}
	}
	return outdirPath, nil
}

/**
* creates file in imgpro-output dir
**/
func CreateFile(filename string) (*os.File, error) {
	outputDir, err := OutputDir()
	if err != nil {
		return nil, err
	}
	filename = filepath.Join(outputDir, filename)
	return os.Create(filename)
}

/*
* reads file and returns byte array
**/
func ReadFile(filePath string) (*[]byte, error) {
	out, err := ioutil.ReadFile(filePath)
	return &out, err
}

/*
* converts byte array into int
* if byte array is nil default value of int is returned
* transformation is coverted by assuming that supplied bytes are big endian
 */

func BytesToUnsignInt(byt *[]byte) uint {
	if byt == nil {
		return 0
	}
	var result uint = 0
	len := len(*byt)
	if len == 0 {
		return 0
	}
	index := 0
	for index < len {
		b := (*byt)[index]
		// changing negative sign to postive if any by making 2's compliment
		bInt := uint(b)
		bInt = bInt << (index * 8)
		result = result | bInt
		index += 1
	}
	return result
}

// equal importance is given to all RGB colors
// take RGB values and calculates it average and returns it
func AvgGrayscale(red, blue, green byte) byte {
	var avg int = (int(red) + int(blue) + int(green)) / 3
	return byte(avg)
}

// since human eyes are more sensitive to green, than red than blue
// 0.59 importance for green, 0.3 for red and 0.11 for blue
func LuminousGrayscale(red, blue, green byte) byte {
	var lim float32 = (0.3 * float32(red)) + (0.59 * float32(green)) + (0.11 * float32(blue))
	return byte(lim)
}

func WriteFile(filename string, data *[]byte) error {
	file, err := CreateFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(*data)
	return err
}
