# go-license [![Build Status](https://travis-ci.org/ryanuber/go-license.svg)](https://travis-ci.org/ryanuber/go-license)

A license management utility for programs written in Golang.

This program handles identifying software licenses and standardizing on a short,
abbreviated name for each known license type.

## Enforcement

License identifier enforcement is not strict. This makes it possible to warn
when an unrecognized license type is used, encouraging either conformance or an
update to the list of known licenses. There is no way we can know all types of
licenses.

## License guessing

This program also provides naive license guessing based on the license body
(text). This makes it easy to just throw a blob of text in and get a
standardized license identifier string out.

It is also possible to have `go-license` guess the file name that contains the
license data. This is done by scanning a directory for well-known license file
names.

## Recognized License Types

`MIT`<br>
The MIT license. ([text](http://opensource.org/licenses/MIT))

`NewBSD`<br>
The "new" or "revised" BSD license.
([text](http://opensource.org/licenses/BSD-3-Clause))

`FreeBSD`<br>
The "simplified" BSD  license.
([text](http://opensource.org/licenses/BSD-2-Clause))

`Apache-2.0`<br>
Apache License, version 2.0 ([text](http://opensource.org/licenses/Apache-2.0))

`MPL-2.0`<br>
The Mozilla Public License v2.0 ([text](http://opensource.org/licenses/MPL-2.0))

`GPL-2.0`<br>
The GNU General Public License v2.0
([text](http://opensource.org/licenses/GPL-2.0))

`GPL-3.0`<br>
The GNU General Public License v3.0
([text](http://opensource.org/licenses/GPL-3.0))

`LGPL-2.1`<br>
GNU Library or "Lesser" General Public License v2.1
([text](http://opensource.org/licenses/LGPL-2.1))

`LGPL-3.0`<br>
GNU Library or "Lesser" General Public License v3.0
([text](http://opensource.org/licenses/LGPL-3.0))

`CDDL-1.0`<br>
Common Development and Distribution License v1.0
([text](http://opensource.org/licenses/CDDL-1.0))

`EPL-1.0`<br>
Eclipse Public License v1.0 ([text](http://opensource.org/licenses/EPL-1.0))

## Example

```go
package main

import (
    "fmt"
    "github.com/ryanuber/go-license"
)

func main() {
    // This case will work if there is a guessable license file in the
    // current working directory.
    l, err := license.NewFromDir(".")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(l.Type)

    // This case will do the exact same thing as above, but uses an explicitly
    // set license file name instead of searching for one.
    l, err = license.NewFromFile("./LICENSE")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(l.Type)

    // This case will work when the license type can be guessed based on text
    l = new(license.License)
    l.Text = "The MIT License (MIT)"
    if err := l.GuessType(); err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(l.Type)

    // This case will work in all cases. The license type and the license data
    // are both being set explicitly. This enables one to use any license.
    l = license.New("MyLicense", "My terms go here")
    fmt.Println(l.Type)

    // This call determines if the license in use is recognized by go-license.
    fmt.Println(l.Recognized())
}
```
