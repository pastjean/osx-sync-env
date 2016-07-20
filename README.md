# osx-env-sync

[![Build Status](https://travis-ci.org/pastjean/osx-sync-env.svg?branch=master)](https://travis-ci.org/pastjean/osx-sync-env)
[![Go Report Card](https://goreportcard.com/badge/github.com/pastjean/osx-sync-env)](https://goreportcard.com/report/github.com/pastjean/osx-sync-env)

> An easy to use environment variable manager. It loads the environment
variables exported in the user shell into the osx GUI app context
using launchctl.

## First time set-up

Here's a series of commands that would set everything up for you

- Download from https://github.com/pastjean/osx-sync-env/releases and set in your path
```
curl https://github.com/pastjean/osx-sync-env/releases/download/latest/osx-sync-env.tar.gz | tar xzf
```

- Move the binary to your favorite bin dir
- Then
```
osx-sync-env install
osx-sync-env sync
```

or for go fu's

```bash
go get -u github.com/pastjean/osx-sync-env
```

## Usage

In your shell type `osx-sync-env --help` and you'll see all the available options.


