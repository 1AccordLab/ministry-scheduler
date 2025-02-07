# Setup

1. get into `app` directory (it's the main entrypoint for our dioxus app):

```bash
cd app/
```

2. compile tailwind css

```bash
npm install
npx @tailwindcss/cli -i ./input.css -o ./assets/tailwind.css
```

3. serve the dioxus app with both client & server:

```bash
cargo install dioxus-cli # install dioxus cli if you haven't already
dx serve
```

4. navigate to <http://localhost:8080> to open the app

## Todos

- oauth2 with line
  - [x] postgres: user table
  - [x] login
  - [] signup (persistence user data in postgres)
  - [x] logout
  - [] session/cookie
    - [x] in-memory
    - [] redis
- ministry scheduler
  - [] postgres: ministry schedule table
  - [] validate if a user can be assigned to a ministry event by series of conditions
    - [] only users belong to the ministry can be assigned
    - [] user must assigned once in single ministry event
    - [] user can schedule day-off before ministry events got scheduled
    - [] admin can determine the number of times a user can be assigned to a ministry event
    - [] admin can set the group of users who can be assigned to a ministry event
  - [] users can exchange their ministries with other users
    - [] request to specific user for help
    - [] request to all users for help
  - [] scheduler viewing
    - [] personal
    - [] by user
    - [] by ministry event
    - [] by time
  - [] websocket: group users by ministry event
  - [] websocket: allow collaboration editing & avoid race conditions
- line integration
  - [] notifications:
    - [] remind user for ministry event
    - [] ministry exchange
      - [] the user who is requested
      - [] require admin to check/approve
- advanced features
  - [] AI-integration
  - [] editing history
