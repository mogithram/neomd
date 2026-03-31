# Changelog

## 2026-03-31
- **Threaded inbox** — related emails are automatically grouped in the inbox list with a Twitter-style vertical connector line (`│`/`╰`); threads detected via `In-Reply-To`/`Message-ID` IMAP envelope headers with a normalized-subject + participant overlap fallback; newest reply on top, root at bottom; threads sorted by most recent email so active conversations float to the top; thread connector colored subtly in the Kanagawa border color
- fix showing recipient in SENT tab (instead of from)
- **IMAP search across all folders (`space /` or `:search`)** — server-side IMAP SEARCH across all configured folders (Inbox, Sent, Archive, Feed, etc.); results displayed in a temporary "Search" tab with `[Folder]` prefix on each subject; supports query prefixes: `from:simon`, `subject:invoice`, `to:team@`, or plain text to search all three fields; press `esc` to close results
- **Filter preserves across actions** — the local `/` filter no longer clears when pressing `n` (toggle read), `m` (mark), `U` (clear marks), or sorting; filter stays active until `esc`
- **Address autocomplete in compose** — To, Cc, and Bcc fields show autocomplete suggestions from screener lists (`screened_in.txt`, `feed.txt`, `papertrail.txt`); navigate with `ctrl+n`/`ctrl+p`/arrows, accept with `tab`; supports multi-address fields (autocomplete applies after the last comma)
- **Everything view (`ge` or `:everything`)** — shows the 50 most recent emails across all folders in a temporary "Everything" tab, sorted by date descending; each subject prefixed with `[Folder]`; useful for finding emails that were screened out or moved to spam
- **Link opener (`space+1-9` in reader)** — links are extracted from the email body, numbered `[1]`-`[0]` in the header; press `space` then a digit to open in `$BROWSER`; up to 10 links per email, deduplicated by URL
- **Draft signature fix** — re-opening a draft (`E`) no longer appends a duplicate signature; the draft body already contains it from the first compose
- **Draft reader footer** — `E draft` now appears in the reader footer when viewing an email from the Drafts folder
- **Android support (`make android`)** — cross-compile for Android ARM64; runs in Termux; documented in `docs/android.md` with install instructions and useful shortcuts
- **Docs restructure** — detailed documentation moved from README to `docs/` folder: `docs/keybindings.md` (auto-generated), `docs/screener.md`, `docs/sending.md`, `docs/configuration.md`, `docs/android.md`; README kept concise with links

## 2026-03-30

