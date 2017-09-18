package bundle

import (
	"bytes"
	"testing"
)

const fineMsg = "\tEverything went fine, \\ʕ◔ϖ◔ʔ/ YAY!"
const assetName = "example"

var assetData = []byte{3, 14, 15, 92, 65}
var compressedData = []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 28, 204, 49, 10, 194, 48, 24, 197, 241, 247, 76, 76, 65, 68, 39, 241, 10, 78, 110, 30, 160, 226, 226, 5, 92, 212, 33, 98, 148, 66, 155, 150, 38, 130, 171, 181, 122, 237, 175, 164, 219, 227, 241, 227, 191, 147, 143, 34, 205, 254, 229, 239, 165, 163, 116, 160, 98, 118, 114, 109, 40, 106, 79, 3, 234, 163, 127, 212, 105, 152, 60, 4, 23, 3, 229, 15, 96, 35, 63, 77, 174, 43, 219, 156, 67, 108, 11, 255, 188, 222, 198, 192, 118, 68, 201, 112, 78, 233, 129, 149, 124, 21, 39, 210, 131, 164, 62, 216, 104, 57, 67, 122, 59, 197, 204, 189, 109, 213, 148, 142, 83, 181, 88, 94, 114, 96, 8, 0, 0, 255, 255, 170, 2, 67, 33, 139, 0, 0, 0}
var flattenData = []byte{53, 255, 129, 3, 1, 1, 6, 66, 117, 110, 100, 108, 101, 1, 255, 130, 0, 1, 3, 1, 7, 86, 101, 114, 115, 105, 111, 110, 1, 6, 0, 1, 4, 73, 110, 102, 111, 1, 6, 0, 1, 6, 65, 115, 115, 101, 116, 115, 1, 255, 134, 0, 0, 0, 40, 255, 133, 4, 1, 1, 23, 109, 97, 112, 91, 115, 116, 114, 105, 110, 103, 93, 98, 117, 110, 100, 108, 101, 46, 65, 115, 115, 101, 116, 1, 255, 134, 0, 1, 12, 1, 255, 132, 0, 0, 21, 255, 131, 3, 1, 2, 255, 132, 0, 1, 1, 1, 4, 68, 97, 116, 97, 1, 10, 0, 0, 0, 21, 255, 130, 3, 1, 7, 101, 120, 97, 109, 112, 108, 101, 1, 5, 3, 14, 15, 92, 65, 0, 0}
var aeskey = []byte{1, 2, 3, 4, 5, 11, 12, 13, 14, 15, 51, 52, 53, 54, 55, 87}

