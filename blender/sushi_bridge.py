"""
Sushi Bridge — Blender addon
Receives import requests from the Sushi asset manager over HTTP.

Install:
  1. Open Blender → Edit → Preferences → Add-ons → Install…
  2. Select this file (sushi_bridge.py)
  3. Enable "Import: Sushi Bridge"

The addon starts a tiny HTTP server on 127.0.0.1:29877 that listens for
import requests from Sushi. All communication is local-only.
"""

bl_info = {
    "name": "Sushi Bridge",
    "author": "sushi",
    "version": (0, 1, 0),
    "blender": (3, 0, 0),
    "location": "Runs in background",
    "description": "Receives 3D assets from the Sushi asset manager",
    "category": "Import-Export",
}

import bpy
import json
import os
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler
from functools import partial

SUSHI_PORT = 29877

# ── HTTP Handler ──────────────────────────────────────────────────────

class SushiHandler(BaseHTTPRequestHandler):
    """Handles incoming requests from the Sushi desktop app."""

    def log_message(self, format, *args):
        """Suppress default stderr logging."""
        pass

    def do_GET(self):
        """Health check / ping."""
        if self.path == "/sushi":
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps({"status": "ok", "app": "sushi-bridge"}).encode())
        else:
            self.send_response(404)
            self.end_headers()

    def do_POST(self):
        """Receive import requests."""
        if self.path != "/sushi":
            self.send_response(404)
            self.end_headers()
            return

        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)

        try:
            data = json.loads(body)
        except json.JSONDecodeError:
            self.send_response(400)
            self.end_headers()
            self.wfile.write(b'{"error": "invalid json"}')
            return

        action = data.get("action")
        files = data.get("files", [])

        if action == "import" and files:
            # Schedule the import on Blender's main thread via a timer
            bpy.app.timers.register(partial(_import_files, files))
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps({
                "status": "ok",
                "queued": len(files),
            }).encode())
        else:
            self.send_response(400)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps({"error": "unknown action or no files"}).encode())


# ── Import Logic ──────────────────────────────────────────────────────

def _import_files(file_paths):
    """Called on Blender's main thread to import the given files."""
    for path in file_paths:
        path = os.path.abspath(path)
        if not os.path.isfile(path):
            print(f"[sushi] file not found: {path}")
            continue

        ext = os.path.splitext(path)[1].lower()
        try:
            if ext == ".glb" or ext == ".gltf":
                bpy.ops.import_scene.gltf(filepath=path)
                print(f"[sushi] imported: {path}")
            else:
                print(f"[sushi] unsupported format: {ext}")
        except Exception as e:
            print(f"[sushi] import error for {path}: {e}")

    return None  # returning None unregisters the timer (one-shot)


# ── Server Lifecycle ──────────────────────────────────────────────────

_server = None
_server_thread = None


def _start_server():
    global _server, _server_thread

    if _server is not None:
        return  # already running

    try:
        _server = HTTPServer(("127.0.0.1", SUSHI_PORT), SushiHandler)
        _server.timeout = 1
        _server_thread = threading.Thread(target=_server.serve_forever, daemon=True)
        _server_thread.start()
        print(f"[sushi] bridge listening on 127.0.0.1:{SUSHI_PORT}")
    except OSError as e:
        print(f"[sushi] could not start bridge: {e}")
        _server = None


def _stop_server():
    global _server, _server_thread

    if _server is not None:
        _server.shutdown()
        _server = None
        _server_thread = None
        print("[sushi] bridge stopped")


# ── Blender Registration ─────────────────────────────────────────────

def register():
    _start_server()


def unregister():
    _stop_server()


if __name__ == "__main__":
    register()
