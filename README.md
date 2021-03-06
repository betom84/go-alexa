[![Go Report Card](https://goreportcard.com/badge/github.com/betom84/go-alexa)](https://goreportcard.com/report/github.com/betom84/go-alexa)
[![codebeat badge](https://codebeat.co/badges/5cf553b7-d574-4a5f-8134-bbdab8ed034c)](https://codebeat.co/projects/github-com-betom84-go-alexa-master)
[![GoDoc](https://godoc.org/github.com/betom84/go-alexa?status.svg)](https://godoc.org/github.com/betom84/go-alexa)
[![CircleCI](https://circleci.com/gh/betom84/go-alexa.svg?style=shield)](https://circleci.com/gh/betom84/go-alexa)
[![Coverage Status](https://coveralls.io/repos/github/betom84/go-alexa/badge.svg?branch=master)](https://coveralls.io/github/betom84/go-alexa?branch=master)

# go-alexa

Go library to connect to [Amazon Alexa Smarthome Skill API v3](https://developer.amazon.com/de/docs/smarthome/understand-the-smart-home-skill-api.html)

## Table of Contents

- [About](#about)
- [Features](#features)
    - [Functional](#functional)
    - [Non-Functional](#non-functional)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
    - [Define devices discoverable for Alexa](#define-devices-discoverable-for-alexa)
    - [Create a DeviceFactory](#create-a-devicefactory)
    - [Custom directive processors](#custom-directive-processors)
- [Frequently Asked Questions (FAQ)](#faq)
- [Roadmap](#roadmap)
- [License](#license)

## About

This project was created to control smarthome devices (e.g. [Homematic](http://www.eq-3.com/products/homematic.html)) via Amazon Echo. It implements Amazon Alexa`s [Smart Home Skill API](https://developer.amazon.com/de/docs/smarthome/understand-the-smart-home-skill-api.html). So far, not all device types or possibilities enabled by Alexa are supported. But it should not be very challenging to add this capabilities as needed. Have a look at the [features](#features) and [roadmap](#roadmap) sections for more information.

## Features

### Functional:
- Discover defined smarthome devices by Alexa ([Alexa.Discovery Interface](https://developer.amazon.com/de/docs/device-apis/alexa-discovery.html))
- Authenticate an Alexa-User and grant access based on his Amazon profile ([Alexa.Authorization Interface](https://developer.amazon.com/de/docs/device-apis/alexa-authorization.html))
- Turn capable devices on or off ([Alexa.PowerController Interface](https://developer.amazon.com/de/docs/device-apis/alexa-powercontroller.html))
- Report health and current state for capable devices ([Alexa Interface](https://developer.amazon.com/de/docs/device-apis/alexa-interface.html))
- Query temperature sensor values ([Alexa.TemperatureSensor Interface](https://developer.amazon.com/de/docs/device-apis/alexa-temperaturesensor.html))

### Non-Functional:
- Easily extend supported functionality on your own
- Supports HTTP-BasicAuth for any requests handled (recommend use of https)
- Optionally all responses, send back to alexa, can be validated with the [official](https://github.com/alexa/alexa-smarthome/wiki/Validation-Schemas) JSON-Schema (extremely helpful for debugging)

## Requirements

- Of course, all those [prerequisites to Smart Home Skill Development](https://developer.amazon.com/de/docs/smarthome/understand-the-smart-home-skill-api.html#prerequisites-to-smart-home-skill-development) whould be very helpful
- I asume you authenticate an Alexa user by using [LWA](https://developer.amazon.com/de/docs/smarthome/authenticate-an-alexa-user-account-linking.html). To support other OAuth2 providers, you need to add a [custom directive processor](#custom-directive-processor) to handle authentication directives.

## Installation

You can install this package in the usual way by using [go get](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies).
```
go get github.com/betom84/go-alexa
```

## Usage

The following example starts an https server listening for Alexa directives. Of course you also can start an http server or don´t use BasicAuth (simply skip the assignments), but i would not recommend that.

```go
package main

import(
    "net/http"
    "os"
    
    "github.com/betom84/go-alexa/smarthome"
    "github.com/betom84/go-alexa/smarthome/common/discoverable"
)

func main() {
    endpoints := []discoverable.Endpoint{
        discoverable.Endpoint{...},
        discoverable.Endpoint{...},
    }

    authority := smarthome.Authority{
        ClientID:        "my-client-id",
        ClientSecret:    "my-client-secret",
        RestrictedUsers: []string{"amzn-user@mail.com"},
    }

    handler := smarthome.NewDefaultHandler(&authority, endpoints)
    handler.BasicAuth.Username = "alexa"
    handler.BasicAuth.Password = "supersecret"
    handler.DeviceFactory = DeviceFactory{}

    err = http.ListenAndServeTLS(":https", "ssl/certificate.pem", "ssl/private-key.pem", handler)
    if err != nil {
        panic(err)
    }
}
```
I asume you use [LWA for account linking](https://developer.amazon.com/de/docs/smarthome/authenticate-an-alexa-user-account-linking.html). You must add the scope `profile` in the "Account Linking" section in the skill developer console. Then the users email address can be compared with the `RestrictedUsers` to grant access.
In that example we only allow the Alexa-User with the e-mail `amzn-user@mail.com` to get access. See [amazon documentation](https://developer.amazon.com/de/docs/smarthome/authenticate-a-customer-permissions.html#getting-authorization) for more information.

### Define devices discoverable for Alexa

The endpoint model contains the information expected by Alexa´s [Discover](https://developer.amazon.com/de/docs/device-apis/alexa-discovery.html) directive. The properties from `Cookie` are used to create a device with the registered [DeviceFactory](#create-a-devicefactory). Also capabilities controllable by Alexa are defined within endpoint model. That means, devices created by the factory needs to satisfy the according [capability interface](https://godoc.org/github.com/betom84/go-alexa/smarthome/common/capabilities).

### Create a DeviceFactory

The `DeviceFactory` used above is needed to create a device which is capable of the action intended by Alexa. This device will be passed to the `DirectiveProcessor` to finally perform the intended action. By using `smarthome.NewDefaultHandler()` to create the handler, all supported processors are automatically added. Therefore devices need to satisfy the appropriate [capability interfaces](https://godoc.org/github.com/betom84/go-alexa/smarthome/common/capabilities) to work with these processors.

```go
type DeviceFactory struct{}

func (f *DeviceFactory) NewDevice(epType string, id string) (interface{}, error) {
    // return anything capable of the intended action
}
```

### Custom directive processors

You can also implement a `DirectiveProcessor` by your own.
```go
import (
    ...
    "github.com/betom84/go-alexa/smarthome/common"
)

type CustomDirectiveProcessor struct{}

func (p CustomDirectiveProcessor) Process(directive *common.Directive, device interface{}) (*common.Response, error) {
    // perform the action intended by the directive at the device
}

func (p CustomDirectiveProcessor) IsCapable(directive *common.Directive) bool {
    // return true if your processor is responsible for the given directive
}
```
And add it to the previously created handler.
```go
handler.AddDirectiveProcessor(CustomDirectiveProcessor{})
```

### Logging

By default only errors get logged to `stderr`. Change default logging behaviour by creating a new instance with custom writer and severenity by using `smarthome.NewDefaultLogger(...)`. Alternativly assign a custom `smarthome.Logger` implementation to `smarthome.Log` to override logging for your needs. 

## FAQ

t.b.d.

## Roadmap

- Support asynchronous responses and automatically use them if processing a directive takes to long
- Looking forward to find a better way to control window blinds (instead of simply turning them on and off)

## License

This project is licensed under the [MIT License](LICENSE).
