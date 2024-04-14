(function refresh () {
  const verboseLogging = false;

  let socketUrl = window.location.origin;
  if (!window.location.origin.match(/:[0-9]+/)) {
    socketUrl = window.location.origin + ":80";
  }

  socketUrl = socketUrl.replace(/^https?:\/\/(.+):(\d+)/, "ws://$1:9765"); 
  let socket;

  function dprint(...msg) {
    if (verboseLogging) console.debug("[kreload] ", ...msg);
  }
  function eprint(...msg) {
    if (verboseLogging) console.error("[kreload] ", ...msg);
  }

  dprint("reload script loaded");

  if (!("WebSocket" in window)) {
    throw new Error("WebSocket not supported in this browser");
  }

  let firstConnection = true;
  
  window.addEventListener("load", function () {
    websocketStart();
  });

  window.addEventListener("beforeunload", function () {
    dprint("page unloading");
  });

  const socketOnMessage = function (msg) {
    if (msg.data === "reload") {
      dprint("reloading...");
      socket.close();
      window.location.reload();
    }
  };

  const socketOnOpen = function (_msg) {
    dprint("ws connected");

    
    if (firstConnection === true) {
      // the page should not reload, if this is the first connection
    } else {
      dprint("conection recovered, reloading...");
      window.location.reload();
    }
  };

  const socketOnClose = function (_msg) {
    dprint("ws closed");
    firstConnection = false;
    websocketStart();
  };

  const socketOnError = function (msg) {
    eprint(msg);
  };

  /** Start WebSocket connection */
  function websocketStart () {
    dprint("starting ws connection...");
    
    setTimeout(function () {
      try {
        socket = new WebSocket(socketUrl);

        socket.onopen = socketOnOpen;
        socket.onclose = socketOnClose;
        socket.onmessage = socketOnMessage;
        socket.onerror = socketOnError;
      }
      catch (err) {
        eprint("connection failed. retrying...");
      }
    }, 250);
  }
})();