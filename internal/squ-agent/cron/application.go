package cron

func (c *Cron) startApp() {
	c.Cron.AddFunc("*/30 * * * *", func() {
		return
	})
}
