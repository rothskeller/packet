package cmd

const (
	deleteSlug = `Delete an unsent message or manual log entry`
	deleteHelp = `
usage: packet delete { «message-id» | «log-entry» }

Given a message ID or log entry of an outgoing message that has not been sent, the "delete" command deletes that message.  For safety, the usual shorthands for message IDs are not accepted by the "delete" command; it must be spelled out fully.  Given a manual log entry, the "delete" command deletes that log entry.  These deletions are not reversible.`
)

func cmdDelete(args []string) (err error) {
	panic("not implemented")
	/*
		flags := pflag.NewFlagSet("delete", pflag.ContinueOnError)
		flags.Usage = func() {} // we do our own
		if err = flags.Parse(args); err == pflag.ErrHelp {
			return cmdHelp([]string{"delete"})
		} else if err != nil {
			cio.Error("%s", err.Error())
			return usage(deleteHelp)
		}
		if len(args) != 1 {
			return usage(deleteHelp)
		}
		args[0] = strings.ToUpper(args[0])
		if !incident.MsgIDRE.MatchString(args[0]) {
			cio.Error(`%q is not a valid, complete message ID`, args[0])
			return usage(deleteHelp)
		}
		env, _, err := incident.ReadMessage(args[0])
		if err != nil {
			return fmt.Errorf("read %s: %s", args[0], err)
		}
		if env.IsFinal() {
			if env.IsReceived() {
				return errors.New("can't delete a received message")
			} else {
				return errors.New("message has already been sent")
			}
		}
		incident.RemoveMessage(args[0])
		cio.Confirm("%s deleted.", args[0])
		return nil
	*/
}
