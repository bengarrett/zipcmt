# zipcmt

[![goreleaser](https://github.com/bengarrett/zipcmt/actions/workflows/release.yml/badge.svg)](https://github.com/bengarrett/zipcmt/actions/workflows/release.yml) &nbsp;
[![Go Reference](https://pkg.go.dev/badge/github.com/bengarrett/zipcmt.svg)](https://pkg.go.dev/github.com/bengarrett/zipcmt)

A zip archive comment batch viewer and extractor.

## Downloads

- [Windows](https://github.com/bengarrett/zipcmt/releases/latest/download/myip_Windows_Intel.zip)
- [macOS](https://github.com/bengarrett/zipcmt/releases/latest/download/myip_macOS_Intel.tar.gz
), [or for the Apple M chip](https://github.com/bengarrett/zipcmt/releases/latest/download/myip_macOS_M-series.tar.gz
)
- [FreeBSD](https://github.com/bengarrett/zipcmt/releases/latest/download/myip_FreeBSD_Intel.tar.gz
)
- [Linux](https://github.com/bengarrett/zipcmt/releases/latest/download/myip_Linux_Intel.tar.gz
)

### Packages

##### macOS [Homebrew](https://brew.sh/)
```sh
brew install bengarrett/homebrew-myip/zipcmt
```


[APK (Alpine package)](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.apk)
```sh
wget https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.apk
apk add zipcmt.apk
```

[DEB (Debian package)](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.deb)
```sh
wget https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.deb
dpkg -i zipcmt.deb
```

[RPM (Redhat package)](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.rpm)
```sh
wget https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.rpm
rpm -i zipcmt.rpm
```

##### Windows [Scoop](https://scoop.sh/)
```sh
scoop bucket add bengarrett https://github.com/bengarrett/zipcmt.git
scoop install bengarrett/zipcmt
```

## Usage

```sh
zipcmt -help

Usage:
    zipcmt [options] [directories]

Examples:
    zipcmt --print --nodupes .		# scan the working directory and only show unique comments
    zipcmt --export ~/Downloads		# scan the download directory and save all comments
    zipcmt -r -d=~/text ~/Downloads	# recursively scan the download directory and save all comments to a directory
    zipcmt -n -p -q -r / | less		# scan the whole system and view unique comments in a page reader

Options:
    -r, --recursive    recursively walk through all subdirectories while scanning for zip archives
    -n, --nodupes      no duplicate comments, only show unique finds
    -p, --print        print the comments to the terminal

    -e, --export       save the comment to a textfile stored alongside the archive (use at your own risk)
    -d, --exportdir    save the comment to a textfile stored in this directory
    -o, --overwrite    overwrite any previously exported comment textfiles

        --raw          use the original comment text encoding instead of Unicode

    -q, --quiet        suppress zipcmt feedback except for errors
    -v, --version      version and information for this program
    -h, --help         show this list of options
```

##### Print example
```sh
ls -l test/
# test-no-comment.zip  test-with-comment.zip  test.txt

zipcmt --print test/

#  ── test/test-with-comment.zip ───────────┐
#    This is an example test comment for zipcmt.
#
# Scanned 2 zip archives and found 1 comment
```

##### Quiet example
```sh
ls -l test/
# test-no-comment.zip  test-with-comment.zip  test.txt

zipcmt --print --quiet test/
#   This is an example test comment for zipcmt.
```
##### No dupes example
```sh
cp test/test-with-comment.zip test/test-with-comment-1.zip

ls -l test/
# test-no-comment.zip  test-with-comment-1.zip test-with-comment.zip  test.txt

zipcmt --print --nodupes test/

#  ── test/test-with-comment-1.zip ─────────┐
#    This is an example test comment for zipcmt.
#
# Scanned 3 zip archives and found 1 unique comment
```

##### Save to directory example
```sh
ls -l test/
# test-no-comment.zip  test-with-comment.zip  test.txt

zipcmt --exportdir=~ test/
# Scanned 2 zip archives and found 1 comment

cat ~/test-with-comment-zipcomment.txt
#   This is an example test comment for zipcmt.
```

## Build

[Go](https://golang.org/doc/install) supports dozens of architectures and operating systems letting zipcmt to [be built for most platforms](https://golang.org/doc/install/source#environment).

```sh
# clone this repo
git clone git@github.com:bengarrett/myip.git

# access the repo
cd myip

# target and build the app for the host system
go build

# target and build for Windows 7+ 32-bit
env GOOS=windows GOARCH=386 go build

# target and build for OpenBSD
env GOOS=openbsd GOARCH=amd64 go build

# target and build for Linux on MIPS CPUs
env GOOS=linux GOARCH=mips64 go build
```

---

#### MyIP uses the following online APIs.

- [ipify API](https://www.ipify.org)
- [MYIP.com](https://www.myip.com)
- [Workshell MyIP](https://www.my-ip.io)
- [SeeIP](https://seeip.org)

The IP region data is from GeoLite2 created by MaxMind, available from
[maxmind.com](https://www.maxmind.com).

I found [Steve Azzopardi's excellent _import "context"_](https://steveazz.xyz/blog/import-context/) post useful for understanding context library in Go.
