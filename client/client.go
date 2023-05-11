package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Haydz6/rich-go/ipc"
)

type AuthenticatedStruct struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type ReceivedPayloadStruct struct {
	Evt  string              `json:"evt"`
	Data AuthenticatedStruct `json:"data"`
}

var logged bool
var Authentication *AuthenticatedStruct
var AuthenticationUpdate = sync.Cond{}

// Login sends a handshake in the socket and returns an error or nil
func Login(clientid string) error {
	if !logged {
		payload, err := json.Marshal(Handshake{"1", clientid})
		if err != nil {
			return err
		}

		err = ipc.OpenSocket()
		if err != nil {
			return err
		}

		go func() {
			for {
				Data := ipc.Read()

				if Data == "Connection Closed" {
					logged = false
					Authentication = nil
					AuthenticationUpdate.Broadcast()

					ipc.CloseSocket()
					break
				}

				var Instruction ReceivedPayloadStruct
				err := json.Unmarshal([]byte(Data), &Instruction)

				if err != nil {
					continue
				}

				if Instruction.Evt == "READY" {
					Authentication = &Instruction.Data
					AuthenticationUpdate.Broadcast()
				}
			}
		}()

		// TODO: Response should be parsed
		ipc.Send(0, string(payload))
	}
	logged = true

	return nil
}

func Logout() {
	logged = false
	Authentication = nil
	AuthenticationUpdate.Broadcast()

	err := ipc.CloseSocket()
	if err != nil {
		panic(err)
	}
}

func SetActivity(activity Activity) error {
	if !logged {
		return nil
	}

	var Arguments Args

	if activity.State != "end" {
		Arguments = Args{
			os.Getpid(),
			mapActivity(&activity),
		}
	} else {
		Arguments = Args{
			os.Getpid(),
			nil,
		}
	}

	payload, err := json.Marshal(Frame{
		"SET_ACTIVITY",
		Arguments,
		getNonce(),
	})

	if err != nil {
		return err
	}

	// TODO: Response should be parsed
	ipc.Send(1, string(payload))
	return nil
}

func getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
