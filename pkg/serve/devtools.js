(function () {
  const socket = new WebSocket("ws://" + location.host + "/ws");

  socket.onmessage = function (event) {
    try {
      const message = JSON.parse(event.data);
      switch (message.type) {
        case "reload":
          location.reload();
          break;
        case "error":
          console.error("Build error:", message.data);
          showErrorOverlay(message.data);
          break;
        case "css-update":
          console.log("CSS update:", message.data.path);
          updateCSS(message.data.path);
          break;
        default:
          console.log("Unknown message type:", message.type);
      }
    } catch (e) {
      console.error("Invalid message from server:", event.data);
    }
  };

  socket.onclose = function () {
    console.log("Connection closed. Trying to reconnect...");
    setTimeout(function () {
      // Attempt to reconnect
      location.reload();
    }, 2000);
  };

  socket.onerror = function (error) {
    console.error("WebSocket error:", error);
    socket.close();
  };

  function showErrorOverlay(error) {
    let overlay = document.getElementById("evoke-error-overlay");
    if (!overlay) {
      overlay = document.createElement("div");
      overlay.id = "evoke-error-overlay";
      document.body.appendChild(overlay);
    }

    const parsedError = parseError(error);

    overlay.innerHTML = `
      <style>
        #evoke-error-overlay {
          position: fixed;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          background-color: rgba(0, 0, 0, 0.8);
          color: #e8e8e8;
          z-index: 9999;
          display: flex;
          justify-content: center;
          align-items: center;
          font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
        }
        .evoke-error-modal {
          background: #2a2a2a;
          border: 1px solid #444;
          border-radius: 8px;
          padding: 2rem;
          max-width: 80%;
          max-height: 80%;
          overflow: auto;
          box-shadow: 0 10px 30px rgba(0,0,0,0.5);
        }
        .evoke-error-modal h2 {
          color: #ff6b6b;
          margin-top: 0;
          border-bottom: 1px solid #444;
          padding-bottom: 1rem;
        }
        .evoke-error-modal .error-details {
          margin-bottom: 1rem;
        }
        .evoke-error-modal .error-details strong {
          color: #ffb86c;
        }
        .evoke-error-modal pre {
          white-space: pre-wrap;
          word-wrap: break-word;
          background: #1e1e1e;
          padding: 1rem;
          border-radius: 5px;
        }
        .evoke-error-modal button {
          background: #ff6b6b;
          color: white;
          border: none;
          padding: 0.8rem 1.5rem;
          border-radius: 5px;
          cursor: pointer;
          font-size: 1rem;
          margin-top: 1rem;
        }
      </style>
      <div class="evoke-error-modal">
        <h2>Build Error</h2>
        <div class="error-details">
          <p><strong>File:</strong> ${parsedError.file}</p>
          <p><strong>Line:</strong> ${parsedError.line}</p>
          <p><strong>Column:</strong> ${parsedError.column}</p>
        </div>
        <pre>${parsedError.description}</pre>
        <button onclick="document.getElementById('evoke-error-overlay').remove()">Close</button>
      </div>
    `;
  }

  function parseError(error) {
    const match = error.match(
      /template: (.*?):(\d+):(\d+): executing ".*" at <(.*?)>: (.*)/
    );
    if (match) {
      return {
        file: match[1],
        line: match[2],
        column: match[3],
        template: match[4],
        description: match[5],
      };
    }
    return {
      file: "N/A",
      line: "N/A",
      column: "N/A",
      description: error,
    };
  }

  function updateCSS(path) {
    const links = document.querySelectorAll(
      `link[rel="stylesheet"][href^="${path}"]`
    );
    links.forEach((link) => {
      const url = new URL(link.href);
      url.searchParams.set("v", new Date().getTime());
      link.href = url.href;
    });
  }
})();
