package sphero

func ExampleAsyncResponse_Sensors() {
	r := &AsyncResponse{}

	// Define a struct containing `int16` fields for each mask value in
	// SetDataStreaming.
	type SensorData struct {
		// Field names don't matter, however the order should match the order
		// defined in the Sphero API spec.
		AccelX, AccelY, AccelZ int16

		// Also, you can ignore fields by giving them the name `_`.
		_ int16
	}

	// Then when receiving an `AsyncResponse`, you would call `Sensors`.
	var data SensorData
	if err := r.Sensors(&data); err != nil {
		// Handle error
	}
}
