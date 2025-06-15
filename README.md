# networkd-ipmon

networkd-ipmon is a dispatcher daemon for systemd-networkd that monitors network interface IP changes., When a network interface IP changes, executables in a specified directory will be called with relevant information. It is similar to [networkd-dispatcher](https://gitlab.com/craftyguy/networkd-dispatcher), but more specialized, and uses systemd-networkd's new DBus APIs.

## Usage

```
$ networkd-ipmon <dir>
```

Where `<dir>` is the path to the directory containing the executables. For each executable named `<executable>`, a `<executable>.json` configuration file must exist in the same directory to specify which interfaces and IP properties to watch for that executable.

## JSON configuration File Format

```json
{
  "interfaces": [],
  "properties": []
}
```

`interfaces` specifies the network interfaces to watch.

Type: `string[]`
Example: `[ "eth0", "eth1" ]`

`properties` specifies the types of IP to watch.

Type: `("IPV6_ADDRS" | "IPV4_ADDRS" | "PD_ADDRS")[]`,
Example: `[ "IPV4_ADDRS", "IPV6_ADDRS" ]`

## Executable Environment

When an interface's IP changes, all monitoring executables will be executed with the following environment variables:

```
OLD_IPV4_ADDRS
IPV4_ADDRS
OLD_IPV6_ADDRS
IPV6_ADDRS
OLD_PD_ADDRS
PD_ADDRS
```

Unwatched properties are omitted. A `OLD_*` variable contains the old value before the change, and is omitted if the interface previously had no IP.
