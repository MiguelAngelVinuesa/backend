package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	assert.Equal(t, "INFO", InfoLevel.String())
	assert.Equal(t, "DEBUG", DebugLevel.String())
	assert.Equal(t, "WARN", WarnLevel.String())
	assert.Equal(t, "ERROR", ErrorLevel.String())
	assert.Equal(t, "PANIC", PanicLevel.String())
	assert.Equal(t, "INFO", Level(0).String())
	assert.Equal(t, "INFO", Level(99).String())
}

func TestLevelFromName(t *testing.T) {
	assert.Equal(t, InfoLevel, LevelFromName("info"))
	assert.Equal(t, DebugLevel, LevelFromName("deBug"))
	assert.Equal(t, WarnLevel, LevelFromName("warning"))
	assert.Equal(t, WarnLevel, LevelFromName("warn"))
	assert.Equal(t, ErrorLevel, LevelFromName("Err"))
	assert.Equal(t, ErrorLevel, LevelFromName("eRRor"))
	assert.Equal(t, PanicLevel, LevelFromName("panic"))
	assert.Equal(t, InfoLevel, LevelFromName(""))
	assert.Equal(t, InfoLevel, LevelFromName("xyz"))
}
