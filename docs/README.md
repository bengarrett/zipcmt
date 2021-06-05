# zipcmt

[![Go Reference](https://pkg.go.dev/badge/github.com/bengarrett/zipcmt.svg)](https://pkg.go.dev/github.com/bengarrett/zipcmt) [![GoReleaser](https://github.com/bengarrett/zipcmt/actions/workflows/release.yml/badge.svg)](https://github.com/bengarrett/zipcmt/actions/workflows/release.yml)

Zipcmt is the super-fast, batch, zip file comment viewer, and extractor.

- Using a modern PC with the zip files stored on a solid-state drive, zipcmt handles many thousands of archives per second.
- Comments convert to Unicode text for easy viewing, editing, or web hosting.<br>
<small>* comments can also be left as-is in their original CP437 or ISO-8859 text encoding.</small>
- Rarely see duplicate comments to avoid those annoying lists of identical site adverts.
- Transfers the source zip file last modification date over to any saved comments.
- Tailored to both Windows and POSIX terminal platforms.

## Downloads

<small>zipcmt is a standalone (portable) terminal application and doesn't require installation.</small>

- [Windows](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_Windows_Intel.zip)
- [macOS](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_macOS_Intel.tar.gz
), [or for the Apple M chip](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_macOS_M-series.tar.gz
)
- [FreeBSD](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_FreeBSD_Intel.tar.gz
)
- [Linux](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_Linux_Intel.tar.gz
)

### Packages

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

[RPM (Red Hat package)](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.rpm)
```sh
wget https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.rpm
rpm -i zipcmt.rpm
```

Windows [Scoop](https://scoop.sh/)
```sh
scoop bucket add bengarrett https://github.com/bengarrett/zipcmt.git
scoop install bengarrett/zipcmt
```

## Usage

![Usage screenshot on Windows](usage.png)

## Example usage
### Print
```sh
zipcmt test/

#  ── test/test-with-comment.zip ───────────┐
#    This is an example test comment for zipcmt.
#
# Scanned 4 zip archives and found 1 unique comment
```

##### Quiet example
```sh
zipcmt --quiet test/

#   This is an example test comment for zipcmt.
```

##### Save to directory examples

Linux, macOS, etc.

```sh
zipcmt --save=~ test/

# Scanned 2 zip archives and found 1 comment

cat ~/test-with-comment-zipcomment.txt

#   This is an example test comment for zipcmt.
```

Windows

```powershell
zipcmt --save="C:\Users\Ben\Documents" test

# Scanned 2 zip archives and found 1 comment

cat "C:\Users\Ben\Documents\test-with-comment-zipcomment.txt

#   This is an example test comment for zipcmt.
```

## Build

[Go](https://golang.org/doc/install) supports dozens of architectures and operating systems letting zipcmt to [be built for most platforms](https://golang.org/doc/install/source#environment).

```sh
# clone this repo
git clone git@github.com:bengarrett/zipcmt.git

# access the repo
cd zipcmt

# target and build the app for the host system
go build

# target and build for Windows 7+ 32-bit
env GOOS=windows GOARCH=386 go build

# target and build for OpenBSD
env GOOS=openbsd GOARCH=amd64 go build

# target and build for Linux on MIPS CPUs
env GOOS=linux GOARCH=mips64 go build
```
