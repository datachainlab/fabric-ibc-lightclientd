# Fabric-IBC light client daemon

## Introduction
`fabric-ibc-lightclientd` is a standalone daemon meant to offload implementation cost of Fabric-IBC light client logic.
By using this internally, IBC implementers can easily develop ClientState for Fabric-IBC in their own IBC implementation.

## Usage
`fabric-ibc-lightclientd` runs as a gRPC server listening to a specified port.
ClientState (for Fabric IBC) in your IBC implementation just works as a shim layer for this server.

```
$ fabric-ibc-lightclientd -port 60000
```
