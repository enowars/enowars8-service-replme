use actix::prelude::*;
use actix::{Actor, Addr, Context, Handler, StreamHandler};
use actix_web::web::{Bytes, BytesMut};
use actix_web_actors::ws;
// use futures::executor::block_on;
// use std::future::Future;
use std::io::{self, Read, Write};
use std::process::Command;
use std::task::Poll;
use std::time::{Duration, Instant};
use tokio::io::AsyncRead;
use tokio_pty_process::{
    AsyncPtyMaster, AsyncPtyMasterReadHalf, AsyncPtyMasterWriteHalf, Child, CommandExt,
};
use tokio_util::codec::{BytesCodec, Decoder, FramedRead};

const HEARTBEAT_INTERVAL: Duration = Duration::from_secs(5);
const CLIENT_TIMEOUT: Duration = Duration::from_secs(10);

struct AsyncPtyMasterReadHalfExt(AsyncPtyMasterReadHalf);

impl AsyncRead for AsyncPtyMasterReadHalfExt {
    fn poll_read(
        mut self: std::pin::Pin<&mut Self>,
        _cx: &mut std::task::Context<'_>,
        buf: &mut tokio::io::ReadBuf<'_>,
    ) -> std::task::Poll<std::io::Result<()>> {
        let dst = buf.initialize_unfilled();
        let r = self.0.read(dst);
        match r {
            Ok(n) => {
                buf.advance(n);
                Poll::Ready(Ok(()))
            }
            Err(err) if err.kind() == io::ErrorKind::WouldBlock => Poll::Pending,
            Err(err) => Poll::Ready(Err(err)),
        }
    }
}

// struct AsyncPtyMasterWriteHalfExt(AsyncPtyMasterWriteHalf);
//
// impl AsyncWrite for AsyncPtyMasterWriteHalfExt {
//     fn poll_write(
//         mut self: std::pin::Pin<&mut Self>,
//         cx: &mut std::task::Context<'_>,
//         buf: &[u8],
//     ) -> Poll<Result<usize, io::Error>> {
//         match self.0.write(buf) {
//             Ok(n) => Poll::Ready(Ok(n)),
//             Err(err) if err.kind() == io::ErrorKind::WouldBlock => Poll::Pending,
//             Err(err) => Poll::Ready(Err(err)),
//         }
//     }
//
//     fn poll_flush(
//         mut self: std::pin::Pin<&mut Self>,
//         cx: &mut std::task::Context<'_>,
//     ) -> Poll<Result<(), io::Error>> {
//         match self.0.flush() {
//             Ok(n) => Poll::Ready(Ok(())),
//             Err(err) if err.kind() == io::ErrorKind::WouldBlock => Poll::Pending,
//             Err(err) => Poll::Ready(Err(err)),
//         }
//     }
//
//     fn poll_shutdown(
//         mut self: std::pin::Pin<&mut Self>,
//         cx: &mut std::task::Context<'_>,
//     ) -> Poll<Result<(), io::Error>> {
//         Poll::Ready(Ok(()))
//     }
// }

// struct ChildExt(Child);
//
// impl Future for ChildExt {
//     type Output = Result<(), io::Error>;
//
//     fn poll(
//         mut self: std::pin::Pin<&mut Self>,
//         _cx: &mut std::task::Context<'_>,
//     ) -> Poll<Self::Output> {
//         let poll = self.0.poll_exit();
//         match poll {
//             Ok(result) => {
//                 if result.is_ready() {
//                     Poll::Ready(Ok(()))
//                 } else {
//                     Poll::Pending
//                 }
//             }
//             Err(err) => Poll::Ready(Err(err)),
//         }
//     }
// }

#[derive(Debug, Eq, PartialEq, Clone)]
pub struct IO(BytesMut);

impl Message for IO {
    type Result = ();
}

impl AsRef<[u8]> for IO {
    fn as_ref(&self) -> &[u8] {
        self.0.as_ref()
    }
}

impl Into<Bytes> for IO {
    fn into(self) -> actix_web::web::Bytes {
        self.0.into()
    }
}

impl From<actix_web::web::Bytes> for IO {
    fn from(b: actix_web::web::Bytes) -> Self {
        Self(b.as_ref().into())
    }
}

// impl From<String> for IO {
//     fn from(s: String) -> Self {
//         Self(s.into())
//     }
// }

impl From<&str> for IO {
    fn from(s: &str) -> Self {
        Self(s.into())
    }
}

pub struct ChildDied();

impl Message for ChildDied {
    type Result = ();
}

pub struct Websocket {
    command: Option<Command>,
    terminal: Option<Addr<Terminal>>,
    hb: Instant,
}

impl Websocket {
    pub fn new(command: Command) -> Self {
        Self {
            hb: Instant::now(),
            terminal: None,
            command: Some(command),
        }
    }

    fn hb(&self, ctx: &mut <Self as Actor>::Context) {
        ctx.run_interval(HEARTBEAT_INTERVAL, |act, ctx| {
            if Instant::now().duration_since(act.hb) > CLIENT_TIMEOUT {
                log::warn!("Client heartbeat timeout, disconnecting.");
                ctx.stop();
                return;
            }

            ctx.ping(b"");
        });
    }
}

impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for Websocket {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        let cons: &mut Addr<Terminal> = match self.terminal {
            Some(ref mut c) => c,
            None => {
                log::error!("Terminalole died, closing websocket.");
                ctx.stop();
                return;
            }
        };

        match msg {
            Ok(msg) => match msg {
                ws::Message::Ping(msg) => {
                    self.hb = Instant::now();
                    ctx.pong(&msg);
                }
                ws::Message::Pong(_) => self.hb = Instant::now(),
                ws::Message::Text(t) => {
                    cons.do_send(IO::from(t.as_ref()));
                }
                ws::Message::Binary(b) => cons.do_send(IO::from(b)),
                ws::Message::Close(_) => ctx.stop(),
                ws::Message::Nop => {}
                ws::Message::Continuation(_cont) => {
                    log::warn!("No support for continuation frames");
                }
            },
            Err(_) => todo!(),
        };
    }
}

impl Actor for Websocket {
    type Context = ws::WebsocketContext<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        // Start heartbeat
        self.hb(ctx);

        let command = self
            .command
            .take()
            .expect("command was None at start of WebSocket.");

        // Start PTY
        self.terminal = Some(Terminal::new(ctx.address(), command).start());

        log::trace!("Started WebSocket");
    }

    fn stopping(&mut self, _ctx: &mut Self::Context) -> Running {
        log::trace!("Stopping WebSocket");

        Running::Stop
    }

    fn stopped(&mut self, _ctx: &mut Self::Context) {
        log::trace!("Stopped WebSocket");
    }
}

impl Handler<IO> for Websocket {
    type Result = ();

    fn handle(&mut self, msg: IO, ctx: &mut <Self as Actor>::Context) {
        ctx.binary(msg);
    }
}

impl Handler<ChildDied> for Websocket {
    type Result = ();

    fn handle(&mut self, _msg: ChildDied, ctx: &mut <Self as Actor>::Context) {
        log::trace!("Websocket <- ChildDied");
        ctx.close(None);
        ctx.stop();
    }
}

pub struct Terminal {
    ws: Addr<Websocket>,
    command: Command,
    stdin: Option<AsyncPtyMasterWriteHalf>,
    child: Option<Child>,
}

impl Terminal {
    pub fn new(ws: Addr<Websocket>, command: Command) -> Self {
        Self {
            ws,
            command,
            stdin: None,
            child: None,
        }
    }
}

impl StreamHandler<Result<<BytesCodec as Decoder>::Item, <BytesCodec as Decoder>::Error>>
    for Terminal
{
    fn handle(
        &mut self,
        msg: Result<<BytesCodec as Decoder>::Item, <BytesCodec as Decoder>::Error>,
        _ctx: &mut Self::Context,
    ) {
        if let Ok(bytes) = msg {
            self.ws.do_send(IO(bytes));
        }
    }
}

impl Actor for Terminal {
    type Context = Context<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        let pty = match AsyncPtyMaster::open() {
            Err(e) => {
                log::error!("Unable to open PTY: {:?}", e);
                ctx.stop();
                return;
            }
            Ok(pty) => pty,
        };

        let child = match self.command.spawn_pty_async(&pty) {
            Err(e) => {
                log::error!("Unable to spawn child: {:?}", e);
                ctx.stop();
                return;
            }
            Ok(child) => child,
        };

        log::info!("Spawned new child process with PID {}", child.id());

        let (pty_read, pty_write) = pty.split();

        self.stdin = Some(pty_write);
        self.child = Some(child);

        let stdout = AsyncPtyMasterReadHalfExt(pty_read);

        Self::add_stream(FramedRead::new(stdout, BytesCodec::new()), ctx);
    }

    fn stopping(&mut self, _ctx: &mut Self::Context) -> Running {
        log::info!("Stopping Terminal");

        let child = self.child.take();
        if child.is_none() {
            return Running::Stop;
        }

        let mut child = child.unwrap();
        match child.kill() {
            Ok(()) => log::info!("Child died"),
            // Ok(()) => match block_on(ChildExt(child)) {
            //     Ok(_) => log::info!("Child died"),
            //     Err(e) => log::error!("Child wouldn't die: {}", e),
            // },
            Err(e) => log::error!("Could not kill child with PID {}: {}", child.id(), e),
        };
        self.ws.do_send(ChildDied());
        Running::Stop
    }

    fn stopped(&mut self, _ctx: &mut Self::Context) {
        log::info!("Stopped Terminal");
    }
}

impl Handler<IO> for Terminal {
    type Result = ();

    fn handle(&mut self, msg: IO, ctx: &mut <Self as Actor>::Context) {
        let pty = match self.stdin {
            Some(ref mut p) => p,
            None => {
                log::error!("Write half of PTY died, stopping Terminal.");
                ctx.stop();
                return;
            }
        };

        if let Err(e) = pty.write(msg.as_ref()) {
            log::error!("Could not write to PTY: {}", e);
            ctx.stop();
        }

        log::trace!("Websocket -> Terminal : {:?}", msg);
    }
}
