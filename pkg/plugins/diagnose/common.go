package diagnose

import "fmt"

func CommonDeafer(c chan *Result) {
	close(c)
	if err := recover(); err != nil {
		c <- &Result{
			Error: fmt.Errorf("%v", err),
		}
	}
}
