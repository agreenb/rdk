package gpsnmea

import (
	"math"
	"strconv"
	"strings"

	"github.com/adrianmo/go-nmea"
	geo "github.com/kellydunn/golang-geo"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

const (
	knotsToMPerSec = 0.51444
	kphToMPerSec   = 0.27778
)

// GPSData struct combines various attributes related to GPS.
type GPSData struct {
	Location   *geo.Point
	Alt        float64
	Speed      float64 // ground speed in m per sec
	VDOP       float64 // vertical accuracy
	HDOP       float64 // horizontal accuracy
	SatsInView int     // quantity satellites in view
	SatsInUse  int     // quantity satellites in view
	valid      bool
	FixQuality int
}

func errInvalidFix(sentenceType, badFix, goodFix string) error {
	return errors.Errorf("type %q sentence fix is not valid have: %q  want %q", sentenceType, badFix, goodFix)
}

// ParseAndUpdate will attempt to parse a line to an NMEA sentence, and if valid, will try to update the given struct
// with the values for that line. Nothing will be updated if there is not a valid gps fix.
func (g *GPSData) ParseAndUpdate(line string) error {
	// add parsing to filter out corrupted data
	ind := strings.Index(line, "$G")
	if ind == -1 {
		line = ""
	} else {
		line = line[ind:]
	}

	var errs error
	s, err := nmea.Parse(line)
	if err != nil {
		return multierr.Combine(errs, err)
	}
	// Most receivers support at least the following sentence types: GSV, RMC, GSA, GGA, GLL, VTG, GNS
	if gsv, ok := s.(nmea.GSV); ok {
		// GSV provides the number of satellites in view
		g.SatsInView = int(gsv.NumberSVsInView)
	} else if rmc, ok := s.(nmea.RMC); ok {
		// RMC provides validity, lon/lat, and ground speed.
		if rmc.Validity == "A" {
			g.valid = true
		} else if rmc.Validity == "V" {
			g.valid = false
			errs = multierr.Combine(errs, errInvalidFix(rmc.Type, rmc.Validity, "A"))
		}
		if g.valid {
			g.Speed = rmc.Speed * knotsToMPerSec
			g.Location = geo.NewPoint(rmc.Latitude, rmc.Longitude)
		}
	} else if gsa, ok := s.(nmea.GSA); ok {
		// GSA gives horizontal and vertical accuracy, and also describes the type of lock- invalid, 2d, or 3d.
		switch gsa.FixType {
		case "1":
			// No fix
			g.valid = false
			errs = multierr.Combine(errs, errInvalidFix(gsa.Type, gsa.FixType, "1 or 2"))
		case "2":
			// 2d fix, valid lat/lon but invalid alt
			g.valid = true
			g.VDOP = -1
		case "3":
			// 3d fix
			g.valid = true
		}
		if g.valid {
			g.VDOP = gsa.VDOP
			g.HDOP = gsa.HDOP
		}
		g.SatsInUse = len(gsa.SV)
	} else if gga, ok := s.(nmea.GGA); ok {
		// GGA provides validity, lon/lat, altitude, sats in use, and horizontal position error
		g.FixQuality, err = strconv.Atoi(gga.FixQuality)
		if err != nil {
			return err
		}
		if gga.FixQuality == "0" {
			g.valid = false
			errs = multierr.Combine(errs, errInvalidFix(gga.Type, gga.FixQuality, "1 to 6"))
		} else {
			g.valid = true
			g.Location = geo.NewPoint(gga.Latitude, gga.Longitude)
			g.SatsInUse = int(gga.NumSatellites)
			g.HDOP = gga.HDOP
			g.Alt = gga.Altitude
		}
	} else if gll, ok := s.(nmea.GLL); ok {
		// GLL provides just lat/lon
		now := toPoint(gll)
		g.Location = now
	} else if vtg, ok := s.(nmea.VTG); ok {
		// VTG provides ground speed
		g.Speed = vtg.GroundSpeedKPH * kphToMPerSec
	} else if gns, ok := s.(nmea.GNS); ok {
		// GNS Provides approximately the same information as GGA
		for _, mode := range gns.Mode {
			if mode == "N" {
				g.valid = false
				errs = multierr.Combine(errs, errInvalidFix(gns.Type, mode, " A, D, P, R, F, E, M or S"))
			}
		}
		if g.valid {
			g.Location = geo.NewPoint(gns.Latitude, gns.Longitude)
			g.SatsInUse = int(gns.SVs)
			g.HDOP = gns.HDOP
			g.Alt = gns.Altitude
		}
	}

	if g.Location == nil {
		g.Location = geo.NewPoint(math.NaN(), math.NaN())
		errs = multierr.Combine(errs, errors.New("no location parsed for nmea gps, using default value of lat: NaN, long: NaN"))
		return errs
	}
	return nil
}
