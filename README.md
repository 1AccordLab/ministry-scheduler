# Setup

1. get into `app` directory (it's the main entrypoint for our dioxus app):

```bash
cd app/
```

2. compile tailwind css

```bash
npx @tailwindcss/cli -i ./input.css -o ./assets/tailwind.css
```

3. serve the dioxus app with both client & server:

```bash
cargo install dioxus-cli # install dioxus cli if you haven't already
dx serve
```

4. navigate to <http://localhost:8080> to open the app
