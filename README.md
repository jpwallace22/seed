# Seed 🌱

Seed is a CLI tool that helps you quickly create directory structures from a tree representation. Whether you have a tree structure in your clipboard or a file, Seed can instantly "grow" it into a real directory structure.

## Installation

```bash
go install github.com/jpwallace22/seed@latest
```

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

### From File (coming soon)

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

## Examples

1. **Basic React Project Structure**
   ```
   my-react-app
   ├── src
   │   ├── components
   │   ├── hooks
   │   ├── utils
   │   └── App.tsx
   ├── public
   │   └── index.html
   └── package.json
   ```

2. **Simple Node.js Project**
   ```
   node-api
   ├── src
   │   ├── controllers
   │   ├── models
   │   ├── routes
   │   └── index.js
   ├── tests
   └── package.json
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

- Implement ability to parse from file path
  - This should come with json/yml parsing
- flag to adjust spacing between 2 and 4 for people who write their own trees with just spaces


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Justin Wallace ([@jpwallace22](https://github.com/jpwallace22))

## Acknowledgments

- Inspired by the Unix `tree` command
- Built with [Cobra](https://github.com/spf13/cobra)
