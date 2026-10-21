package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// ANSI Color Codes
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Italic      = "\033[3m"
	Underline   = "\033[4m"
	Black       = "\033[30m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	Gray        = "\033[90m"
	BrightCyan  = "\033[96m"
	BrightGreen = "\033[92m"
	BrightBlue  = "\033[94m"
)

// Spinner displays an animated progress spinner on stderr.
type Spinner struct {
	message string
	stopCh  chan struct{}
	wg      sync.WaitGroup
	active  bool
	mu      sync.Mutex
}

// NewSpinner creates a spinner with an initial message.
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stopCh:  make(chan struct{}),
	}
}

// Start begins animating the spinner.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stopCh:
				// Clear line
				fmt.Fprintf(os.Stderr, "\r\033[K")
				return
			case <-time.After(80 * time.Millisecond):
				frame := frames[i%len(frames)]
				fmt.Fprintf(os.Stderr, "\r%s%s%s %s%s%s", Cyan, frame, Reset, Dim, s.message, Reset)
				i++
			}
		}
	}()
}

// SetMessage updates the spinner text.
func (s *Spinner) SetMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = msg
}

// Stop terminates the spinner animation.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	close(s.stopCh)
	s.mu.Unlock()
	s.wg.Wait()
}

// PrintBanner prints the Nexus logo and version banner.
func PrintBanner(w io.Writer) {
	banner := fmt.Sprintf(`%s%s
    _   __                     ________    ____
   / | / /__  _  ____  _______/ ____/ /   /  _/
  /  |/ / _ \| |/_/ / / / ___/ /   / /    / /  
 / /|  /  __/>  </ /_/ (__  ) /___/ /____/ /   
/_/ |_/\___/_/|_|\__,_/____/\____/_____/___/   
%s%s   ⚡ DeepSeek V3 & R1 Web Bridge CLI for Golang %s
`, BrightCyan, Bold, Reset, Gray, Reset)
	fmt.Fprintln(w, banner)
}
