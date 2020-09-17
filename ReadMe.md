# Echelon - hierarchical progress in terminals

[![PkgGoDev](https://pkg.go.dev/badge/github.com/roberChen/echelon)](https://pkg.go.dev/github.com/roberChen/echelon)

Cross-platform library to organize logs in a hierarchical structure.

Here is an example how it looks for running Dockerized tasks via [Cirrus CLI](https://github.com/roberChen/cirrus-cli):

![Cirrus CLI Demo](images/cirrus-cli-demo.gif)

Logger with progress bar:
[![asciicast](https://asciinema.org/a/GwKKnu5Z5J7hqQB4pON5o806W.svg)](https://asciinema.org/a/GwKKnu5Z5J7hqQB4pON5o806W)
## Features

* Customizable and works with any VT100 compatible terminal
* Supports simplified output for dumb terminals
* Implements incremental drawing algorithm to optimize drawing performance
* Can be used from multiple goroutines
* Pluggable and customizable renderers
* Works on Windows!

## Example

Please check [demo.go](demo/main.go) for a simple example or how *echelon* is used in [Cirrus CLI](https://github.com/roberChen/cirrus-cli).
