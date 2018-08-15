# luna


## Prerequisites
* [golang](https://golang.org/doc/install)
* [aws-sdk-go](https://github.com/aws/aws-sdk-go)
* [AWS Security Credentials](https://docs.aws.amazon.com/ko_kr/general/latest/gr/aws-security-credentials.html)

## Install
go get github.com/lunamint/luna<br/>
cd $GOPATH/src/github.com/lunamint/luna<br/>
go install

## planb Notification
* Register the topic name "sipchanged" in Amazon Simple Notification Service.<br/>
  ( Planb enabled without setting )

## Usage
* luna net
* luna planb

## Warning
* Never use planb on validator node.
