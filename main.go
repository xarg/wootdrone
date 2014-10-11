package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
	"github.com/hybridgroup/gobot/platforms/leap"
)

func main() {
	gbot := gobot.NewGobot()

	droneAdaptor := ardrone.NewArdroneAdaptor("Drone")
	d := ardrone.NewArdroneDriver(droneAdaptor, "Drone")

	leapMotionAdaptor := leap.NewLeapMotionAdaptor("leap", "127.0.0.1:6437")
	l := leap.NewLeapMotionDriver(leapMotionAdaptor, "leap")

	var frame leap.Frame
	var altitude float64
	altitude = 5
	workLeap := func() {
		gobot.On(l.Event("message"), func(data interface{}) {
			frame = data.(leap.Frame)
			for _, hand := range frame.Hands {
				altitude = hand.Y() / 10
				fmt.Println(altitude)
			}
		})
	}

	workDrone := func() {
		d.TakeOff()
		gobot.On(d.Events()["Flying"], func(data interface{}) {
			if altitude < 4 {
				d.Land()
			}
			if altitude > 10 {
				d.Up(0.1)
			}

			if altitude < 10 {
				d.Down(0.1)
			}
		})
	}

	robotLeap := gobot.NewRobot("leapBot",
		[]gobot.Connection{leapMotionAdaptor},
		[]gobot.Device{l},
		workLeap,
	)
	robotDrone := gobot.NewRobot("drone",
		[]gobot.Connection{droneAdaptor},
		[]gobot.Device{d},
		workDrone,
	)

	gbot.AddRobot(robotLeap)
	gbot.AddRobot(robotDrone)

	gbot.Start()
}
