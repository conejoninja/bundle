package main

import (
	"fmt"

	"io/ioutil"
	"log"

	"bufio"
	"encoding/binary"
	"hash/crc32"
	"os"

	"path/filepath"

	"github.com/conejoninja/bundle"
	"github.com/spf13/cobra"
	"strings"
)

var output string
var key string
var development bool
var bpath string

type asset struct {
	name string
	data []byte
	path string
}
var basepath string

func main() {

	var cmdMake = &cobra.Command{
		Use:   "make [options] dir1 dir2 ... dirN",
		Short: "Make a new bundle from a folder",
		Long:  "Make a new bundle from a folder",
		Args:  cobra.MinimumNArgs(1),
		Run:   exec,
	}
	cmdMake.Flags().StringVarP(&output, "output", "o", "", "Output file name")
	cmdMake.Flags().StringVarP(&key, "key", "k", "", "If set, bundle will be encrypted")
	cmdMake.Flags().BoolVarP(&development, "dev", "d", false, "Files will not be included in the bundle, but will be loaded from the disk")
	cmdMake.Flags().StringVarP(&bpath, "basepath", "b", "", "When in dev mode, real path will be replaced by this")

	cwd, err := os.Getwd()
	if err != nil && development {
		log.Fatal("Can not determine current directory")
	}
	basepath = cwd

	var rootCmd = &cobra.Command{Use: "bundler"}
	rootCmd.AddCommand(cmdMake)
	rootCmd.Execute()
}

func exec(cmd *cobra.Command, args []string) {

	assets := make(map[string]asset)
	for _, d := range args {
		if !filepath.IsAbs(d) {
			d = filepath.Join(basepath, d)
		}
		tmp := readDir(d)
		for k, v := range tmp {
			if development && bpath != "" {
				v.path = strings.Replace(v.path, d, bpath, 1)
			}
			assets[k] = v
		}

	}

	var b bundle.Bundle
	b.Assets = make(map[string]bundle.Asset)
	if development {
		for _, v := range assets {
			b.AddAsset(v.name, []byte(v.path))
		}
		flat, _ := b.Flat()
		writeFile(output, flat)
		return
	}

	for _, v := range assets {
		b.AddAsset(v.name, v.data)
	}

	enc := bundle.NONE
	if key != "" {
		enc = bundle.ENCRYPTED
	}
	if development {
		enc = bundle.DEVELOPMENT
	}
	b.Info = byte(enc)
	b.Version = 1

	flat, _ := b.Flat()

	var a []byte
	if enc == bundle.ENCRYPTED || enc == bundle.COMPRESSED {
		a = bundle.Compress(flat)
	} else {
		a = flat
	}

	if key != "" {
		keyByte := bundle.PadAESKey([]byte(key))
		data, nonce := bundle.AES256GCMEncrypt(a, keyByte)
		a = append(nonce, data...)

	}

	writeFile(output, a)
}

func readDir(dir string) map[string]asset {

	assets := make(map[string]asset)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			tmp := readDir(filepath.Join(dir, f.Name()))
			for k, v := range tmp {
				assets[k] = v
			}
		} else {
			data, err := readFile(filepath.Join(dir, f.Name()))
			if err == nil {
				assets[f.Name()] = asset{
					name: f.Name(),
					data: data,
					path: filepath.Join(dir, f.Name()),
				}
			}
		}
	}
	return assets
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

	enc := bundle.NONE
	if key != "" {
		enc = bundle.ENCRYPTED
	}
	if development {
		enc = bundle.DEVELOPMENT
	}

	fullData := []byte{3, 14, 1, byte(enc)}
	fullData = append(fullData, checkSumByte...)
	fullData = append(fullData, dataLengthByte...)
	fullData = append(fullData, data...)

	err := ioutil.WriteFile(filename, fullData, 0644)
	if err != nil {
		fmt.Println("Unable to create File", err)
		return
	}
}
