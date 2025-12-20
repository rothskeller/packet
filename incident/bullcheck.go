package incident

import (
	"log/slog"
	"time"
)

// BulletinAreasToCheck returns the list of bulletin areas that are overdue to
// be checked for new bulletins.
func (inc *Incident) BulletinAreasToCheck() (areas []string) {
	var now = time.Now()

	for area, freq := range inc.Config.BulletinChecks {
		last := inc.BulletinChecks[area]
		delta := now.Sub(last)
		if (freq.Duration != 0 && delta >= freq.Duration) || last.IsZero() {
			areas = append(areas, area)
		}
	}
	return areas
}

// BullletinAreaChecked marks the named bulletin area as having been checked at
// the specified time (usually now, but sometimes time.Time{} in order to force
// a check).
func (inc *Incident) BulletinAreaChecked(area string, at time.Time) {
	inc.BulletinChecks[area] = at
	if inc.Config.BulletinChecks[area].Duration == 0 {
		delete(inc.Config.BulletinChecks, area)
	}
	slog.Info("recorded bulletin area check", "inc", inc.Dir, "area", area, "at", at.Format("2006-01-02T15:04:05"))
}
