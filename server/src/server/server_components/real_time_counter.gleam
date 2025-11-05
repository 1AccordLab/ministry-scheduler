import client/components/counter
import ewe.{
  type Request, type Response, type WebsocketConnection, type WebsocketMessage,
  type WebsocketNext,
}
import gleam/erlang/process.{type Subject}
import gleam/json
import lustre
import lustre/element.{type Element}
import lustre/server_component.{
  type ClientMessage, client_message_to_json, runtime_message_decoder,
}

pub fn serve(req: Request) -> Response {
  ewe.upgrade_websocket(req, on_init, handler, on_close)
}

pub fn element() -> Element(_) {
  server_component.element([server_component.route("/ws/counter")], [])
}

type CounterSocket {
  CounterSocket(
    self: Subject(CounterMsg),
    component: lustre.Runtime(counter.Msg),
  )
}

type CounterMsg =
  ClientMessage(counter.Msg)

fn on_init(
  _connection: WebsocketConnection,
  selector: process.Selector(CounterMsg),
) {
  // start the server component runtime on the server-side
  let assert Ok(component) =
    lustre.start_server_component(counter.component(), Nil)

  // client <-> websocket server(`process.subject`) <-> server component runtime
  let self = process.new_subject()

  // `handler()` method will handle messages received from client/server component runtime
  let selector = process.select(selector, self)

  // establish the connection between:
  // `component`: client
  // `self`: server component runtime
  server_component.register_subject(self)
  |> lustre.send(to: component)

  #(CounterSocket(self, component), selector)
}

fn handler(
  connection: WebsocketConnection,
  socket: CounterSocket,
  msg: WebsocketMessage(CounterMsg),
) -> WebsocketNext(CounterSocket, CounterMsg) {
  case msg {
    // the websocket server receives a message from the client first,
    // decode it, and pass to the server component runtime
    ewe.Text(json) -> {
      echo "message received from client: " <> json

      // send decoded message to server component runtime
      let assert Ok(runtime_msg) = json.parse(json, runtime_message_decoder())
      lustre.send(socket.component, runtime_msg)

      ewe.websocket_continue(socket)
    }

    ewe.Binary(_) -> {
      ewe.websocket_continue(socket)
    }

    // the server component sends a message to the websocket server first,
    // and forward the message to the client
    ewe.User(client_msg) -> {
      let text_frame =
        client_message_to_json(client_msg)
        |> json.to_string
      let assert Ok(_) = ewe.send_text_frame(connection, text_frame)
      echo "message sent to client: " <> text_frame

      ewe.websocket_continue(socket)
    }
  }
}

fn on_close(_connection: WebsocketConnection, socket: CounterSocket) {
  server_component.deregister_subject(socket.self)
  |> lustre.send(to: socket.component)
}
