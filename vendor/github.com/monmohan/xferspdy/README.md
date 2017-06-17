# xferspdy

This library provides binary diff and patch API in golang. 


#### Supported today:
* Command line utilities to diff and patch binary files
* Library for fingerprint generation, rolling hash and block matching
* *NEW* : For large files, fingerprint generation automatically switches to parallel mode wherein multiple go routines are used concurrently. For files > 20 MB, there is an improvement of ~50% compared to sequential fingerprint generation. 

Reference :
[Rsync Algorithm] (https://rsync.samba.org/tech_report/node2.html)

### Setup
* Need to have go installed, [golang downloads] (https://golang.org/dl/)
* Do go get

       `go get github.com/monmohan/xferspdy`

* Install the command line utilities

  Run  `go install ./...` from the _xferspdy_ directory

### Using the API
See [GoDoc](https://godoc.org/github.com/monmohan/xferspdy). 
The docs also contain an example usage of the API.
 
### Using the fpgen, diff and patch CLI utilities:
The library also provides CLI wrappers on API.

#### Example
* You can see the usage of any of these commands using `$ GOPATH/bin/<command> --help`

* Lets say you have a binary file  (e.g. power point presentation _MyPrezVersion1.pptx_).

* First generate a fingerprint of version 1

  `$ GOPATH/bin/fpgen -file <path>/MyPrezVersion1.pptx`
  
  This will generate the fingerprint file _<path>/MyPrezVersion1.pptx.fingerprint_.
* Lets say that the file was changed now (for example add a slide or image) and saved as _MyPrezVersion2.pptx_
* Now Generate a diff (*doesn't require original file*)

   `$ GOPATH/bin/diff -fingerprint <path>/MyPrezVersion1.pptx.fingerprint -file <path>/MyPrezVersion2.pptx`

 It will create a patch file `<path>/MyPrezVersion2.pptx.patch`

* Now patch the Version 1 file to get the Version 2
 
   `$ GOPATH/bin/patch -patch <path>/MyPrezVersion2.pptx.patch -base <path>/MyPrezVersion1.pptx`

* This will generate _<path>/Patched_MyPrezVersion1.pptx_. 
 *This file would exactly be same as MyPrezVersion2.pptx.*

*NOTE:* diff and patch are also common utilities present on most distributions so its better to give explicit path to these binaries. for example use _$GOPATH/bin/diff_ and _$GOPATH/bin/patch_


