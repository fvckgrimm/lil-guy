# lil-guy

lil-guy is a fun, customizable command-line animation tool that displays cute characters with messages and can output from stdin.

<img src="/assets/showcase.gif"/>

## Features

- Multiple pre-defined characters
- Support for custom characters via TOML configuration
- Multi-line character support
- Animated character display
- Customizable messages
- Debug mode for troubleshooting

## TODO:

- Fix handling of err being piped into lil-guy


## Installation

### Releases

Prebuilt binaries can be found in [releases](https://github.com/fvckgrimm/lil-guy/releases)

### Build from source

1. Clone the repository

```bash
git clone https://github.com/fvckgrimm/lil-guy.git && cd lil-guy
```

2. Build the project

```bash
go build -v -o lil-guy
```

## Configuration

lil-guy uses TOML file for character configuration, the default location which it's read from is:

```
~/.config/lil-guy/characters.toml
```

## Usage

Run lil-guy with the following command:

```bash
./lil-guy [flags]
```

Available flags:

* -m, -message: Set the message to display (default: "Hello, I'm lil guy!")
* -c, -character: Choose the character to display (default: "default")
* -debug: Run in debug mode for troubleshooting

Examples:

```bash
./lil-guy -message "Hello, World!" -character cat
./lil-guy -message "I'm a dog" -character dog
./lil-guy -message "Multi-line test" -character multiline_example
./lil-guy -debug -character fumo_1
```

You can also pipe other commands into lil-guy: 

```bash
pip install discord | lil-guy -message "I can handle this myself..." -character blinkcat
```

For use in scripts:

```bash
# in something.sh
lil-guy -message "Doing something..." -character cat &
PID=$!

# Do your task here
echo "thanks for the work"

kill $PID 2>/dev/null
```


## Adding New Characters

To add a new character, simply add a new entry to your `characters.toml` file. No changes to the source code are required.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

