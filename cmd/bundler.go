package main

import (
	"fmt"

	"io/ioutil"
	"log"

	"bufio"
	"encoding/binary"
	"hash/crc32"
	"os"

	"github.com/conejoninja/bundle"
	"github.com/spf13/cobra"
)

var output string

func main() {

	var cmdMake = &cobra.Command{
		Use:   "make [options] dir1 dir2 ... dirN",
		Short: "Make a new bundle from a folder",
		Long:  "Make a new bundle from a folder",
		Args:  cobra.MinimumNArgs(1),
		Run:   exec,
	}
	cmdMake.Flags().StringVarP(&output, "output", "o", "", "Output file name")

	var rootCmd = &cobra.Command{Use: "bundler"}
	rootCmd.AddCommand(cmdMake)
	rootCmd.Execute()
}

func exec(cmd *cobra.Command, args []string) {

	b := make(bundle.Bundle)
	for _, d := range args {
		readDir(d, &b)
	}
	a, _ := b.Compress()
	writeFile(output, a)
}

func readDir(dir string, b *bundle.Bundle) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			readDir(dir+"/"+f.Name(), b)
		} else {
			data, err := readFile(dir + "/" + f.Name())
			if err == nil {
				b.AddAsset(f.Name(), data)
			}
		}
		fmt.Println(dir + "/" + f.Name())
	}
}

func readFile(filename string) ([]byte, error) {
	var empty []byte

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return empty, err
	}

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return empty, statsErr
	}
	var size int64 = stats.Size()
	fw := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(fw)
	return fw, err
}

func writeFile(filename string, data []byte) {

	dataLength := len(data)
	dataLengthByte := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLengthByte, uint64(dataLength))

	checkSum := crc32.Checksum(data, crc32.IEEETable)
	checkSumByte := make([]byte, 4)
	binary.BigEndian.PutUint32(checkSumByte, checkSum)

	fullData := append([]byte{3, 14, 1, 1}, checkSumByte...)
	fullData = append(fullData, dataLengthByte...)
	fullData = append(fullData, data...)

	err := ioutil.WriteFile(filename, fullData, 0644)
	if err != nil {
		fmt.Println("Unable to create File", err)
		return
	}
}
