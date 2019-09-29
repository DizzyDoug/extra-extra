package types

type Source interface {
	InitSource() error
	CheckForUpdates() error
	GetName() string
}
