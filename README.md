# srcds_proxy - An UDP proxy for the SRCDS protocol

This is a small Go project of mine to proxy all connections established to a SRCDS server. It basically NAT every connection.

The purpose of this proxy is to cache, filter or alter the requests sent to the server.
It can be used to protect the server against some DOS attacks, it can reduce the load on the server and allow to introduce custom behavior.
