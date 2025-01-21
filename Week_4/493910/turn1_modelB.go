package main  
import (  
    "fmt"
    "log"
    "sync"
    "time"
)

// FunnelStage represents each stage of the user engagement funnel
type FunnelStage struct {
    Name      string
    Handle    func(user *User, wg *sync.WaitGroup)
    NextStage *FunnelStage
}

// User represents a user interacting with the e-commerce website
type User struct {
    ID          int
    Name        string
    Engaged     bool
    Conversion  bool
    ClickThrough bool
}

func (fs *FunnelStage) Run(user *User, wg *sync.WaitGroup) {
    defer wg.Done()
    fs.Handle(user, wg)
    if fs.NextStage != nil {
        fs.NextStage.Run(user, wg)
    }
}

// SimulateUserInteraction simulates a user interacting with the website.
func SimulateUserInteraction(user *User, wg *sync.WaitGroup) {
    time.Sleep(time.Duration(2 * time.Second)) // Simulate interaction time
    user.Engaged = true
    log.Printf("User %d is engaged in the funnel", user.ID)
}

// CalculateClickThroughRate calculates the click-through rate based on user engagement.
func CalculateClickThroughRate(user *User, wg *sync.WaitGroup) {
    time.Sleep(time.Duration(1 * time.Second)) // Simulate calculation time
    if user.Engaged {
        user.ClickThrough = true
    }
    log.Printf("User %d has a click-through rate of %t", user.ID, user.ClickThrough)
}

// TrackConversion tracks user conversion based on click-through rate and other factors.
func TrackConversion(user *User, wg *sync.WaitGroup) {
    time.Sleep(time.Duration(3 * time.Second)) // Simulate conversion tracking time
    if user.ClickThrough {
        user.Conversion = true
    }
    log.Printf("User %d has converted: %t", user.ID, user.Conversion)
}

func main() {  
    numUsers := 5
    var wg sync.WaitGroup
    funnelStart := &FunnelStage{
        Name:      "Start",
        Handle:    SimulateUserInteraction,
    }
    
    funnelStart.NextStage = &FunnelStage{
        Name:      "Click-Through Rate Calculation",
        Handle:    CalculateClickThroughRate,
    }
    
    funnelStart.NextStage.NextStage = &FunnelStage{
        Name:      "Conversion Tracking",
        Handle:    TrackConversion,
    }
    
    for i := 0; i < numUsers; i++ {
        user := &User{ID: i + 1, Name: fmt.Sprintf("User %d", i+1)}
        wg.Add(1)
        go funnelStart.Run(user, &wg)
    }
    wg.Wait()
    log.Println("All users' funnel metrics tracked successfully.")
}
