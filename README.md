# N9E-Compatible Alarm Repeater

[English](./README.md) | [简体中文](./doc/README_zh-cn.md)

> Note: The json data structure of the N9E in this project has been adjusted. It is not the original N9E data structure. If you want to be compatible with the N9E data structure, please adjust the structure on the receiver side.

## Main functions

1. Use N9E's alarm callback function to receive alarm information and convert it into a file
2. Alarm information is contained in the file, which can be processed at will, such as uploading it to other places
3. At the sender, the file will be converted into a remote call from N9E and sent to the real callback address to forward alarm information

## Application scenarios

1. The N9E operating environment does not support direct network invocation of remote services
2. Unified collection and recycling of N9E alarm information

## Run environment

- go version 1.23.0+
- x86 linux
- arm linux

## TODO

- [x] Information encryption during forwarding
- [x] Forward backup information and clean it up regularly
- [x] Support log level customization
- [x] Customize encryption file prefixes and suffixes to be compatible with various types of gatekeepers
- [x] Send failed retry
- [] Support alarm information filtering, forward alarms containing specific characters, or used to block certain alarms
- [x] Record forwarding log content to local files, to facilitate search tracking, to avoid disputes, support regular clean-up

