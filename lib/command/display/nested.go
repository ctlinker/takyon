package display

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	width  int
	height int
	name   string
)

var nestedDisplayCmd = &cobra.Command{
	Use:   "nested",
	Short: "Start a nested X display using Xephyr",
	Run: func(cmd *cobra.Command, args []string) {
		display := pickDisplay()

		fmt.Printf("[→] Starting Xephyr on :%d (%dx%d)\n", display, width, height)

		xephyr := exec.Command(
			"Xephyr",
			fmt.Sprintf(":%d", display),
			"-screen", fmt.Sprintf("%dx%d", width, height),
			"-ac", // disable access control (optional)
			"-title",
			fmt.Sprintf("Xephyr-[%s]", name),
		)

		xephyr.Stdout = os.Stdout
		xephyr.Stderr = os.Stderr

		if err := xephyr.Start(); err != nil {
			fmt.Println("[!] Failed to start Xephyr:", err)
			return
		}

		// Give Xephyr time to initialize
		time.Sleep(500 * time.Millisecond)

		fmt.Printf("[→] Launching session inside DISPLAY=:%d\n", display)

		session := exec.Command("bash")
		session.Env = append(os.Environ(),
			"DISPLAY=:"+strconv.Itoa(display),
		)

		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		if err := session.Run(); err != nil {
			fmt.Println("[!] Session exited with error:", err)
		}

		fmt.Println("[→] Shutting down Xephyr...")
		xephyr.Process.Kill()
	},
}

func pickDisplay() int {
	for i := 100; i < 200; i++ {
		if _, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", i)); os.IsNotExist(err) {
			return i
		}
	}
	return 100
}

func init() {
	nestedDisplayCmd.Flags().IntVar(&width, "width", 1280, "Display width")
	nestedDisplayCmd.Flags().IntVar(&height, "height", 720, "Display height")
	nestedDisplayCmd.Flags().StringVar(&name, "name", "default", "Display name")

	DisplayCMD.AddCommand(nestedDisplayCmd)
}
