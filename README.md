<p align="center">
  <code>gofilter</code> is a string filter for PF on OpenBSD using divert(4). 
</p>

#

### Main Features

- Filters packets based on strings
- No noticeable degradation of latency
- 30% of the original bandwidth available

#

> [!IMPORTANT]
> `gofilter` is pre-alpha software.

> [!NOTE]
> In [CHANGELOG.md] you can follow recent changes.
> [ROADMAP.md] shows our future plans.

***

### Usage

```
Usage of gofilter:
  -f string
        strings to filter, comma separated
  -p int
        divert socket listening port (default 700)
```

`gofilter` makes use of the kernel packet diversion mechanism [divert(4)].
Therefore, PF has to be configured accordingly. For example to filter inbound
DNS traffic:

```
pass in proto udp to any port 53 divert-packet port 700
```

### Caveats

> [!CAUTION]
> TCP segmentation offload will need to be disabled for the filter to not choke on
> large TCP packets:
>
> ```
> sysctl net.inet.tcp.tso=0
> ```

> [!CAUTION]
> IPv6 is currently broken at all.

### Performance

- `gofilter` reduces the available bandwidth down to 30%
- `gofilter` worsens the reliability of the traffic, as the standard deviation of the
  available bandwidth is very high

```
# Without gofilter
bandwidth min/avg/max/std-dev = 927.681/934.177/935.895/2.475 Mbps

# gofilter, without filter
bandwidth min/avg/max/std-dev = 0.023/310.585/925.562/293.994 Mbps

# gofilter, 1 filter
bandwidth min/avg/max/std-dev = 0.000/308.867/935.003/282.638 Mbps

# gofilter, 2 filter
bandwidth min/avg/max/std-dev = 0.023/313.504/916.121/261.767 Mbps

# gofilter, 10 filter
bandwidth min/avg/max/std-dev = 0.092/315.832/910.908/264.350 Mbps
```

### Contributing

See [CONTRIBUTING.md]

### Security

See [SECURITY.md]

### License

The package may be used under the terms of the ISC License a copy of
which may be found in the file [LICENSE].

Unless you explicitly state otherwise, any contribution submitted for inclusion
in the work by you shall be licensed as above, without any additional terms or
conditions.

[ROADMAP.md]:
https://github.com/huhnsystems/gofilter/blob/master/docs/ROADMAP.md
[CHANGELOG.md]:
https://github.com/huhnsystems/gofilter/blob/master/docs/CHANGELOG.md
[CONTRIBUTING.md]:
https://github.com/huhnsystems/gofilter/blob/master/docs/CONTRIBUTING.md
[SECURITY.md]:
https://github.com/huhnsystems/gofilter/blob/master/docs/SECURITY.md
[LICENSE]: https://github.com/huhnsystems/gofilter/blob/master/LICENSE
[divert(4)]: https://man.openbsd.org/divert.4
