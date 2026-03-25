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

> The tables below are generated from [`internal/ui/keys.go`](internal/ui/keys.go).
> To update both the help overlay and this section at once, edit that file and run `make docs`.

<!-- keybindings-start -->

### Navigation

| Key | Action |
|-----|--------|
| `j / k` | move down / up |
| `gg` | jump to top |
| `G` | jump to bottom |
| `enter / l` | open email |
| `h / q / esc` | back to inbox (from reader) |
| `?` | toggle help overlay (type to filter) |


### Folders

| Key | Action |
|-----|--------|
| `L / tab` | next folder tab |
| `H / shift+tab` | previous folder tab |
| `gi` | go to Inbox |
| `ga` | go to Archive |
| `gf` | go to Feed |
| `gp` | go to PaperTrail |
| `gt` | go to Trash |
| `gs` | go to Sent |
| `gk` | go to ToScreen |
| `go` | go to ScreenedOut |
| `gw` | go to Waiting |
| `gm` | go to Someday |


### Screener  (marked or cursor, any folder)

| Key | Action |
|-----|--------|
| `I` | approve sender → screened_in.txt + move to Inbox |
| `O` | block sender → screened_out.txt + move to ScreenedOut |
| `F` | mark as Feed → feed.txt + move to Feed |
| `P` | mark as PaperTrail → papertrail.txt + move to PaperTrail |
| `A` | archive (move to Archive, no screener update) |
| `S` | dry-run screen inbox (loaded emails), then y/n |


### Move  (marked or cursor, no screener update)

| Key | Action |
|-----|--------|
| `x` | delete → Trash |
| `Mi` | move to Inbox |
| `Ma` | move to Archive |
| `Mf` | move to Feed |
| `Mp` | move to PaperTrail |
| `Mt` | move to Trash |
| `Mo` | move to ScreenedOut |
| `Mw` | move to Waiting |
| `Mm` | move to Someday |


### Multi-select

| Key | Action |
|-----|--------|
| `space` | mark / unmark email + advance cursor |
| `U` | clear all marks |


### Sort  (, prefix)

| Key | Action |
|-----|--------|
| `,m` | date newest first (default) |
| `,M` | date oldest first |
| `,a` | from A→Z |
| `,A` | from Z→A |
| `,s` | size smallest first |
| `,S` | size largest first |
| `,n` | subject A→Z |
| `,N` | subject Z→A |


### Email actions

| Key | Action |
|-----|--------|
| `N` | toggle read/unread  (marked or cursor) |
| `ctrl+n` | mark all in current folder as read |
| `R` | reload / refresh folder |
| `r` | reply  (from reader) |
| `c` | compose new email |
| `e  (reader)` | open in $EDITOR read-only — search, copy, vim motions |
| `O  (reader)` | open in browser — $BROWSER or w3m |
| `ctrl+a` | switch account  (if multiple configured) |


### Command line  (: to open, tab to complete)

| Key | Action |
|-----|--------|
| `:screen  / :s` | dry-run screen loaded inbox emails |
| `:screen-all  / :sa` | dry-run screen ALL inbox emails (no limit) |
| `:reset-toscreen  / :rts` | move all ToScreen emails back to Inbox |
| `:mark-read  / :mr` | mark all emails in current folder as read |
| `:reload  / :r` | reload current folder |
| `:check  / :ch` | show screener classification for selected email |
| `:delete-all  / :da` | permanently delete ALL emails in current folder (y/n) |
| `:create-folders  / :cf` | create missing IMAP folders from config (safe, idempotent) |
| `:quit  / :q` | quit neomd |


### Composing

| Key | Action |
|-----|--------|
| `tab / enter` | move to next field |
| `enter  (on Subject)` | open $EDITOR with a .md temp file |
| `esc` | cancel |


### General

| Key | Action |
|-----|--------|
| `/` | filter emails in current folder |
| `?` | toggle this help |
| `q` | quit  (from inbox) |

<!-- keybindings-end -->

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
