# Fabric-IBC light client daemon

## Introduction
`fabric-ibc-lightclientd` is a standalone daemon meant to offload implementation cost of Fabric-IBC light client logic.
By using this internally, IBC implementers can easily write ClientState for Fabric-IBC in their own IBC implementation.

## Usage
`fabric-ibc-lightclientd` listens to a specified local TCP port.
A main body of IBC implementation just have to have a shim layer of this.

```
$ fabric-ibc-lightclientd --port 60000
```
