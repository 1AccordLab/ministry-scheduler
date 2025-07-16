# Update dependencies
update-deps:
  cd client && gleam update
  cd server && gleam update

# Rebuild and restart the app on change
dev:
  # Create this directory so watchexec can detect it during the first `just dev` run
  mkdir -p server/priv/static

  # Watch frontend for changes and rebuild on update
  watchexec -w client/src \
            --clear --restart \
            --wrap-process=session \
            --stop-signal=SIGKILL \
            just build &

  # Watch server and frontend built static assets, restart on update
  watchexec -w server/src \
            -w server/priv/static --no-vcs-ignore \
            --clear --restart \
            --wrap-process=session \
            --stop-signal=SIGKILL \
            just run &

  wait

# Build frontend and output to server static dir
build:
  cd client && gleam run -m lustre/dev build --minify --outdir=../server/priv/static

# Run backend server
run:
  cd server && gleam run

# Stop all running watchexec processes
stop:
  pkill -x watchexec

