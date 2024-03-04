
# QR Code Generator

This program generates QR codes from URLs and saves them as PNG or SVG files.

## Installation

To install and run this program, you'll need to have Go installed on your machine. 

1. Clone the repository to your local machine:

```bash
git clone https://github.com/mtzvd/qr-generator.git
```

1. Navigate to the cloned directory:

```bash
cd qr-generator
```

1. Build the program:

```bash
go build
```

This will create an executable file in the current directory.

## Usage

To generate a QR code, you can use the following flags:

- `-u`: URL to generate QR code for (required, max length 2048)
- `-l`: Correction level (options: L, M, Q, H; default "M")
- `-f`: Output format (options: png, svg; default "png")
- `-s`: Size of the QR code (default 256, min 100, max 4096)
- `-d`: Directory to save the file (default is current directory)

### Examples

Generate a QR code as PNG with medium correction level and save it to the current directory:

```bash
./qr-generator -u 'https://www.example.com' -s 256 -l M -f png
```

Generate a QR code as SVG with high correction level and save it to a specific directory:

```bash
./qr-generator -u 'https://www.example.com' -s 512 -l Q -f svg -d /path/to/save
```

## Contributing

Contributions are welcome. Feel free to open a pull request with any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
