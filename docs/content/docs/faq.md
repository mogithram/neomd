---
title: FAQ
weight: 50
---

Questions that came up when people using neomd.


## Is it possible to create new directories/tabs

You basically create the folder in your web mail and configure it in your `config.toml` and add the new folder under `[folder]` and in the `tab_order` so neomd knows where to place it:


```toml
[folders]
  ...existing folders
  new = "NewMissingFolder"
  tab_order = ["inbox", "to_screen", "feed", "papertrail", "waiting", "someday", "scheduled", "sent", "work", "archive", "screened_out", "trash", "new"]
```

If you want to move emails to that folder, or just move to it, that's currently not possible. You can always move through the tabs with `[]HL` or `space+1-10`, but you can't move emails to them yet.

## Does the signature appear only in new messages, not in replies?

Currently the signature is only automatically added if you create and compose a **new email**.

But you can add the signature in any email, e.g. if you reply with `[html-signature]` like this:

```markdown

# [neomd: to: email@domain.com]
# [neomd: from: Simon <my-email@ssp.sh>]
# [neomd: subject: Re: Subject Title]


Here's my reply.
BR Simon

[html-signature]

---

> **Previous sender <email@domain.com>** wrote:
>
> * * *
......

```

The html-signature is the placeholder for adding the HTML signature, but yes, it will always be added at the end of the email (e.g. in this case the reply).
