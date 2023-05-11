package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Haydz6/rich-go/ipc"
)

type AuthenticatedStruct struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type ReceivedPayloadStruct struct {
	Evt  string `json:"evt"`
	Data struct {
		User AuthenticatedStruct `json:"user"`
	} `json:"data"`
}

var logged bool
var CachedClientId string
var Authentication *AuthenticatedStruct
var AuthenticationUpdate = sync.WaitGroup{}
var LogLooping bool

// Login sends a handshake in the socket and returns an error or nil

func LoginLoop() {
	if LogLooping {
		return
	}

	LogLooping = true
	for {
		time.Sleep(time.Second * 5)

		if !logged && CachedClientId != "" {
			Login(CachedClientId)
		}
	}
}

func CheckForClosure(Result string) bool {
	if Result == "The pipe is being closed." {
		logged = false
		Authentication = nil
		AuthenticationUpdate.Done()

		ipc.CloseSocket()
		return true
	}
	return false
}

func Login(clientid string) error {
	go LoginLoop()
	CachedClientId = clientid
	if !logged {
		payload, err := json.Marshal(Handshake{"1", clientid})
		if err != nil {
			return err
		}

		err = ipc.OpenSocket()
		if err != nil {
			return err
		}

		// TODO: Response should be parsed
		Result := ipc.Send(0, string(payload))
		if !CheckForClosure(Result) {
			var Body ReceivedPayloadStruct
			JSONErr := json.Unmarshal([]byte(Result), &Body)
			if JSONErr == nil {
				if Body.Evt == "READY" {
					Authentication = &Body.Data.User
					AuthenticationUpdate.Done()
				}
			}
		}
	}

	logged = true

	return nil
}

func Logout() {
	logged = false
	Authentication = nil
	AuthenticationUpdate.Done()

	err := ipc.CloseSocket()
	if err != nil {
		panic(err)
	}
}

func SetActivity(activity Activity) error {
	if !logged {
		if CachedClientId == "" {
			return nil
		}
		Login(CachedClientId)
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
	CheckForClosure(ipc.Send(1, string(payload)))
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
