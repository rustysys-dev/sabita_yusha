# Sabita Yusha

This project is geared towards allowing for custom / programmable buttons for VIA/QMK keyboards. It is primarily built using the yushakobo Quick Paint keyboard, and while i imagine it should work with other keyboards compliant to QMK on Linux I cannot guarantee it.

## Features

This project implements the following features most of which are explained in the [config.yaml.example](./config.yaml.example) file.

- Automatic Detection of Existing Target Keyboard
- Automatic Detection of Adding / Removal of Target Keyboard
- Automatic Config Generation When Non-existant
- Execution of Handlers for Linux `KEY_MACRO1` ~ `KEY_MACRO30`
  - Bash Command Handler
  - Slack Message Command Handler

This project assumes the user has setup their keyboard via VIA/QMK or REMAP software to emit `PROGRAMMABLE_BUTTON1` ~ `PROGRAMMABLE_BUTTON30` which emit Linux `KEY_MACRO` events of similar number.

This project assumes the user knows the name of the HID report which handles these buttons. For reference, the yushakobo Quick Paint will generate many report files depending on the configuration. The report name for `yushakobo Quick Paint` is `yushakobo Quick Paint Consumer Control`.

## Usage

### Testing

For testing purposes the application can be run to discover required information
in order to properly config the app. Using the following command, you will be
prompted to tap an appropriately configured button on your keyboard. If your
device is properly configured you will be given the `target_device_name` for
your config.

```
go run . -test
```

### Run

To test run officially...

```
go run .
```

### Package / Install

To officially install for RPM based environments:

```
# edit this file as you see fit
go run cmd/package_app.go
```

And Install it.

#### Fedora (dnf/yum)

```
sudo dnf install dist/*.rpm
```

### Auto-Start

Various window managers and desktop environments have their ways of setting
applications to automatically start, please reference your environment's
documentation if it isn't listed here. Feel free to submit a request or PR
to add your environment to this list.

#### Hyperland

Set to start on login (hyperland):

```
# Add the following to $HOME/.config/hypr/hyprland.conf
exec-once = sabita_yusha
```
