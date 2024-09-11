package pic18

type SleepController struct {
	OnSleep  func()
	OnWakeUp func()
}

func (sleep *SleepController) Sleep() {
	if sleep != nil && sleep.OnSleep != nil {
		sleep.OnSleep()
	}
}

func (sleep *SleepController) WakeUp() {
	if sleep != nil && sleep.OnWakeUp != nil {
		sleep.OnWakeUp()
	}
}
