package bundle

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"math"
	"io/ioutil"
)

const (
	NONE = iota
	COMPRESSED
	ENCRYPTED
	DEVELOPMENT
)

type Asset struct {
	Data []byte
}

type Bundle struct {
	Version byte
	Info byte
	Assets map[string]Asset
}

func (b Bundle) Asset(name string) ([]byte, error) {
	if _, ok := b.Assets[name]; !ok {
		return []byte{}, errors.New(fmt.Sprintf("asset %s does not exist", name))
	}
	if b.Info == DEVELOPMENT {
		return readFile(string(b.Assets[name].Data))
	}
	return b.Assets[name].Data, nil
}

func (b Bundle) AddAsset(name string, data []byte) {
	b.Assets[name] = Asset{
		Data: data,
	}
}

func (b Bundle) DeleteAsset(name string) error {
	if _, ok := b.Assets[name]; !ok {
		return errors.New(fmt.Sprintf("asset %s does not exist", name))
	}
	delete(b.Assets, name)
	return nil
}

func (b Bundle) Flat() (data []byte, err error) {
	encBuf := new(bytes.Buffer)
	err = gob.NewEncoder(encBuf).Encode(b)
	if err != nil {
		return
	}
	data = encBuf.Bytes()
	return
}

func Deflat(data []byte) (b Bundle, err error) {
	decBuf := new(bytes.Buffer)
	decBuf.Write(data)
	err = gob.NewDecoder(decBuf).Decode(&b)
	return
}

func Compress(data []byte) []byte {
	var c bytes.Buffer
	w := gzip.NewWriter(&c)
	w.Write(data)
	w.Close()
	return c.Bytes()
}

func Decompress(data []byte) (d []byte, err error) {
	encBuf := bytes.NewBuffer(data)
	r, err := gzip.NewReader(encBuf)

	if err != nil {
		return
	}

	d, err = ioutil.ReadAll(r)
	return
}

func LoadBundle(filename string, key []byte) (b Bundle, err error) {
	rawData, err := readFile(filename)

	if err != nil {
		return
	}

	if rawData[0] != 3 || rawData[1] != 14 {
		err = errors.New("Not a bundle file")
		return
	}

	//version
	if rawData[2] != 1 {
		err = errors.New("Wrong version")
		return
	}


	checkSum := binary.BigEndian.Uint32(rawData[4:8])
	dataLength := binary.BigEndian.Uint64(rawData[8:16])
	data := rawData[16:]

	if uint64(len(data)) != dataLength {
		err = errors.New("Wrong data size, file may be incompleted")
		return
	}

	if crc32.Checksum(data, crc32.IEEETable) != checkSum {
		err = errors.New("Wrong checksum, file may be corrupted")
		return
	}

	// encrypted
	if rawData[3] == ENCRYPTED {
		key = PadAESKey(key)
		data, err = AES256GCMDecrypt(data[12:], key, data[:12])
		if err != nil {
			return
		}

		data, err = Decompress(data)
		if err != nil {
			return
		}
	} else if rawData[3] == COMPRESSED {
		data, err = Decompress(data)
		if err != nil {
			return
		}
	}

	b, err = Deflat(data)

	b.Version = rawData[2]
	b.Info = rawData[3]
	return b, err
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

func AES256GCMEncrypt(plainText, key []byte) ([]byte, []byte) {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	cipheredText := aesgcm.Seal(nil, nonce, plainText, nil)
	return cipheredText, nonce
}

func AES256GCMDecrypt(cipheredText, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	plainText, err := aesgcm.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return []byte{}, err
	}

	return plainText, nil
}

func PadAESKey(key []byte) []byte {
	l := len(key)
	if l == 0 {
		return make([]byte, 16)
	}
	mod := int(math.Mod(float64(l), 16))
	d := (l+16-mod)/16;
	padded := make([]byte, d*16)
	i := 0
	k := 0
	for ;i<(d*16); {
		for j:=0;j<l;j++ {
			i = (k*l)+j
			if i>=(d*16) {
				break
			}
			padded[i] = key[j]
		}
		k++
	}
	return padded
}