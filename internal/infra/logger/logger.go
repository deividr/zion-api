package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Level string

const (
	INFO  Level = "INFO"
	ERROR Level = "ERROR"
	WARN  Level = "WARN"
	DEBUG Level = "DEBUG"
)

type Logger struct {
	log zerolog.Logger
}

func New() *Logger {
	// Configura o zerolog para usar formato mais legível de tempo
	zerolog.TimeFieldFormat = time.RFC3339

	// Cria um logger com saída formatada e colorida para console
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	// Configura o nível mínimo de log e adiciona timestamp por padrão
	logger := zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{
		log: logger,
	}
}

func (l *Logger) Error(message string, err error) {
	l.log.Error().
		Err(err).
		Msg(message)
}

func (l *Logger) Info(message string) {
	l.log.Info().
		Msg(message)
}

func (l *Logger) Warn(message string) {
	l.log.Warn().
		Msg(message)
}

func (l *Logger) Debug(message string) {
	l.log.Debug().
		Msg(message)
}

// Métodos auxiliares para adicionar contexto aos logs
func (l *Logger) WithFields(fields map[string]interface{}) *zerolog.Event {
	event := l.log.Info()
	for k, v := range fields {
		event.Interface(k, v)
	}
	return event
}

// Exemplo de uso com contexto estruturado
func (l *Logger) InfoWithContext(message string, fields map[string]interface{}) {
	event := l.log.Info()
	for k, v := range fields {
		event.Interface(k, v)
	}
	event.Msg(message)
}
