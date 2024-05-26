# tldr-ai

`tldr-ai` is a Go-based wrapper for the `tldr` command-line utility. It extends the functionality of `tldr` by utilizing the GPT API to generate descriptions for commands that are not present in the `tldr` database. This ensures you always have a concise and helpful explanation of any command.

## Features

- **Extended Command Support**: Provides descriptions for commands not found in the `tldr` database by leveraging the GPT API.
- **User-Friendly**: Maintains the simplicity and conciseness of the `tldr` utility.
- **Fast and Efficient**: Written in Go for performance and ease of deployment.

## Installation

1. **Install Go**: Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

2. **Clone the Repository**:
    ```sh
    git clone https://github.com/yourusername/tldr-ai.git
    cd tldr-ai
    ```

3. **Build the Project**:
    ```sh
    go build -o tldr-ai
    ```

4. **Add tldr-ai to Your PATH**:
    ```sh
    export PATH=$PATH:/path/to/tldr-ai
    ```

## Usage

`tldr-ai` works just like `tldr`. Simply use it followed by the command you need help with.

```sh
tldr-ai <command>
```

## Examples
- Get a description for ls:
```sh
tldr-ai ls
```

- Get a description for foo (not in tldr database):
```sh
tldr-ai foo
```

## Configuration
To use the GPT API, you need to set your API key. You can do this by setting the GPT_API_KEY environment variable.
```sh
export OPENAI_API_KEY=your-api-key
```
