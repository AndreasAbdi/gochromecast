package media

import (
	"github.com/AndreasAbdi/go-castv2/primitives"
)

const seekPositionCommandType = "SEEK"
const seekPosition = "PLAYBACK_START"

//SeekCommand is a type of command payload to be passed along the media channel.
type SeekCommand struct {
	primitives.PayloadHeaders
	CurrentTime float64 `json:"currentTime"`
	ResumeState string  `json:"resumeState"`
}

//CreateSeekCommand creates a seekcommand object.
func CreateSeekCommand(position float64) SeekCommand {
	command := SeekCommand{
		CurrentTime: position,
		ResumeState: seekPosition,
	}
	command.PayloadHeaders = primitives.PayloadHeaders{
		Type: seekPositionCommandType,
	}
	return command
}
