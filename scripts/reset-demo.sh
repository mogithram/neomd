#!/bin/bash
# Reset neomd demo to first-run state.
# Only touches demo-specific files — never modifies production config or cache.
# Usage: ./scripts/reset-demo.sh [config-dir]

set -e

CONFIG_DIR="${1:-$HOME/.config/neomd-demo}"
DEMO_CACHE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/neomd-demo"

echo "Resetting neomd demo state..."
echo "  Config dir: $CONFIG_DIR"
echo "  Cache dir:  $DEMO_CACHE_DIR"
echo

# 1. Remove demo welcome marker
if [ -f "$DEMO_CACHE_DIR/welcome-shown" ]; then
    rm -f "$DEMO_CACHE_DIR/welcome-shown"
    echo "  [x] Removed welcome marker"
else
    echo "  [ ] Welcome marker already absent"
fi

# 2. Remove demo screener lists directory (recreated on next launch)
if [ -d "$CONFIG_DIR/lists" ]; then
    rm -rf "$CONFIG_DIR/lists"
    echo "  [x] Removed lists directory"
else
    echo "  [ ] Lists directory already absent"
fi

# 3. Clear demo command history
if [ -f "$DEMO_CACHE_DIR/cmd_history" ]; then
    rm -f "$DEMO_CACHE_DIR/cmd_history"
    echo "  [x] Cleared command history"
fi

echo
echo "Done! Next launch will show the welcome screen."
echo "Run: make demo"
