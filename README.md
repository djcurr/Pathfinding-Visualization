# Pathfinding Visualization

A fast WebAssembly application that allows users to visualize pathfinding algorithms. Users can place walls, move the start and end nodes, and see the algorithm process in action. 
This was built in Go using WebAssembly, Go, and JavaScript.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
- [Contributing](#contributing)

## Features

Pathfinding Visualization offers the following features:

- Visualize pathfinding algorithms such as Dijkstra's algorithm
- Place walls and move the start and end nodes
- See the algorithm process in action

## Technologies Used

Pathfinding Visualization is built using the following technologies:

- Go
- WebAssembly: A binary instruction format for a stack-based virtual machine that enables running code written in multiple languages on the web
- HTML/CSS/JavaScript

## Getting Started

To get started with Pathfinding Visualization, follow these steps:

1. Clone the repository
2. Build the WebAssembly binary using `GOOS=js GOARCH=wasm go build -o main.wasm`
3. Serve the files using a web server such as [Go's built-in web server](https://golang.org/pkg/net/http/) or [Node.js's http-server](https://www.npmjs.com/package/http-server)
4. Open the web page in a web browser

## Contributing

Contributions to Pathfinding Visualization are welcome! If you have any bug reports, feature requests, or other suggestions, please open an issue or submit a pull request.

