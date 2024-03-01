package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "main/voxieproto"
	"log"
	"io"
	"io/ioutil"
	"os"
	"time"
	"context"
	"sync/atomic"
	"encoding/json"
)

type Tile struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

var (
	Placement PlacementInfo
	startTile Tile
 	mapName string
 	teamId int32
 	randomNumberClue int32
 	id atomic.Uint32
 	charId int32
 	timeBattleStarted time.Time
)

type PlacementInfo []struct {
	MapName        string `json:"mapName"`
	PlacementTiles []Tile `json:"placementTiles"`
}

func init() {
    jsonFile, err := os.Open("placement.json")
    if err != nil {
        log.Fatalln(err)
    }

    log.Println("Successfully Opened placement.json")
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal(byteValue, &Placement)
}

func EasyProcessRequest(client pb.PlayingClient, request *pb.Request) {
    response, err := client.ProcessRequest(context.Background(), request)
    if err != nil {
        grpclog.Fatalf("failed to dial: %v", err)
    }

    log.Println("SEND:", request)
	log.Println("RECV:", response)
	id.Add(1)
}


func (usr *Instance) SkipTurn(client pb.PlayingClient) {
	endTurn := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_Make_Turn{
			Make_Turn: &pb.MakeTurnRequest{
				CharacterIndex: charId,
				Data: &pb.MakeTurnRequest_EndTurn{
					EndTurn: &pb.EndTurnRequest{
						ReduceTimer: true,
					},
				},
			},
		},
	}

	EasyProcessRequest(client, endTurn)
}

func (usr *Instance) KillOpponent(client pb.PlayingClient) {
	var oppTile Tile
	var oppId int32

	if teamId == 1 {
		oppId = 0
	} else {
		oppId = 1
	}

    for _, e := range Placement {
    	if e.MapName == mapName {
    		oppTile = e.PlacementTiles[oppId]
    		break
    	}
    }

	attackRequest := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_Make_Turn{
			Make_Turn: &pb.MakeTurnRequest{
				CharacterIndex: charId,
				Data: &pb.MakeTurnRequest_Action{
					Action: &pb.ActionRequest{
						Data: &pb.ActionRequest_Attack{
							Attack: &pb.AttackRequest{
								Target_Position: &pb.Position{X: oppTile.X, Y: oppTile.Y},
								DamageDone: 500,
								HitChance: 100,
								HitOverride: true,
								SpellProcOverride: true,
							},
						},
					},
				},
			},
		},
	}

	EasyProcessRequest(client, attackRequest)
}

func (usr *Instance) SendWinRequest(client pb.PlayingClient) {
    var randomNumberAnswer int32 = 0

    switch mapName {
    case "map1":
    	randomNumberAnswer = (randomNumberClue - 8) * 8 + 190
    case "map2":
    	randomNumberAnswer = (randomNumberClue - 7) * 3 + 679
    case "map3":
    	randomNumberAnswer = (randomNumberClue - 5) * 7 - 222
    case "map4":
    	randomNumberAnswer = (randomNumberClue - 1) * 9 + 587
    case "map5":
    	randomNumberAnswer = (randomNumberClue - 3) * 3 + 51
    case "map6":
    	randomNumberAnswer = (randomNumberClue - 2) * 7 + 117
    case "map7":
    	randomNumberAnswer = (randomNumberClue - 6) * 8 + 338
    }

    answerStart := []string{"tess", "zatter", "ios", "plead", "f4tt", "opp-osi", "oxion"}
    answerEnd := []string{"yawn", "rex", "szuickk", "wwwrww", "trye", "fal-fal", "ingly", "stop", "reampoi"}

    var s1, s2 string

    s1 = answerStart[randomNumberAnswer % 7]
    s2 = answerEnd[randomNumberAnswer % 9]

	winRequest := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_Win{
			Win: &pb.WinRequest{
				Team_Id: teamId,
				RandomNumberAnswer: randomNumberAnswer,
				RandomStringAnswer: s1+s2,
			},
		},
	}

	EasyProcessRequest(client, winRequest)
}

