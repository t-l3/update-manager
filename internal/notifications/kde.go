package notifications

import (
	"github.com/godbus/dbus"
)

type notification struct {
	conn *dbus.Conn
	path *string
}

func New(msg string, icon string) notification {
	conn, _ := dbus.SessionBus()
	obj := conn.Object("org.freedesktop.Notifications", "/JobViewServer")
	call := obj.Call("requestView", 0, "update-manager", icon, 0)

	var response string
	call.Store(&response)

	notif := notification{
		path: &response,
		conn: conn,
	}

	notif.SetInfoMessage(msg)

	return notif
}

func (n *notification) SetPercent(percent int) {
	<-n.conn.Object("org.freedesktop.Notifications", dbus.ObjectPath(*n.path)).Call("org.kde.JobViewV2.setPercent", 1, uint32(percent)).Done
}

func (n *notification) SetInfoMessage(msg string) {
	<-n.conn.Object("org.freedesktop.Notifications", dbus.ObjectPath(*n.path)).Call("org.kde.JobViewV2.setInfoMessage", 1, msg).Done
}

func (n *notification) Terminate(msg string) {
	<-n.conn.Object("org.freedesktop.Notifications", dbus.ObjectPath(*n.path)).Call("org.kde.JobViewV2.terminate", 1, msg).Done
}
