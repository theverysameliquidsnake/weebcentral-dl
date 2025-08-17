# Manga Downloader
A script for downloading manga from weebcentral.com written in Go.

## Usage

### Installation
1. **Clone this repository:**
```
git clone https://github.com/theverysameliquidsnake/weebcentral-dl.git
cd weebcentral-dl
```

2. **Compile executable file:**
```
go build -o weebcentral-dl
```

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

### Examples

1. **Search for manga and download allchapters:**
```
weebcentral-dl -t Dandadan
```

2. **Search for manga and download chapters within specified range:**
```
weebcentral-dl -t Dandadan -f 10 -l 25
```
3. **Search for manga, download chapters from chapter 10 and compress it to .cbz:**
```
weebcentral-dl -t Dandadan -f 10 -c cbz
```
4. **Search for manda, download chapters up to 50 and save it to `~/manga` folder:**
```
weebcentral-dl -t Dandadan -l 50 -o ~/manga
```

5. **Install Playwright dependencies if needed:**
```
weebcentral-dl -i
```
