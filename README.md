# Fisherman 🎣
Fisherman is a simple cross-platform command-line tool for managing Git pre-commit hooks in a friendly way. It allows you to easily lint, run tests and other code-quality functions every time you commit your code.

## Features
- **Improve Code Quality**: Run code quality checks automatically on every commit.
- **Enable/Disable Hooks**: Quickly disable/re-enable commit hooks with 1 command.
- **Temporarily Disable Hooks**: Temporarily disable hooks for a set period if you need to get work done quickly with the `-disabled-for 1h` flag.
- **Cross-Platform Support**: Works on macOS, Windows, and Linux.

## Installation
You can download the executable for your operating system from the [GitHub Releases](https://github.com/TomStowe/fisherman/releases) page, or you can use `go install github.com/TomStowe/fisherman@latest`

### macOS
1. Download the `fisherman` executable from the [GitHub Releases](https://github.com/yourusername/fisherman/releases) page.
2. Move the executable to a directory in your `PATH`, e.g., `/usr/local/bin`:
   ```bash
   mv /path/to/downloaded/fisherman /usr/local/bin/fisherman
   chmod +x /usr/local/bin/fisherman
   ```
3. Verify the installation:
    ```bash
    fisherman -h
    ```

### Windows
1. Download the fisherman.exe file from the GitHub Releases page.
2. Move the executable to a directory included in your PATH, e.g., `C:\Windows\System32`.
3. Verify the installation by opening the Command Prompt and running:
    ```cmd
    fisherman -h
    ```
### Linux
1. Download the fisherman executable from the GitHub Releases page.
2. Move the executable to a directory in your PATH, e.g., /usr/local/bin:
    ```bash
    mv /path/to/downloaded/fisherman /usr/local/bin/fisherman
    chmod +x /usr/local/bin/fisherman
    ```
3. Verify the installation:
    ```bash
    fisherman -h
    ```

## Usage

Fisherman allows you to enable or disable Git pre-commit hooks with simple commands.

### Setup and Enable a Pre-Commit Hook
To enable a pre-commit hook using a script file:

```bash
fisherman -file=path/to/your/hook_script.sh
```

### Disable the Pre-Commit Hook
To disable the pre-commit hook:

```bash
fisherman -disable
```

### Temporarily Disable the Pre-Commit Hook
To temporarily disable the pre-commit hook for 2 hours:

```bash
fisherman -disabled-for 2h
```
Where units are:
| Unit | Value  |
| ---- | ------ |
| w    | week   |
| d    | day    |
| h    | hour   |
| m    | minute |
| s    | second |

### Re-Enable the Pre-Commit Hook
To re-enable an existing pre-commit hook:

```bash
fisherman -enable
```

## Contributing

Contributions are welcome! If you have ideas, find bugs, or want to add new features, feel free to submit a pull request or open a new issue.

### Development

To build Fisherman from source, ensure you have Go installed, and then run the following commands:

1. Clone the repository:
    ```bash
    git clone https://github.com/TomStowe/fisherman.git
    cd fisherman
    ```
2. Build the executable:
    ```bash
    go build -o fisherman main.go
    ```
3. Run the tool:
    ```bash
    ./fisherman -h
    ```

## License

Fisherman is licensed under the MIT License. See the [LICENSE file](LICENSE.md) for more details.