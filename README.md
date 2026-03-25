# neomd

A minimal terminal email client for people who write in Markdown and live in Neovim.

[Neomd](https://www.ssp.sh/brain/neomd/) This is my way of implementing an email TUI based on my experience with Neomutt, focusing on [Neovim](https://www.ssp.sh/brain/neovim) (input) and reading/writing in [Markdown](https://www.ssp.sh/brain/markdown) and navigating with [Vim Motions](https://www.ssp.sh/brain/vim-language-and-motions).

![neomd](images/neomd.png)

Compose emails in your editor, read them rendered with [glamour](https://github.com/charmbracelet/glamour), and manage your inbox with a [HEY-style screener](https://www.hey.com/features/the-screener/) — all from the terminal. 

## Features

- **Write in Markdown, send beautifully** — compose in `$EDITOR` (defaults to `nvim`), send as `multipart/alternative`: raw Markdown as plain text + goldmark-rendered HTML so recipients get clickable links and formatting
- **Glamour reading** — incoming emails rendered as styled Markdown in the terminal
- **HEY-style screener** — unknown senders land in `ToScreen`; press `I/O/F/P` to approve, block, mark as Feed, or mark as PaperTrail; reuses your existing `screened_in.txt` lists from neomutt
- **Folder tabs** — Inbox, ToScreen, Feed, PaperTrail, Archive, Waiting, Someday, Scheduled, Sent, Trash, ScreenedOut
- **Multi-select** — `space` marks emails, then batch-delete, move, or screen them all at once
- **Kanagawa theme** — colors from the [kanagawa.nvim](https://github.com/rebelot/kanagawa.nvim) palette
- **IMAP + SMTP** — direct connection via RFC 6851 MOVE, no local sync daemon required

## Install

```sh
git clone https://github.com/sspaeti/neomd
cd neomd
make install   # installs to ~/.local/bin/neomd
```

Or just build locally:

```sh
make build
./neomd
```

## Configuration

On first run, neomd creates `~/.config/neomd/config.toml` with placeholders:

```toml
[[accounts]]
name     = "Personal"
imap     = "imap.example.com:993"   # :993 = TLS, :143 = STARTTLS
smtp     = "smtp.example.com:587"
user     = "me@example.com"
password = "app-password"
from     = "Me <me@example.com>"

# Multiple accounts supported — add more [[accounts]] blocks
# Switch between them with `ctrl+a` in the inbox

[screener]
# reuse your existing neomutt allowlist files
screened_in  = "~/.dotfiles/mutt/.lists/screened_in.txt"
screened_out = "~/.dotfiles/mutt/.lists/screened_out.txt"
feed         = "~/.dotfiles/mutt/.lists/feed.txt"
papertrail   = "~/.dotfiles/mutt/.lists/papertrail.txt"

[folders]
inbox        = "INBOX"
sent         = "Sent"
trash        = "Trash"
drafts       = "Drafts"
to_screen    = "ToScreen"
feed         = "Feed"
papertrail   = "PaperTrail"
screened_out = "ScreenedOut"
archive      = "Archive"
waiting      = "Waiting"
scheduled    = "Scheduled"
someday      = "Someday"

[ui]
theme       = "dark"   # dark | light | auto
inbox_count = 50
```

Use an app-specific password (Gmail, Fastmail, Hostpoint, etc.) rather than your main account password.

## Keybindings

Press `?` inside neomd to open the interactive help overlay. Start typing to filter shortcuts.

### Navigation

| Key | Action |
|-----|--------|
| `j` / `k` | Move down / up |
| `gg` | Jump to top |
| `G` | Jump to bottom |
| `enter` / `l` | Open email |
| `h` / `q` / `esc` | Back to inbox (from reader) |
| `?` | Toggle help overlay (type to filter) |

### Folders

| Key | Action |
|-----|--------|
| `L` / `tab` | Next folder tab |
| `H` / `shift+tab` | Previous folder tab |
| `gi` | Go to Inbox |
| `ga` | Go to Archive |
| `gf` | Go to Feed |
| `gp` | Go to PaperTrail |
| `gt` | Go to Trash |
| `gs` | Go to Sent |
| `gk` | Go to ToScreen |
| `go` | Go to ScreenedOut |
| `gw` | Go to Waiting |
| `gm` | Go to Someday |

### Multi-select & Batch operations

| Key | Action |
|-----|--------|
| `space` | Mark / unmark email + advance cursor |
| `U` | Clear all marks |
| `x` | Delete marked (or cursor) → Trash |
| `A` | Archive marked (or cursor) → Archive |

All screener and move actions below apply to **all marked emails**, or just the cursor email if nothing is marked.

### Screener (any folder)

| Key | Action |
|-----|--------|
| `I` | Approve sender → `screened_in.txt` + move to Inbox |
| `O` | Block sender → `screened_out.txt` + move to ScreenedOut |
| `F` | Mark as Feed → `feed.txt` + move to Feed |
| `P` | Mark as PaperTrail → `papertrail.txt` + move to PaperTrail |
| `S` | Dry-run screen Inbox (shows preview, then `y` to apply / `n` to cancel) |

### Move (no screener update)

| Key | Action |
|-----|--------|
| `Mi` | Move to Inbox |
| `Ma` | Move to Archive |
| `Mf` | Move to Feed |
| `Mp` | Move to PaperTrail |
| `Mt` | Move to Trash |
| `Mo` | Move to ScreenedOut |
| `Mw` | Move to Waiting |
| `Mm` | Move to Someday |

### Email actions

| Key | Action |
|-----|--------|
| `N` | Toggle read/unread (applies to marked or cursor) |
| `R` | Reload / refresh folder |
| `r` | Reply (from reader) |
| `c` | Compose new email |
| `e` | Open in `$EDITOR` (read-only) — search, copy, vim motions (from reader) |
| `O` | Open in browser — `$BROWSER` or `w3m` (from reader) |
| `ctrl+a` | Switch account (if multiple configured) |
| `/` | Filter emails |
| `q` | Quit |

### Composing

| Key | Action |
|-----|--------|
| `tab` / `enter` | Move to next field |
| `enter` (on Subject) | Open `$EDITOR` with a `.md` temp file |
| `esc` | Cancel |

After saving and closing the editor, the email is sent automatically.

## Screener Workflow

The screener classifies senders into four buckets using plain-text allowlists. Unknown senders land in `ToScreen` until you make a decision.

### How classification works

| List file | Category | Where email lands |
|-----------|----------|-------------------|
| `screened_in.txt` | Approved | stays in Inbox |
| `screened_out.txt` | Blocked | ScreenedOut |
| `feed.txt` | Newsletter / feed | Feed |
| `papertrail.txt` | Receipts / notifications | PaperTrail |
| *(not in any list)* | Unknown | ToScreen |

### Day-to-day: screen new arrivals

Press `S` (or run `:screen`) to dry-run the screener against the emails currently loaded in your Inbox. A preview shows what would move where — press `y` to apply, `n` to cancel.

For individual senders, use `I` / `O` / `F` / `P` from any folder or the ToScreen queue.

### Bulk re-classification after updating your lists

When you add many senders to `feed.txt` or `papertrail.txt` at once (e.g. after importing from HEY), use this workflow:

```
1. Edit feed.txt / papertrail.txt / screened_in.txt with the new senders
2. Restart neomd  (lists are loaded at startup)
3. :reset-toscreen   →  shows "Move N emails from ToScreen → Inbox? y/n"
                         (moves everything back so it can be re-classified)
4. y to confirm
5. :screen-all       →  dry-run against ALL inbox emails (not just the loaded subset)
6. y to apply
```

`:screen-all` (alias `:sa`) scans every email in your Inbox — read and unread — and proposes moves for any sender that is now in a list. It does **not** touch emails already in Feed, PaperTrail, or other folders.

### Colon commands

Press `:` to open the command line. Tab cycles through completions; Enter runs the command.

| Command | Alias | Description |
|---------|-------|-------------|
| `:screen` | `:s` | dry-run screen currently loaded Inbox emails |
| `:screen-all` | `:sa` | dry-run screen **every** Inbox email (no count limit) |
| `:reset-toscreen` | `:rts` | move all ToScreen emails back to Inbox |
| `:reload` | `:r` | reload the current folder |
| `:quit` | `:q` | quit neomd |

## How Sending Works

neomd sends every email as `multipart/alternative`:

- **`text/plain`** — the raw Markdown you wrote (readable as-is in any client)
- **`text/html`** — rendered by [goldmark](https://github.com/yuin/goldmark) with a clean CSS wrapper

This means recipients using Gmail, Apple Mail, Outlook, etc. see properly formatted links, bold, headers, and code blocks — while you write nothing but Markdown.

## Make Targets

```
make build    compile ./neomd
make run      build and run
make install  install to ~/.local/bin
make test     run tests
make vet      go vet
make fmt      gofmt -w .
make tidy     go mod tidy
make clean    remove compiled binary
make help     print this list
```

## Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — list, viewport, textinput components
- [Glamour](https://github.com/charmbracelet/glamour) — Markdown → terminal rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — styling
- [go-imap/v2](https://github.com/emersion/go-imap) — IMAP client (RFC 6851 MOVE)
- [go-message](https://github.com/emersion/go-message) — MIME parsing
- [goldmark](https://github.com/yuin/goldmark) — Markdown → HTML for sending
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — config parsing


## Inspirations

- [Neomutt](https://neomutt.org) — the gold standard terminal email client; neomd reuses its screener list format and borrows many keybindings
- [HEY](https://www.hey.com/features/the-screener/) — the Screener concept: unknown senders wait for a decision before reaching your inbox
- [hey-cli](https://github.com/sspaeti/hey-cli) — a Go CLI for HEY; provided the bubbletea patterns used here
- [newsboat](https://newsboat.org) — RSS reader whose `O` open-in-browser binding and vim navigation feel inspired neomd's reader view
- [emailmd.dev](https://www.emailmd.dev) — the idea that email should be written in Markdown
- [charmbracelet/pop](https://github.com/charmbracelet/pop) — minimal Go email sender from Charm
- [charmbracelet/glamour](https://github.com/charmbracelet/glamour) — Markdown rendering in the terminal
- [kanagawa.nvim](https://github.com/rebelot/kanagawa.nvim) — the color palette used for the inbox
- [msgvault](https://github.com/sspaeti/msgvault) — Go IMAP archiver; the IMAP client code in neomd is adapted from it

If you rather see the prompt, check out [initial prompt](_prompts/prompt.md) and its [plan](_prompts/prompt-plan.md).

## Roadmap

See at my second brain at [Roadmap](https://www.ssp.sh/brain/neomd#roadmap).
