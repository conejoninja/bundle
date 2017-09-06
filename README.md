Bundle
======
Tool to pack different assets in a single *external* file and include it in your applications. This project is similar to [bindata](https://github.com/jteeuwen/go-bindata), [go.rice](https://github.com/GeertJohan/go.rice), [statik](https://github.com/rakyll/statik),..., but instead of include the assets in your executable and making it bigger, it creates a file that could be loaded by your application. 

## Usage (executable bundler)

The command *bundler make* will pack all the files in the given directories and create a g-zipped file of them
```bash
$ bundler make [options] dir1 dir2 ... dirN
``` 

## Usage (in your go application)

```go
b, err := bundle.LoadBundle("path/to/your/file.bundle")
if err != nil {
    // error handling
}
// access asset
data, _ := b.Asset("assetName")
fmt.Println(data)
```

Check the [examples folder](https://github.com/conejoninja/bundle/tree/master/examples)


## Contributing to this project:

If you find any improvement or issue you want to fix, feel free to send me a pull request.

## License

The MIT License (MIT)

Copyright 2017 Daniel Esteban - conejo@conejo.me

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE. 


There are several libraries used in the code, described in the section below that may have different licenses.




## Notes


![](https://raw.githubusercontent.com/conejoninja/cerrojo/master/assets/ribbon.png)

If you would like to donate via Bitcoin, please send your donation to this wallet:

   ![](https://raw.githubusercontent.com/conejoninja/cerrojo/master/assets/qr.png)

Bitcoin: **1G9d7uVvioNt8Emsv6fVmCdAPc41nX1c8J**