# jsonencoder

A command line tool written in Go that can accept a JSON string or file and encode or decode it. This tool helps solve the problem of providing plain text JSON that needs to be parseable when there are escape characters.

## Features

- **Encode JSON**: Convert JSON to an escaped string format that can be safely embedded in other contexts
- **Decode JSON**: Convert escaped JSON strings back to their original format  
- **File Support**: Read JSON input from files or provide it directly as command line arguments
- **Validation**: Ensures input is valid JSON before processing
- **Error Handling**: Clear error messages for invalid input

## Installation

### From Source

```bash
git clone https://github.com/Knighton-Dev/jsonencoder.git
cd jsonencoder
go build -o jsonencoder
```

### System-Wide Installation

After building the binary, you can install it system-wide to use from anywhere:

#### Option 1: Install to /usr/local/bin (macOS/Linux - Recommended)

```bash
# Build the binary
go build -o jsonencoder

# Make it executable and install system-wide
chmod +x jsonencoder
sudo mv jsonencoder /usr/local/bin/

# Verify installation
jsonencoder --help
```

#### Option 2: Install to user directory (macOS/Linux)

```bash
# Create user bin directory if it doesn't exist
mkdir -p ~/.local/bin

# Build and install to user directory
go build -o jsonencoder
chmod +x jsonencoder
mv jsonencoder ~/.local/bin/

# Add to PATH (add this line to ~/.zshrc, ~/.bashrc, or ~/.bash_profile)
export PATH="$HOME/.local/bin:$PATH"

# Reload your shell or source the profile
source ~/.zshrc  # or ~/.bashrc or ~/.bash_profile
```

#### Option 3: Add current directory to PATH

```bash
# Build the binary in the repository
go build -o jsonencoder

# Add the repository directory to your PATH
export PATH="$PWD:$PATH"

# To make permanent, add to your shell profile
echo 'export PATH="'$PWD':$PATH"' >> ~/.zshrc
```

After installation, you can run `jsonencoder` from any directory without the `./` prefix.

## Usage

```
jsonencoder - A CLI tool to encode and decode JSON strings

Usage:
  jsonencoder [options] <command> <input>

Commands:
  encode    Encode JSON (escape for embedding)
  decode    Decode JSON (unescape)

Options:
  -f, --file    Read input from file instead of command line argument
  -h, --help    Show this help message
```

## Examples

### Encoding JSON

Encode a JSON string (adds escape characters):

```bash
jsonencoder encode '{"key": "value"}'
# Output: "{\"key\": \"value\"}"
```

Encode JSON from a file:

```bash
jsonencoder -f encode input.json
```

### Decoding JSON

Decode an escaped JSON string:

```bash
jsonencoder decode '"{\"key\": \"value\"}"'
# Output: {"key": "value"}
```

Decode JSON from a file:

```bash
jsonencoder -f decode encoded.json
```

### Round Trip Example

```bash
# Start with original JSON
echo '{"name": "John", "age": 30}' > original.json

# Encode it
jsonencoder -f encode original.json > encoded.json

# Decode it back
jsonencoder -f decode encoded.json > decoded.json

# Compare original and decoded
diff original.json decoded.json
# Should show no differences
```

> **Note**: Examples above assume system-wide installation. If running from the repository directory without installation, prefix commands with `./` (e.g., `./jsonencoder encode ...`)

## Use Cases

This tool is particularly useful when you need to:

- Embed JSON data as a string literal in code
- Store JSON data in configuration files that require escaping
- Process JSON data that has been escaped for transmission
- Debug JSON escaping issues
- Prepare JSON for APIs that expect escaped JSON strings

## Testing

Run the test suite:

```bash
go test -v
```

## License

See [LICENSE](LICENSE) file for details.
