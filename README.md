# dimpsled

## dimpsled - Dim PS LED
This project is for use with a used rule to set a random sweet dim light for your connected PS4 controller (maybe PS5).

Install the binary in `/usr/local/bin`
```sh
sudo GOBIN=/usr/local/bin go install github.com/bastienbc/dimpsled
```

Create a 99-psled.rules with this content:
```
ACTION=="add", KERNEL=="js*", SUBSYSTEM=="input", ATTRS{id/vendor}=="054c", ATTRS{id/product}=="09cc" RUN+="/usr/local/bin/dimpsled -d '/sys%p'"
```

Beware of the `/sys%p` in the udev rule:
`udev` gives the device path INSIDE `/sys`, but `dimpsled` requires an absolute file path.