func TestAddAsset(t *testing.T) {
	t.Log("Testing AddAsset")
	{
		var b Bundle
		b.Assets = make(map[string]Asset)
		b.AddAsset(assetName, assetData)
		howManyAssets := len(b.Assets)
		if howManyAssets != 1 {
			t.Errorf("\tExpected howManyAssets=1, received %d", howManyAssets)
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestAsset(t *testing.T) {
	t.Log("Testing AddAsset")
	{
		var b Bundle
		b.Assets = make(map[string]Asset)
		b.AddAsset(assetName, assetData)
		data, err := b.Asset(assetName)
		if err != nil {
			t.Errorf("\tAsset doesn't exists: %s", err)
		} else if !bytes.Equal(data, assetData) {
			t.Errorf("\tExpected %v, received: %v", assetData, data)
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestDeleteAsset(t *testing.T) {
	t.Log("Testing DeleteAsset")
	{
		var b Bundle
		b.Assets = make(map[string]Asset)
		b.AddAsset(assetName, assetData)
		howManyAssets := len(b.Assets)
		if howManyAssets != 1 {
			t.Errorf("\tPrevious step: expected howManyAssets=1, received %d", howManyAssets)
		} else {
			b.DeleteAsset(assetName)
			howManyAssets = len(b.Assets)
			if howManyAssets != 0 {
				t.Errorf("\tExpected howManyAssets=0, received %d", howManyAssets)
			} else {
				t.Log(fineMsg)
			}
		}
	}
}

func TestFlat(t *testing.T) {
	t.Log("Testing Flat")
	{
		var b Bundle
		b.Assets = make(map[string]Asset)
		b.AddAsset(assetName, assetData)
		data, err := b.Flat()
		if err != nil {
			t.Errorf("\tFlat failed, error: %s", err)
		} else if len(data) != 139 {
			t.Errorf("\tExpected len(data)=139, received %d", len(data))
		} else if !bytes.Equal(data, flattenData) {
			t.Error("\tExpected [53 255 129 3 ..., received something else")
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestDeflat(t *testing.T) {
	t.Log("Testing Deflat")
	{
		b, err := Deflat(flattenData)
		if err != nil {
			t.Errorf("\tDeflat failed, error: %s", err)
		} else if len(b.Assets) != 1 {
			t.Errorf("\tExpected bundle to has one asset, received %d", len(b.Assets))
		} else if _, ok := b.Assets[assetName]; !ok {
			t.Errorf("\tExpected %s to exists, it doesn't", assetName)
		} else if !bytes.Equal(b.Assets[assetName].Data, assetData) {
			t.Error("\tAssets are not the same")
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestCompress(t *testing.T) {
	t.Log("Testing Compress")
	{
		data := Compress(flattenData)
		if len(data) != 149 {
			t.Errorf("\tExpected len(data)=149, received %d", len(data))
		} else if !bytes.Equal(data, compressedData) {
			t.Error("\tExpected [31 139 8 0..., received something else")
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestDecompress(t *testing.T) {
	t.Log("Testing Decompress")
	{
		data, err := Decompress(compressedData)
		if err != nil {
			t.Errorf("\tDecompress failed, error: %s", err)
		} else if len(data) != 139 {
			t.Errorf("\tExpected len(data)=139, received %d", len(data))
		} else if !bytes.Equal(data, flattenData) {
			t.Error("\tExpected [53 255 129 3..., received something else")
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestPadAESKey(t *testing.T) {
	t.Log("Testing PadAESKey")
	{
		data := PadAESKey([]byte{1, 2, 3, 4, 5})
		if len(data) != 16 {
			t.Errorf("\tExpected len(data)=16, received %d", len(data))
		} else if !bytes.Equal(data, []byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1}) {
			t.Errorf("\tExpected [1 2 3 4 5 1 2 3 4 5 1 2 3 4 5 1] received: %s", data)
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestAES256GCMEncryptDecrypt(t *testing.T) {
	t.Log("Testing AES256GCM Encrypt & Decrypt (at the same time, since the nonce is random)")
	{
		data, nonce := AES256GCMEncrypt(assetData, aeskey)
		decData, err := AES256GCMDecrypt(data, aeskey, nonce)
		if err != nil {
			t.Errorf("\tDecrypt failed, error: %s", err)
		} else if !bytes.Equal(decData, assetData) {
			t.Error("\tExpected %s received %s", assetData, decData)
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestLoadBundle(t *testing.T) {
	t.Log("Testing LoadBundle")
	{
		b, err := LoadBundle("tests/test.bundle", []byte{})
		if err != nil {
			t.Errorf("\tLoadBundle failed, error: %s", err)
		} else if len(b.Assets) != 2 {
			t.Errorf("\tExpected bundle to has 2 assets, received %d", len(b.Assets))
		} else if _, ok := b.Assets["rabbits.jpg"]; !ok {
			t.Errorf("\tExpected %s to exists, it doesn't", "rabbits.jpg")
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestReadFile(t *testing.T) {
	t.Log("Testing readFile")
	{
		rawData, err := readFile("tests/test.bundle")
		if err != nil {
			t.Errorf("\treadFile failed, error: %s", err)
		} else if len(rawData) != 44212 {
			t.Errorf("\tExpected len(rawData) = 44212, received %d", len(rawData))
		} else if !bytes.Equal(rawData[:10], []byte{3, 14, 1, 0, 33, 83, 216, 127, 0, 0}) {
			t.Errorf("\tFiles don't look the same, received %v", rawData[:10])
		} else {
			t.Log(fineMsg)
		}
	}
}
