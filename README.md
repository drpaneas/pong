# Pong

Pong is a classic arcade game developed by Atari and originally released back in 1972.
This is a remake of the game, using Go programming language and Ebiten game library.
The gameplay mechanics here are slightly different from the original.

## How to Play

The goal of Pong is to score points by hitting the ball past your opponent's paddle.
You control the right paddle using the arrow keys (up and down).
The game ends when one player reaches 10 points.

## Play Online

You can play the game online via your web browser at <https://drpaneas.net/pong/>

## Features

- Single player only.
- Three levels of progressive difficulty.
- Sound effects and background music.

## How to Build and Run

To build and run Pong, you will need Go 1.20 installed on your system.
Once you have Go installed, follow these steps:

1. Clone the repository: `git clone https://github.com/your-username/pong.git`
2. Change into the repository directory: `cd pong`
3. Run the game: `go run main.go`

## How to Build for the Browser

1. Add the following build tag to the `main.go`: `// +build js,wasm`
2. Copy `wasm_exec.js` into the game's dir: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .`
3. Have an HTML file that the `BODY` looks like this:

```html
<canvas id="canvas" width="640" height="480"></canvas>
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("game.wasm"), go.importObject)
    .then((result) => {
      go.run(result.instance);
    });
</script>
```

4. Use the `GOOS` and `GOARCH` environment variables to compile your game for WebAssembly: `env GOOS=js GOARCH=wasm go build -o game.wasm`
5. The previous step should have created `game.wasm`.
6. Start a server (e.g. `python -m http.server`) and visit `http://localhost:8000` in your web browser to play your game.