func (usr *Instance) StartGame(client pb.PlayingClient) {
    time.Sleep(3 * time.Second)

    for _, e := range Placement {
    	if e.MapName == mapName {
    		startTile = e.PlacementTiles[teamId]
    		break
    	}
    }

	placeChar := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_Place_Characters{
			Place_Characters: &pb.PlaceCharacterRequest{
				CharacterIndex: charId,
				Position: &pb.Position{X: startTile.X, Y: startTile.Y},
			},
		},
	}

	EasyProcessRequest(client, placeChar)

	readyRequest := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_Ready_To_Play{
			Ready_To_Play: &pb.ReadyToPlayRequest{Ready: true},
		},
	}

	EasyProcessRequest(client, readyRequest)
}

func (usr *Instance) Play(addr, attribute string, winningMode bool) {
	id.Store(0)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUserAgent("grpc-c++/1.41.0 grpc-c/19.0.0 (windows; chttp2)"),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPlayingClient(conn)

	id.Add(1)

	go func() {
		for {
			time.Sleep(4 * time.Second)

			pingRequest := &pb.Request{
				Request_Id: id.Load(),
				User_Id: usr.username,
				Data: &pb.Request_Client_Ping{},
			}

		    response, err := client.ProcessRequest(context.Background(), pingRequest)
		    if err != nil {
		        log.Println("failed to ping: %v", err)
		        break
		    }

		    log.Println("SEND:", pingRequest)
			log.Println("RECV:", response)
			id.Add(1)
		}
	}()

	joinRequest := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_Join_Game{
			Join_Game: &pb.JoinGameRequest{
				Client_Access_Token: usr.token,
				Custom_Attribute: attribute,
			},
		},
	}

	EasyProcessRequest(client, joinRequest)

	notifRequest := &pb.NotificationSubscribeRequest{
		User_Id: usr.username,
	}

	waitc := make(chan struct{})

	go func() {
	    stream, err := client.ProcessNotification(context.Background(), notifRequest)
	    if err != nil {
	        grpclog.Fatalf("failed to dial: %v", err)
	    }

	    log.Println("SEND:", notifRequest)

		for {
		    notif, err := stream.Recv()
		    if err == io.EOF {
		    	close(waitc)
		        return
		    }

		    if initNotif := notif.GetInit(); initNotif != nil {
		    	mapName = initNotif.MapName
		    	randomNumberClue = initNotif.RandomNumberClue
		    	teamId = initNotif.Team_Id
		    }

		    if getMpTeamNotif := notif.GetMultiPlayer_Team(); getMpTeamNotif != nil {
		    	if getMpTeamNotif.User_Id == usr.username{
		    		charId = getMpTeamNotif.Characters[2].VoxieIndex
		    	}
		    }

		    if startNotif := notif.GetStart_Game(); startNotif != nil {
		    	usr.StartGame(client)
		    }

		    if battleNotif := notif.GetBattle_Started(); battleNotif != nil {
		    	timeBattleStarted = time.Now()
		    	if teamId == 1 {
		    		time.Sleep(1 * time.Second)
		    		usr.SkipTurn(client)
		    	}
		    } 

		    if endTurn := notif.GetEndTurn_Character(); endTurn != nil {
		    	if endTurn.User_Id != usr.username {
		    		time.Sleep(10*time.Second)
		    		if winningMode && (time.Since(timeBattleStarted).Seconds() > 270) {
		    			usr.KillOpponent(client)
		    		}
		    		usr.SkipTurn(client)	
		    	}
		    }

			if attackNotif := notif.GetAttack_Character(); attackNotif != nil {
				usr.SendWinRequest(client)
			}

		    if err != nil {
		        log.Fatalf("%v.ProcessNotification(_) = _, %v", client, err)
		    }
		    log.Println("RECV:", notif)
		}
	}()

	for mapName == "" {

	}

	loadingComplete := &pb.Request{
		Request_Id: id.Load(),
		User_Id: usr.username,
		Data: &pb.Request_LoadingComplete{
			LoadingComplete: &pb.LoadingCompleteRequest{},
		},
	}

	EasyProcessRequest(client, loadingComplete)

    <-waitc
}