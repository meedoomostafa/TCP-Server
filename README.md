# Go TCP Chat Server

A lightweight, concurrent TCP chat server written in Go. This project demonstrates low-level network programming using Go’s `net` package, implementing the **handle-per-connection** concurrency pattern and `sync.Mutex` for thread safety.

---

## Features

- **Concurrent Connections**  
  Handles multiple TCP clients simultaneously using lightweight goroutines.

- **Bi-Directional Communication**
  - **Client → Server**: Logs incoming messages from clients with their remote address.
  - **Server → Client**: Supports broadcasting messages from the server terminal to all connected clients.

- **Thread Safety**  
  Uses `sync.Mutex` to safely manage the map of active peers (clients).

- **Channel-Based Architecture**  
  Decouples message reading from message processing.

---

## 🛠️ Architecture

The server uses a **fan-in / fan-out** pattern for message handling and mutex locking for peer management.

```mermaid
sequenceDiagram
    participant Admin as Server Terminal (Stdin)
    participant Server as Go Server (Main)
    participant Client as TCP Client (Telnet)

    Note over Server: Start Listener :3000

    Client->>Server: Connect (TCP Handshake)
    Server->>Client: Accept Connection & Spawn ReadLoop
    
    par Client Sending
        Client->>Server: "Hello"
        Server->>Server: ReadLoop -> msgch
        Server->>Admin: Log: "Client says: Hello"
    and Server Broadcasting
        Admin->>Server: Type "Welcome!" (InputLoop)
        Server->>Server: Lock Peers Map
        loop Every Connected Client
            Server->>Client: Write "Server: Welcome!"
        end
        Server->>Server: Unlock Peers Map
    end
````

---

## 📦 Getting Started

### Prerequisites

* Go **1.20.0+**
* Telnet or Netcat (for testing clients)

---

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/go-tcp-server.git
cd TCP-Server
```

Run the server:

```bash
go run main.go
```

**Output**

```
Server started on :3000. Type to broadcast messages.
```

---

## 🎮 How to Use

### 1. Start the Server

Run the Go program in your main terminal.
This terminal acts as the **Server Admin Console**.

### 2. Connect a Client

Open a new terminal window and connect via Telnet:

```bash
telnet localhost 3000
```

> **Note:** Use a space, not a colon, between `localhost` and `3000`.

### 3. Chat

* **From Client**: Type in the Telnet window and hit Enter.
  The server logs the message.
* **From Server**: Type in the Server terminal and hit Enter.
  All connected clients receive the message.

---

## 📂 Code Structure

* **NewServer**
  Initializes the server struct, channels, and the peer map.

* **Start**
  Bootstraps the listener, the `acceptLoop`, and the `inputLoop`.

* **ReadLoop**
  A dedicated goroutine spawned for every client to read incoming data.

* **inputLoop**
  Reads from server `stdin` and broadcasts messages to all active peers.

---

## 🐛 Common Issues

* **Connection Refused**
  Ensure the server is running before connecting via Telnet.

* **Telnet Syntax**
  Use:

  ```bash
  telnet localhost 3000
  ```

  ❌ Do not use `localhost:3000`.

* **Output Formatting**
  The server uses `fmt.Printf` for clean logs; ensure your terminal supports standard output.

---

