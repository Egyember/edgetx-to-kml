package layout

import (
	"flag"
)

type Layout struct {
	Gps, Date, Time, Sats, Alt, Mode, Baro int
}

var (
	date = flag.String("ld", "Date", "log file to convert from")
	sats = flag.String("ls", "Sats", "log file to convert from")
	time = flag.String("lt", "Time", "log file to convert from")
	gps  = flag.String("lg", "GPS", "log file to convert from")
	alt  = flag.String("la", "Alt(m)", "log file to convert from")
	baro = flag.String("lb", "Alt2(m)", "log file to convert from")
	mode = flag.String("lm", "FM", "log file to convert from")
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
