# zipcmt

[![Go Reference](https://pkg.go.dev/badge/github.com/bengarrett/zipcmt.svg)](https://pkg.go.dev/github.com/bengarrett/zipcmt) [![GoReleaser](https://github.com/bengarrett/zipcmt/actions/workflows/release.yml/badge.svg)](https://github.com/bengarrett/zipcmt/actions/workflows/release.yml)

Zipcmt is the super-fast batch, zip file comment viewer, and extractor.

- Using a modern PC with the zip files stored on a solid-state drive, zipcmt handles many thousands of archives per second.
- Comments convert to Unicode text for easy viewing, editing, or web hosting.<br>
<small>* comments can also be left as-is in their original CP437 or ISO-8859 text encoding.</small>
- Rarely see duplicate comments to avoid those annoying lists of identical site adverts.
- Transfer the source zip file's last modification date over to any saved comments.
- Tailored to both Windows and POSIX terminal platforms.

https://user-images.githubusercontent.com/513842/120908127-a7fc4580-c6aa-11eb-8a71-b3734fa53185.mp4

## Downloads

zipcmt is a standalone (portable) terminal application and doesn't require installation.

[Windows](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_Windows_Intel.zip), 
[macOS](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_macOS.tar.gz), 
[FreeBSD](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_FreeBSD_Intel.tar.gz
), 
[Linux](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_Linux_amd64.tar.gz
)
 and [Linux for ARM](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt_Linux_arm64.tar.gz
)

<small>Windows requires Windows 10 or newer, [users of Windows 7 and 8 can use zipcmt v1.3.10](https://github.com/bengarrett/zipcmt/releases/download/v1.3.10/zipcmt_Windows_Intel.zip).</small>

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

[ZST (Arch Linux package)](https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.pkg.tar.zst)
```sh
wget https://github.com/bengarrett/zipcmt/releases/latest/download/zipcmt.pkg.tar.zst
pacman -U zipcmt.pkg.tar.zst
```

## macOS unverified developer

Unfortunately, newer macOS versions do not permit the running of unsigned terminal applications out of the box. But there is a workaround.

In System Settings, Privacy & Security, Security, toggle **Allow applications downloaded from App store and identified developers**.

1. Use Finder to extract the download for macOS, `zipcmt_macOS.tar.tar.gz`
2. Use Finder to select the extracted `zipcmt` binary.
3. <kbd>^</kbd> control-click the binary and choose _Open_.
4. macOS will ask if you are sure you want to open it.
5. Confirm by choosing _Open_, which will open a terminal and run the program.
6. After this one-time confirmation, you can run this program within the terminal.


## Windows Performance

It is highly encouraged that Windows users temporarily disable **Virus & threat protection, Real-time protection**, or [create **Windows Security Exclusion**s](https://support.microsoft.com/en-us/windows/add-an-exclusion-to-windows-security-811816c0-4dfd-af4a-47e4-c301afe13b26) for the folders to be scanned before running `zipcmt`. Otherwise, the hit to performance is amazingly stark!

```
zipcmt -noprint 'C:\examples\'
```

This is the time taken with the default Microsoft Defender settings.

> Scanned 11331 zip archives and found 412 unique comments, taking 1m38.9398945s

**1 minute and 38 seconds** to scan 11,000 zip archives.

---

This is the expected performance on an SSD with a Windows Security Exclusion in place.

> Scanned 11331 zip archives and found 412 unique comments, taking 1.593534s

**1.6 seconds** to scan the same 11,000 zip archives!

## Usage

![Usage screenshot on Windows](usage.png)

## Example usage
#### Scan and print the comments
```sh
zipcmt test/

#  ── test/test-with-comment.zip ───────────┐
#    This is an example test comment for zipcmt.
#
# Scanned 4 zip archives and found 1 unique comment
```

#### Only print the comments
```sh
zipcmt --quiet test/

#   This is an example test comment for zipcmt.
```

#### Scan and save the comments

Linux, macOS, etc.

```sh
zipcmt --noprint --save=~ test/

# Scanned 4 zip archives and found 1 unique comment

cat ~/test-with-comment-zipcomment.txt

#   This is an example test comment for zipcmt.
```

Windows PowerShell

```powershell
zipcmt.exe --noprint --save='C:\Users\Ben\Documents' .\test\

# Scanned 4 zip archives and found 1 unique comment

cat 'C:\Users\Ben\Documents\test-with-comment-zipcomment.txt'

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

# target and build for Windows 10, 32-bit
env GOOS=windows GOARCH=386 go build

# target and build for OpenBSD
env GOOS=openbsd GOARCH=amd64 go build

# target and build for Linux on MIPS CPUs
env GOOS=linux GOARCH=mips64 go build
```

## Usages online

Reddit user [Iron_Slug](https://www.reddit.com/user/Iron_Slug/) in [r/bbs](https://www.reddit.com/r/bbs/) created a [huge online collection of BBS ads](https://www.ipingthereforeiam.com/bbs/gallery/zip/), many of which were captured using zipcmt.

The website Defacto2 has [a large collection of uncurated BBS ads](https://defacto2.net/f/b428b6e) also captured using zipcmt.
