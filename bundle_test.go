package bundle

import (
	"bytes"
	"testing"
)

const fineMsg = "\tEverything went fine, \\ʕ◔ϖ◔ʔ/ YAY!"
const assetName = "example"

var assetData = []byte{3, 14, 15, 92, 65}
var compressedData = []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 18, 255, 223, 204, 194, 200, 200, 230, 84, 154, 151, 146, 147, 202, 248, 191, 133, 129, 145, 135, 241, 127, 19, 3, 131, 232, 255, 70, 102, 70, 166, 255, 77, 12, 140, 140, 140, 44, 46, 137, 37, 137, 140, 92, 12, 12, 12, 34, 32, 5, 236, 169, 21, 137, 185, 5, 57, 169, 140, 172, 204, 124, 252, 49, 142, 12, 128, 0, 0, 0, 255, 255, 12, 250, 1, 82, 67, 0, 0, 0}

func TestAddAsset(t *testing.T) {
	t.Log("Testing AddAsset")
	{
		b := make(Bundle)
		b.AddAsset(assetName, assetData)
		howManyAssets := len(b)
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
		b := make(Bundle)
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
		b := make(Bundle)
		b.AddAsset(assetName, assetData)
		howManyAssets := len(b)
		if howManyAssets != 1 {
			t.Errorf("\tPrevious step: expected howManyAssets=1, received %d", howManyAssets)
		} else {
			b.DeleteAsset(assetName)
			howManyAssets = len(b)
			if howManyAssets != 0 {
				t.Errorf("\tExpected howManyAssets=0, received %d", howManyAssets)
			} else {
				t.Log(fineMsg)
			}
		}
	}
}

func TestCompress(t *testing.T) {
	t.Log("Testing Compress")
	{
		b := make(Bundle)
		b.AddAsset(assetName, assetData)
		data, err := b.Compress()
		if err != nil {
			t.Errorf("\tCompress failed, error: %s", err)
		} else if len(data) != 90 {
			t.Errorf("\tExpected len(data)=90, received %d", len(data))
		} else if !bytes.Equal(data, compressedData) {
			t.Error("\tExpected 1f8b08..., received something else")
		} else {
			t.Log(fineMsg)
		}
	}
}

func TestDecompress(t *testing.T) {
	t.Log("Testing Decompress")
	{
		b, err := Decompress(compressedData)
		if err != nil {
			t.Errorf("\tDecompress failed, error: %s", err)
		} else if len(b) != 1 {
			t.Errorf("\tExpected bundle to has one asset, received %d", len(b))
		} else if _, ok := b[assetName]; !ok {
			t.Errorf("\tExpected %s to exists, it doesn't", assetName)
		} else if !bytes.Equal(b[assetName].Data, assetData) {
			t.Error("\tAssets are not the same")
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
		} else if len(b) != 3 {
			t.Errorf("\tExpected bundle to has 3 assets, received %d", len(b))
		} else if _, ok := b["logo.png"]; !ok {
			t.Errorf("\tExpected %s to exists, it doesn't", "logo.png")
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
		} else if len(rawData) != 42835 {
			t.Errorf("\tExpected len(rawData) = 42835, received %d", len(rawData))
		} else if !bytes.Equal(rawData[:10], []byte{3, 14, 1, 1, 14, 150, 196, 62, 0, 0}) {
			t.Error("\tFiles don't look the same")
		} else {
			t.Log(fineMsg)
		}
	}
}
