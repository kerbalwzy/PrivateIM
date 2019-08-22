package DataLayer

import "log"

// todo record the message which send failed into database
func RecordWaitSendMessage(id int64, data []byte) {
	log.Printf("save a WaitSendMessage for user<%d>, message: %s\n", id, data)
}
