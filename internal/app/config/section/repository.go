package section

import "time"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		// TODO: Реализуйте поля для подключения к PostgreSQL
		// Address, Username, Password, Name — string, тег required:"true"
		Address  string `required:"true"`
		Username string `required:"true"`
		Password string `required:"true"`
		Name     string `required:"true"`
		// ReadTimeout, WriteTimeout — time.Duration, теги:
		//   split_words:"true" default:"30s"
		ReadTimeout  time.Duration `split_words:"true" default:"30s"`
		WriteTimeout time.Duration `split_words:"true" default:"30s"`
		// Важно: дефолт таймаутов именно "30s" (не "5s") — как в .env / .env.example ниже
	}
)
