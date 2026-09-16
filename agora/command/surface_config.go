package command

type SurfaceSettings struct {
	Bind      string
	ReachedAt string
}

type SurfaceLoader func() (SurfaceSettings, error)
