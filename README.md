# AppWrap

Creates a very minimal application bundle from a binary.

## Install

```shell
$ go install github.com/fzwoch/appwrap@latest
```

## Usage

```shell
$ appwrap /path/to/some/binary
```

will create:

```
  binary.app/
└── Contents
    ├── Info.plist
    └── MacOS
        └── binary

```
in the current working directory.
