package knockback

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/restartfu/gophig"
	"os"
	"path/filepath"
	"time"
)

var (
	// force is the force of the knockback.
	force, height = 0.4, 0.4
	// hitDelay is the delay between hits.
	hitDelay = 500 * time.Millisecond
)

// settings is a struct that holds the settings for the knockback library.
type settings struct {
	Force    float64
	Height   float64
	HitDelay int64
}

// Load loads the settings from the file at the path passed.
func Load(path string) error {
	s := settings{
		Force:    force,
		Height:   height,
		HitDelay: hitDelay.Milliseconds(),
	}

	_ = os.MkdirAll(filepath.Dir(path), 0777)
	err := gophig.GetConfComplex(path, gophig.JSONMarshaler{}, &s)
	if err != nil {
		if os.IsNotExist(err) {
			save()
			return nil
		}
		return err
	}

	force, height, hitDelay = s.Force, s.Height, time.Duration(s.HitDelay)*time.Millisecond
	return nil
}

// ApplyForce applies the force to the knockback.
func ApplyForce(f *float64) {
	*f = force
}

// ApplyHeight applies the height to the knockback.
func ApplyHeight(h *float64) {
	*h = height
}

// ApplyHitDelay applies the hit delay to the knockback.
func ApplyHitDelay(hd *time.Duration) {
	*hd = hitDelay
}

// save saves the settings to the file at the path passed.
func save() {
	s := settings{
		Force:    mgl64.Round(force, 3),
		Height:   mgl64.Round(height, 3),
		HitDelay: hitDelay.Milliseconds(),
	}

	_ = os.MkdirAll(filepath.Dir("assets/knockback.json"), 0777)
	_ = gophig.SetConfComplex("assets/knockback.json", gophig.JSONMarshaler{Indent: true}, s, os.ModePerm)
}
