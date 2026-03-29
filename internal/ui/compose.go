package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// composeStep tracks which field is active in the compose form.
type composeStep int

const (
	stepTo      composeStep = iota
	stepCC                  // only reachable when extraVisible=true
	stepBCC                 // only reachable when extraVisible=true
	stepSubject             // after Subject is filled, launch editor
)

// composeModel holds state for the compose view.
type composeModel struct {
	to           textinput.Model
	cc           textinput.Model
	bcc          textinput.Model
	subject      textinput.Model
	step         composeStep
	extraVisible bool // ctrl+b toggles Cc+Bcc together; off by default
}

func newComposeModel() composeModel {
	to := textinput.New()
	to.Placeholder = "recipient@example.com"
	to.Focus()
	to.CharLimit = 256
	to.Width = 60
	to.Prompt = ""

	cc := textinput.New()
	cc.Placeholder = "cc@example.com (optional)"
	cc.CharLimit = 512
	cc.Width = 60
	cc.Prompt = ""

	bcc := textinput.New()
	bcc.Placeholder = "bcc@example.com (optional)"
	bcc.CharLimit = 512
	bcc.Width = 60
	bcc.Prompt = ""

	sub := textinput.New()
	sub.Placeholder = "Subject"
	sub.CharLimit = 256
	sub.Width = 60
	sub.Prompt = ""

	return composeModel{to: to, cc: cc, bcc: bcc, subject: sub, step: stepTo}
}

// reset clears all fields and refocuses on To.
func (c *composeModel) reset() {
	c.to.Reset()
	c.cc.Reset()
	c.bcc.Reset()
	c.subject.Reset()
	c.step = stepTo
	c.extraVisible = false
	c.to.Focus()
	c.cc.Blur()
	c.bcc.Blur()
	c.subject.Blur()
}

// update handles key input for the compose form.
// Returns (updated model, cmd, launchEditor bool).
func (c composeModel) update(msg tea.Msg) (composeModel, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+b":
			c.extraVisible = !c.extraVisible
			if !c.extraVisible {
				c.cc.Reset()
				c.bcc.Reset()
				// If cursor was on CC or BCC, jump to Subject
				if c.step == stepCC || c.step == stepBCC {
					c.cc.Blur()
					c.bcc.Blur()
					c.step = stepSubject
					c.subject.Focus()
				}
			}
			return c, nil, false

		case "tab", "enter":
			switch c.step {
			case stepTo:
				if c.extraVisible {
					c.step = stepCC
					c.to.Blur()
					c.cc.Focus()
				} else {
					c.step = stepSubject
					c.to.Blur()
					c.subject.Focus()
				}
				return c, nil, false
			case stepCC:
				c.step = stepBCC
				c.cc.Blur()
				c.bcc.Focus()
				return c, nil, false
			case stepBCC:
				c.step = stepSubject
				c.bcc.Blur()
				c.subject.Focus()
				return c, nil, false
			case stepSubject:
				return c, nil, true
			}
		}
	}

	var cmd tea.Cmd
	switch c.step {
	case stepTo:
		c.to, cmd = c.to.Update(msg)
	case stepCC:
		c.cc, cmd = c.cc.Update(msg)
	case stepBCC:
		c.bcc, cmd = c.bcc.Update(msg)
	default:
		c.subject, cmd = c.subject.Update(msg)
	}
	return c, cmd, false
}

// view renders the compose header form.
func (c composeModel) view() string {
	toLabel := styleInputLabel.Render("To:")
	subLabel := styleInputLabel.Render("Subject:")

	toField := c.to.View()
	subField := c.subject.View()

	switch c.step {
	case stepTo:
		toField = styleInputField.Render(toField)
	default:
		subField = styleInputField.Render(subField)
	}

	out := toLabel + toField + "\n"

	if c.extraVisible {
		ccLabel := styleInputLabel.Render("Cc:")
		bccLabel := styleInputLabel.Render("Bcc:")
		ccField := c.cc.View()
		bccField := c.bcc.View()
		if c.step == stepCC {
			ccField = styleInputField.Render(ccField)
			subField = c.subject.View() // undo active styling on Subject
		}
		if c.step == stepBCC {
			bccField = styleInputField.Render(bccField)
			subField = c.subject.View()
		}
		out += ccLabel + ccField + "\n"
		out += bccLabel + bccField + "\n"
	} else {
		out += styleHelp.Render("  ctrl+b to add Cc/Bcc") + "\n"
	}

	out += subLabel + subField
	return out
}
