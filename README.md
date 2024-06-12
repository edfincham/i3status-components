# i3Status Components

The contents of this repository have been begged, borrowed, and stolen from other sources with a hefty dose of ChatGPT. Of particular note is [go-i3status](https://github.com/ghedamat/go-i3status/tree/master).

The module uses the [i3bar-protocol](http://i3wm.org/docs/i3bar-protocol.html) to build an interactive statusbar.

## Requirements
- [Font Awesome 4](https://fontawesome.com/v4/). Old, but you try installing a newer icon pack on Linux


## Usage
To get the default behaviour just install the *go-i3status* and set it as the `status_command` in your `.i3config`. For instance:
```
bar {
        ...
        status_command exec /path/to/i3status-components/bin/i3br
        ...
}
```

## Development
Relevant commands for running, building, and installing can be found in the `Makefile`.
