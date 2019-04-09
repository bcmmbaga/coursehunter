# coursehunter
coursehunters.net courses downloader

[![Build Status](https://travis-ci.org/bcmmbaga/coursehunter.svg?branch=master)](https://travis-ci.org/bcmmbaga/coursehunter) [![Go Report Card](https://goreportcard.com/badge/github.com/bcmmbaga/coursehunter)](https://goreportcard.com/report/github.com/bcmmbaga/coursehunter)

## Download executable from links below

- [linux-amd64](https://github.com/bcmmbaga/coursehunter/releases/download/v0.1.0-beta/hunterD-linux-amd64)
- [windows-amd64](https://github.com/bcmmbaga/coursehunter/releases/download/v0.1.0-beta/hunterD-windows-amd64.exe)

## Usage

```bash

coursehunter [command] [options...]

    COMMAND:
        resume resume interrupted downloads

    OPTIONS:
        -n      coursename
        -e      email
        -p      password
        -start  video index to start at (default: 1)

Example:
    download:
    coursehunter -n <coursename> -e <email> -p <password>

    resume:
    coursehunter resume -start <start index> 

```