package axvdo

// GetSnapshot captures a JPEG snapshot from the specified video channel and returns it as a byte slice.
// It sets up the required settings for capturing the snapshot, captures it, and then returns the snapshot data or an error if the capture fails.
func GetSnapshot(video_channel int) ([]byte, error) {
	settings := NewVdoMap()                              // Create a new settings map for the snapshot.
	settings.SetUint32("channel", uint32(video_channel)) // Set the video channel.
	settings.SetUint32("format", uint32(VdoFormatJPEG))  // Set the snapshot format to JPEG.
	defer settings.Unref()                               // Ensure settings are unreferenced after use.

	snapshotBuffer, err := Snapshot(settings) // Capture the snapshot.
	if err != nil {
		return nil, err
	}
	defer snapshotBuffer.Unref() // Ensure the snapshot buffer is unreferenced after use.

	return snapshotBuffer.GetBytes() // Return the snapshot data.
}
