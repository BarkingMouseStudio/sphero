package main

import (
	"fmt"
	"github.com/Freeflow/sphero"
	"time"
)

func main() {
	fmt.Println("Connecting...")
	mask1 := sphero.ApplyMasks32([]uint32{
		sphero.ACCEL_AXIS_X_RAW, sphero.ACCEL_AXIS_Y_RAW, sphero.ACCEL_AXIS_Z_RAW,
		sphero.GYRO_AXIS_X_RAW, sphero.GYRO_AXIS_Y_RAW, sphero.GYRO_AXIS_Z_RAW,
		sphero.MOTOR_RIGHT_EMF_RAW, sphero.MOTOR_LEFT_EMF_RAW,
		sphero.MOTOR_LEFT_PWM_RAW, sphero.MOTOR_RIGHT_PWM_RAW,
		sphero.IMU_PITCH_ANGLE_FILTERED, sphero.IMU_ROLL_ANGLE_FILTERED, sphero.IMU_YAW_ANGLE_FILTERED,
	})

	mask2 := sphero.ApplyMasks32([]uint32{
		sphero.QUATERNION_Q0, sphero.QUATERNION_Q1, sphero.QUATERNION_Q2, sphero.QUATERNION_Q3,
		sphero.ODOMETER_X, sphero.ODOMETER_Y,
		sphero.VELOCITY_X, sphero.VELOCITY_Y,
	})
	fmt.Printf("%#x %#x\n", mask1, mask2)
	return

	async := make(chan *sphero.AsyncResponse, 256)
	go func() {
		for r := range async {
			fmt.Println(r)
		}
	}()

	var err error
	var s *sphero.Sphero
	if s, err = sphero.NewSphero("/dev/cu.Sphero-YBR-RN-SPP", async); err != nil {
		fmt.Println(err)
		fmt.Println("NOTE: Try toggling your Bluetooth or putting the Sphero asleep (and reawakening) to fix 'resource busy' errors")
		return
	}

	fmt.Println("Connected.")

	var res *sphero.Response
	ch := make(chan *sphero.Response, 1)

	// PING

	fmt.Println("Ping...")
	s.Ping(ch)
	res = <-ch
	fmt.Printf("Pong %#x\n", res)

	// COLOR

	fmt.Println("Setting color...")
	s.SetRGBLEDOutput(0, 0, 255, ch)
	res = <-ch
	fmt.Printf("Set Color %#x\n", res)

	fmt.Println("Getting color...")
	s.GetRGBLED(ch)
	res = <-ch
	fmt.Printf("Get Color %#x\n", res)

	// ASYNC

	fmt.Println("Turning on async...")
	// 1000ms / (400hz / 40hz) = 1 @ 100ms
	// 1000ms / (400hz / 4hz) = 1 @ 10ms
	s.SetDataStreaming(40, 0, mask1, 0, mask2, ch)
	res = <-ch
	fmt.Printf("Data streaming %#x\n", res)

	<-time.Tick(10 * time.Second)

	// CLEANUP

	fmt.Println("Sleeping...")
	s.Sleep(time.Duration(0), 0, 0, ch)
	res = <-ch
	fmt.Printf("Slept %#x\n", res)

	fmt.Println("Closing...")
	s.Close()
	fmt.Println("Closed.")
}
