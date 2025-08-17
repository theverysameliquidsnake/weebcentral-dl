# Manga Downloader
A script for downloading manga from weebcentral.com written in Go.

## Usage

### Command Line Arguments

| Option | Description |
|:---|:---|:---|
| `-t, --title` | The title of the manga to search for. |
| `-f, --first` | Download only chapters equals or newer than specified. |
| `-l, --last` | Download only chapters equals or older than specified. |
| `-c, --compress` | Compress downloaded chapters to .zip or .cbz format. |
| `-o, --output` | The folder where downloaded manga will be saved. |
| `-v, --verbose` | Enable detailed debug output. |
| `-i, --install` | Install required Playwright dependencies (required to be enabled for the first run). |
