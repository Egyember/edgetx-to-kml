package layout

import (
	"flag"
)

type Layout struct {
	Gps, Date, Time, Sats, Alt, Mode, Baro int
}

var (
	date = flag.String("ld", "Date", "data in header")
	sats = flag.String("ls", "Sats", "satelite count in header")
	time = flag.String("lt", "Time", "time in header")
	gps  = flag.String("lg", "GPS", "gps cordinates in header")
	alt  = flag.String("la", "Alt(m)", "gps altitude in header")
	baro = flag.String("lb", "Alt2(m)", "barometer altitude in header")
	mode = flag.String("lm", "FM", "flight mode in header")
)

func GetLayout(header []string) (r Layout) {
	for k, v := range header {
		switch v {
		case *date:
			r.Date = k
		case *sats:
			r.Sats = k
		case *time:
			r.Time = k
		case *gps:
			r.Gps = k
		case *alt:
			r.Alt = k
		case *mode:
			r.Mode = k
		case *baro:
			r.Baro = k
		}
	}
	return
}
