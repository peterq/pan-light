package vnc_password

import (
	"github.com/ThomasRooney/gexpect"
	"log"
)

func SetPassword(op, view string) {
	//return
	child, err := gexpect.Spawn("vncpasswd")
	if err != nil {
		panic(err)
	}
	child.Expect("Password")
	child.SendLine(op)
	child.Expect("Verify")
	child.SendLine(op)

	child.Expect("Would you like to enter a view-only password")
	child.SendLine("y")

	child.Expect("Password")
	child.SendLine(view)
	child.Expect("Verify")
	child.SendLine(view)

	child.Interact()
	child.Close()
	log.Println(op, view)
}
