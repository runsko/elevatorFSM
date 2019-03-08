package main

import "./driver/elevio"
import "fmt"

/*
func elevatorServer(chan newOrder int, chan FloorSensor int, chan DoorControl bool, chan ButtonInput elevio.ButtonEvent) {

}
*/

func main() {
	numberOfFloors := 4

	elevio.Init("localhost:15657", numberOfFloors)

	driver_btn := make(chan elevio.ButtonEvent)
	driver_floors := make(chan int)

	go elevio.PollButtons(driver_btn)
	go elevio.PollFloorSensor(driver_floors)

	for {
		select {
		case input := <-driver_btn:
			switch input.Button {
			case elevio.BT_HallUp:
				fmt.Printf("Hall Call Up at Floor %d\n", input.Floor)
				elevio.SetMotorDirection(elevio.MD_Up)
			case elevio.BT_HallDown:
				fmt.Printf("Hall Call Down at Floor %d\n", input.Floor)
				elevio.SetMotorDirection(elevio.MD_Down)
			case elevio.BT_Cab:
				fmt.Printf("Cab call to Floor %d\n", input.Floor)
			}
			elevio.SetButtonLamp(input.Button, input.Floor, true)
		}
	}

}
