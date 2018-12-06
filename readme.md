Container based inline terminal demos for Google Slides

### Usage

- Open `chrome://extensions/` and load directory `chrome_extension`
- Make sure `docker` is installed
- `presentty --config ./examples/docker.toml`
- In Google Slides draw a shape and add a link `#term=dockerdemo`
- When presenting, on clicking the link, the area of the shape will be replaced with a demo terminal with a separate Docker daemon running


### Extra features

- `#term=id,autostart` will inject the terminal automatically on load
- There is a small gray target area in the bottom right corner that can be used the reload the terminal without refreshing the page if needed.
- `presentty --asciinema` runs demos through asciinema (needs to be installed on the host) and records them to the local directory.
- After uploading asciinema recordings you can replace the links with `https://asciinema.org/...#term=...` so that they now point to recordings for people who have not installed the extension.
- If UI fails on live the demos can be started with regular `docker exec` while `presentty` is running. If all fails `asciinema` recordings can be replayed with `asciinema play`.
- It's recommended to use `tmux a` as an entrypoint of the command so you can continue from the place you left off if you accidentally switch slides. Look at the example demo for instructions.