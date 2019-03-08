package main

import (
	"fmt"

	//. "./elevFSMFunctions"

	. "./driver/elevio"
)

type state int

const (
	Idle     state = 1
	Moving         = 2
	DoorOpen       = 3
)

type elev struct {
	State     state
	Direction int
	Floor     int
	Queue     [4][3]bool //array of floors and buttontypes
}

func main() {

	numberOfFloors := 4
	//numberOfButtontypes := 3

	elevator := elev{
		State:     Idle,
		Direction: 0,
		Floor:     1,
		Queue:     [4][3]bool{},
	}

	Init("localhost:15657", numberOfFloors)

	driver_btn := make(chan ButtonEvent)
	driver_floors := make(chan int)

	go PollButtons(driver_btn)
	go PollFloorSensor(driver_floors)

	for {
		select {
		case input := <-driver_btn:

			SetButtonLamp(input.Button, input.Floor, true)
			fmt.Printf("%d -type call at floor %d\n", input.Button, input.Floor)

			switch elevator.State {
			case Idle:
				if elevator.Floor == input.Floor {
					SetDoorOpenLamp(true)
					//timer either here or in a function
					elevator.State = DoorOpen
				} else {
					elevator.Queue[input.Floor][input.Button] = true
					fmt.Printf("Queue is updated: %t", elevator.Queue[input.Floor][input.Button])

				}

			case Moving:
				elevator.Queue[input.Floor][input.Button] = true
				fmt.Printf("Queue is updated: %t", elevator.Queue[input.Floor][input.Button])
			case DoorOpen:
				if elevator.Floor == input.Floor {
					// reset timer
				} else {
					elevator.Queue[input.Floor][input.Button] = true
				}

			}

		case sensor := <-driver_floors:
			switch elevator.State {
			case Idle: //Do nothing

			case Moving:
				if elevator.Queue[sensor][BT_HallUp] == true && elevator.Direction == 1 {
					SetDoorOpenLamp(true)
					//timer
					elevator.State = DoorOpen
					elevator.Queue[sensor][BT_HallUp] = false //Clear queue
					elevator.Queue[sensor][BT_Cab] = false
				} else if elevator.Queue[sensor][BT_HallDown] == true && elevator.Direction == 1 {
					SetDoorOpenLamp(true)
					//timer
					elevator.State = DoorOpen
					elevator.Queue[sensor][BT_HallDown] = false //Clear queue
					elevator.Queue[sensor][BT_Cab] = false
				} else {
					continue
				}

			case DoorOpen: //Do nothing

			}

		}

	}
}
