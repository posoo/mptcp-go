# mptcp-go
MPTCP testing suite in Go

## Prerequisites

### Kernel Support

Check and enable MPTCP support in your kernel. This has to be done on both the client and server sides. You can check if your kernel supports MPTCP by running:
```bash
sysctl net.mptcp.enabled
```
If it returns `net.mptcp.enabled = 1`, then MPTCP is enabled. If not, you can enable it by running:

```bash
sudo sysctl -w net.mptcp.enabled=1
```
### Endpoint Configuration

Depending on the network setup, there could be multiple subflow topoligies within a single MPTCP connection.
In our case, we assume that the client is multihomed (connected with multiple accessiable network interfaces), while the server is single-homed.
In addition, subflows are actively established by the client through `MP_JOIN` on different paths.

Therefore, we set up the client with multiple `subflow` endpoints as follows:
```bash
# subsitute with your actual devices and IP addresses
sudo ip mptcp endpoint add 10.0.1.4 dev eth0 id 1 subflow
sudo ip mptcp endpoint add 10.0.2.4 dev eth1 id 2 subflow
```

## MPTCP Checker

Similar to [multipath-tcp/mptcp-check](https://github.com/multipath-tcp/mptcp-check), we developed a [webserver](/cmd/webserver) to verify whether an HTTP connection is established over MPTCP. 
The difference is that our webserver directly acquires the MPTCP status from the underlying connection, rather than calling the `ss` command and parsing its output.

In addition, we provide a [Dockerfile](/Dockerfile) to run the webserver in a container persistently.
```bash
docker build -t mptcp-checker .
docker run -d --name mptcp-checker --restart=always -p 8684:8684 mptcp-checker
```

## References
- [ip-mptcp(8) — Linux manual page](https://www.mankier.com/8/ip-mptcp)
- [mptcp.dev](https://www.mptcp.dev/)
- [mptcpize - Man Page](https://www.mankier.com/8/mptcpize) - to force applications to use MPTCP
    - `mptcpize run curl check.mptcp.dev`
- [mptcpd - Man Page](https://www.mankier.com/8/mptcpd) - to perform MPTCP path management related operations
