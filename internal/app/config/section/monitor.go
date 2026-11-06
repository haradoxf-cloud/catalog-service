package section

type Monitor struct {
	// TODO: Добавьте поля LogLevel и Environment
	// LogLevel по умолчанию "debug"
	// Environment по умолчанию "development"
	LogLevel    string `split_words:"true" default:"debug"`
	Environment string `default:"development"`
}
