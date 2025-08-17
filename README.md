# Manga Downloader
A script for downloading manga from weebcentral.com written in Go.

## Usage

### Command Line Arguments

| Option | Description |
|:---|:---|
| `-h, --help` | display help message and exit |
| `-t, --title` | search manga by specified title |
| `-f, --first` | filter chapters equal or newer than specified number |
| `-l, --last` | filter chapters equal or older than specified number |
| `-p, --prefix` | filter chapters by specified chapter prefix |
| `-o, --output` | download to specified directory |
| `-c, --compress` | compress to specified format (cbz or zip) |
| `-i, --install` | install required Playwright dependencies |
| `-v, --verbose` | enable detailed log output |