- added  preview email in $BROWSER (images rendered, same as recipient sees)  with `p`
- **Multiple From addresses / SMTP aliases** — add `[[senders]]` blocks to config to define extra From identities (e.g. `s@ssp.sh` as an alias through an existing account's SMTP); cycle through all accounts + senders with `ctrl+f` in both compose and pre-send screens; the `account =` field matches by account `name =` (not email address)
- **Sent folder** — after sending, neomd APPENDs a copy to the configured Sent IMAP folder with `\Seen` flag; the same raw MIME bytes used for SMTP delivery are reused for the APPEND (no double-build)
- **Attachment column in inbox** — `@` appears in a dedicated column next to the date when an email has attachments (detected from IMAP BODYSTRUCTURE including inline images)
- **Attachment downloads in reader** — the email header now lists all attachments as `[1] report.pdf  [2] photo.png`; press `1`–`9` to download attachment N to `~/Downloads/` and open it with `xdg-open`; filenames are deduplicated automatically
- **Inline images as downloads** — images embedded inline in emails (`Content-Disposition: inline`, e.g. PNG screenshots) are now shown alongside regular attachments in the reader header and downloadable with `1`–`9`; previously only `Content-Disposition: attachment` parts were listed
- **Inline image placeholders in reader body** — `<img src="cid:...">` tags now show `[Image: filename.png]` at their position in the body text instead of being silently stripped; uses Content-ID → filename mapping from MIME parts
- **Undo move / delete** — `u` reverses the last single or batch move/delete (`x`, `A`, `M*`); uses the UIDPLUS destination UID so undo still works even when the server reassigns UIDs on MOVE; screener actions (`I`, `O`, `F`, `P`, `$`) are intentionally excluded because they also modify `.txt` list files
- **Subject (and headers) re-parsed from editor** — editing `# [neomd: subject: ...]`, `# [neomd: to: ...]`, etc. in neovim now correctly updates those fields; previously the values were captured in a closure before the editor opened and changes were silently discarded; all three editor entry points (new compose, reply, continue draft) now call `editor.ParseHeaders` on the saved file content
- **`ctrl+f` for cycling From** — changed from `f` (which conflicts with typing in text fields) to `ctrl+f`; works in both the compose form and the pre-send review screen
- **Forward (`f`)** — forward an email from the reader or inbox; opens the editor with the original message quoted, `Fwd:` subject prefix, and empty `To:` field; from inbox the body is fetched automatically before opening the editor
- **Permanent delete (`X`, Trash only)** — permanently deletes marked or cursor email(s) from the Trash folder via IMAP STORE `\Deleted` + UID EXPUNGE; blocked in other folders with a warning message
- **`:empty-trash` / `:et`** — permanently delete all emails in Trash with y/n confirmation; works from any folder without navigating to Trash first
- **First-run welcome popup** — on the very first launch, a centered popup shows quick-start keybindings and screener basics; any key dismisses it; marker at `~/.cache/neomd/welcome-shown` ensures it only appears once
- **Auto-create IMAP folders on startup** — `ensureFoldersCmd` runs during `Init()` so new users don't need to manually run `:create-folders`; idempotent for existing users
- **Auto-create screener list directories** — parent directories for screener list paths are created automatically during config load; prevents errors when pressing `I`/`O`/`F`/`P` on a fresh install
- **Default screener paths** — changed from `~/.config/mutt/` to `~/.config/neomd/lists/` for new installs; existing configs with custom paths are unaffected
- **Go prerequisite check in Makefile** — `make build`/`make install` now prints clear Go installation instructions instead of a cryptic error when `go` is not found
- **Pre-send preview (`p`)** — press `p` in the pre-send screen to open a browser preview of the composed email; renders through the same goldmark pipeline as sending, with local image paths converted to `file://` URLs so inline images from `[attach]` lines display correctly

## 2026-03-29

- **CC field** — compose and reply forms now include an optional Cc field (Tab/Enter to skip); CC recipients receive the email and appear in the `Cc:` header
- **BCC field** — hidden by default; toggle with `ctrl+b` in compose; BCC recipients receive the email but are not visible in the message headers (standard BCC privacy)
- **Reply-all** — `R` in the reader replies to the original sender + all CC recipients; your own address is excluded automatically; uses `Reply-To` header when present
- **Pre-send review screen** — after closing the editor, neomd shows a summary (To, Subject, body preview) before sending; press `enter` to send, `a` to attach files via yazi (auto-detected, no config needed; override with `$NEOMD_FILE_PICKER`), `D` to remove last attachment, `d` to save to Drafts, `e` to re-open the editor, `esc` to cancel; avoids tmux/terminal key-capture issues since `a` needs no modifier
- **Save to Drafts** — `d` in the pre-send screen APPENDs the composed message to the configured Drafts IMAP folder with `\Draft` + `\Seen` flags; navigate to it with `gd`
- **Attachments from neovim** — `<leader>a` in a `neomd-*.md` buffer opens yazi in a floating terminal; selected files are inserted as `[attach] /path/to/file` lines (visible in markdown, not hidden HTML comments); neomd strips them before sending and adds them as MIME attachments
- **Inline code and code blocks** — `` `inline code` `` and fenced ` ``` ` blocks are rendered in HTML emails (goldmark CommonMark + GFM; styled with monospace font and light grey background)

## 2026-03-27

- **`gd` Drafts navigation** — jump to Drafts folder with `gd` even when it's not in the tab rotation
- **Off-tab folder indicator** — when viewing Spam (`gS`) or Drafts (`gd`), the folder name appears highlighted in the tab bar with a `│` separator; no regular tab stays falsely active
- **Security hardening** — IMAP refuses unencrypted connections (non-993/143 ports error out instead of `DialInsecure`); email-extracted URLs validated to `http/https` only before opening in browser (case-insensitive, RFC 3986); `SECURITY.md` added documenting credential storage, TLS guarantees, screener list handling, and temp file lifecycle with links to source
- **Spam folder** — `$` marks a sender as spam (writes to `spam.txt`, moves to Spam IMAP folder). Separate from ScreenedOut so you never have to look at it again. Navigate with `gS` or `:go-spam` — kept out of the tab rotation intentionally
- **Cross-list cleanup** — reclassifying a sender removes them from conflicting lists automatically: `I` (approve) removes from screened_out + spam; `O` (block) removes from screened_in; `$` (spam) removes from screened_in + screened_out. No manual `.txt` editing needed
- **`:` command history** — `↑`/`↓` cycles through the last 5 distinct commands; `→` accepts the ghost completion; `ctrl+n`/`ctrl+p` cycle forward/backward through completions. Persists across restarts in `~/.cache/neomd/cmd_history` (outside dotfiles version control)
- **Leader key** — `space` is the leader; `<space>1`–`<space>9` jumps to a folder tab by number
- **Auto-screen on inbox load** — screener applies automatically on every Inbox load (startup, `R`). Disable with `auto_screen_on_load = false` in `[ui]`
- **Background sync** — inbox re-fetched and screened every 5 minutes while neomd is open. Configure with `bg_sync_interval` in `[ui]`; `0` disables it
- **`n` / `m` rebind** — `n` toggles read/unread (was `N`); `m` marks for batch ops (was `space`)

## 2026-03-25

- **Signature** — auto-appended to new compose buffers; configure in `[ui]` with `signature`
- **Compose abort** — closing the editor with `ZQ` / `:q!` cancels the email; only `ZZ` / `:wq` sends
- **Browser image workflow** — `O` opens email as HTML in `$BROWSER`; `ctrl+o` opens the canonical web/newsletter URL (extracted from `List-Post` header); `o` opens in w3m
- **`:create-folders` / `:cf`** — creates any missing IMAP folders defined in config (idempotent)
