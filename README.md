# Seed 🌱

Seed is a CLI tool that helps you quickly create directory structures from a tree representation. Whether you have a tree structure in your clipboard or a file, Seed can instantly "grow" it into a real directory structure.

<!--toc:start-->
- [Seed 🌱](#seed-🌱)
  - [Installation](#installation)
    - [Golang](#golang)
    - [Homebrew](#homebrew)
    - [Download Binary](#download-binary)
  - [Usage](#usage)
    - [From Clipboard](#from-clipboard)
    - [From String](#from-string)
    - [From File](#from-file)
  - [Input Format](#input-format)
    - [Using ASCII characters](#using-ascii-characters)
    - [Using spaces](#using-spaces)
    - [Using JSON](#using-json)
  - [Features](#features)
  - [Contributing](#contributing)
  - [Todo](#todo)
  - [License](#license)
  - [Author](#author)
  - [Acknowledgments](#acknowledgments)
<!--toc:end-->

## Installation

Currently supported is golang, homebrew, and good ol' downloading the binary yourself. More support to come. 

### Golang

```bash
go install github.com/jpwallace22/seed@latest
```

### Homebrew

```bash
brew tap jpwallace22/seed
brew install seed
```

### Download Binary

1. Go to the [Releases](https://github.com/jpwallace22/seed/releases) page
2. Download the appropriate archive for your system:

   **macOS (Apple Silicon / M1, M2)**
   - `seed_darwin_arm64.tar.gz`

   **macOS (Intel)**
   - `seed_darwin_amd64.tar.gz`

   **Linux (64-bit)**
   - `seed_linux_amd64.tar.gz`

   **Linux (ARM64)**
   - `seed_linux_arm64.tar.gz`

   **Windows**
   - `seed_windows_amd64.zip`

3. Extract the archive:
   
   **macOS/Linux**
   ```sh
   tar -xzf seed_darwin_arm64.tar.gz  # Replace with your downloaded file
   ```

   **Windows**
   - Right-click the downloaded zip file and select "Extract All"

4. (Optional) Move the binary to a directory in your PATH for easier access:

   **macOS/Linux**
   ```sh
   sudo mv seed /usr/local/bin/
   ```

You can verify the download using the provided checksums file.

## Usage

Seed can create directory structures in two ways:

### From Clipboard

```bash
seed -c
# or
seed --clipboard
```

This will read a tree structure from your clipboard and create the corresponding directories and files.

### From String 

```bash
seed "my-react-app
   ├── src
   │   ├── components
   │   ├── hooks
   │   ├── utils
   │   └── App.tsx
   ├── public
   │   └── index.html
   └── package.json"
```

### From File 

```bash
seed -f path/to/file
# or
seed --file path/to/file
```

## Input Format

Seed accepts tree structures in the common tree command format. For example:

### Using ASCII characters

```bash
my-project
├── src
│   ├── components
│   │   ├── Button.tsx
│   │   └── Card.tsx
│   ├── utils
│   │   └── helpers.ts
│   └── App.tsx
├── public
│   └── index.html
└── package.json
```
### Using spaces

> [!NOTE]  
> Only 4 spaces is supported at this time

```bash
my-project
    src
        components
            Button.tsx
            Card.tsx
        utils
            helpers.ts
        App.tsx
    public
        index.html
    package.json
```


You can generate this format using:
- The `tree` command in Unix-like systems
- VS Code extensions like "File Tree Generator"
- Or manually create it following the format above

### Using JSON

Seed also accepts JSON input that describes the directory structure. The JSON format should be an array containing directory/file objects and an optional report object:

```json
[
  {
    "type": "directory",
    "name": "my-project",
    "contents": [
      {
        "type": "directory",
        "name": "src",
        "contents": [
          {
            "type": "file",
            "name": "main.go"
          },
          {
            "type": "directory",
            "name": "utils",
            "contents": [
              {
                "type": "file",
                "name": "helper.go"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    "type": "report",
    "directories": 3,
    "files": 2
  }
]
```

Each object in the structure must have:
- `type`: Either "directory" or "file"
- `name`: The name of the directory or file
- `contents`: (Optional) An array of nested files and directories (only valid for directory type)

The report object is *optional* and contains:
- `directories`: Total number of directories
- `files`: Total number of files

Seed with throw if the report does not match what was created.

Example usage with JSON:

As string
```bash
seed -F json '{"type":"directory","name":"project","contents":[{"type":"file","name":"README.md"}]}'
# or
seed --format json '{"type":"directory","name":"project","contents":[{"type":"file","name":"README.md"}]}'
```

From clipboard

```bash
seed -F json -c
# or
seed --format json -c
```
From file

> [!NOTE] 
> Soon the filetype will be auto selected

```bash
seed -F json -f path/to/structure.json
# or
seed --format json -f path/to/structure.json
```

## Features

- 🚀 Fast directory structure creation
- 📋 Direct clipboard support
- 🌲 Supports standard tree format
- 📁 Creates both files and directories

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Todo

- ~~Implement ability to parse from file path~~
- ~~Add JSON support ~~
- Increase package manager distribution 
  - apt
  - pacman
  - choco
  - yum
- Add YAML support 
- Support StdIn
- flag to adjust spacing between 2 and 4 for people who write their own trees with just spaces



## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Justin Wallace ([@jpwallace22](https://github.com/jpwallace22))

## Acknowledgments

- Inspired by the Unix `tree` command
- Built with [Cobra](https://github.com/spf13/cobra)
