# Pong

Pong is a classic arcade game developed by Atari and originally released back in 1972.
This is a remake of the game, using Go programming language and [Ebitenengine](https://github.com/hajimehoshi/ebiten).
The gameplay mechanics here are slightly different from the original.

![screenshot](screenshot.png)

## How to Play

The goal of Pong is to score points by hitting the ball past your opponent's paddle.
You control the right paddle using the arrow keys (up and down).
The game ends when one player reaches 10 points.

## Play Online

You can play the game online via your web browser at drpaneas.net/pong

## Features

- Single player only.
- Three levels of progressive difficulty.
- Sound effects and background music.

## How to Build and Run

To build and run Pong, you will need Go 1.20 installed on your system.
Once you have Go installed, follow these steps:

1. Clone the repository: `git clone https://github.com/your-username/pong.git`
2. Change into the repository directory: `cd pong`
3. Run the game: `go build && ./pong`

## How to Build for the Browser

1. Copy `wasm_exec.js` into the game's wasm dir: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./wasm/`
2. Build the Web Assembly version of the game: `GOOS=js GOARCH=wasm go build -o ./wasm/game.wasm`
3. Have an HTML file that the `BODY` looks like this:

```html
<canvas id="canvas" width="1280" height="720"></canvas>
<script src="./wasm/wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("./wasm/game.wasm"), go.importObject)
    .then((result) => {
      go.run(result.instance);
    });
</script>
```

4. Start a server (e.g. `python -m http.server`) and visit `http://localhost:8000` in your web browser to play your game.
