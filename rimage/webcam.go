package rimage

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"regexp"

	"github.com/edaniels/golog"
	"github.com/edaniels/gostream"
	"github.com/edaniels/gostream/media"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"go.uber.org/multierr"
)

var cameraConstraints = mediadevices.MediaStreamConstraints{
	Video: func(constraint *mediadevices.MediaTrackConstraints) {
		constraint.Width = prop.IntRanged{640, 4096, 1920}
		constraint.Height = prop.IntRanged{400, 2160, 1080}
		constraint.FrameRate = prop.FloatRanged{0, 200, 60}
		constraint.FrameFormat = prop.FrameFormatOneOf{
			frame.FormatI420,
			frame.FormatI444,
			frame.FormatYUY2,
			frame.FormatUYVY,
			frame.FormatRGBA,
			frame.FormatMJPEG,
			frame.FormatNV12,
			frame.FormatNV21, // gives blue tinted image?
		}
	},
}

var intelConstraints = mediadevices.MediaStreamConstraints{
	Video: func(constraint *mediadevices.MediaTrackConstraints) {
		constraint.Width = prop.IntRanged{1280, 1280, 1280}
		constraint.Height = prop.IntRanged{720, 720, 720}
		constraint.FrameRate = prop.FloatRanged{0, 200, 60}
		constraint.FrameFormat = prop.FrameFormatOneOf{
			frame.FormatI420,
			frame.FormatI444,
			frame.FormatYUY2,
			frame.FormatUYVY,
			frame.FormatRGBA,
			frame.FormatMJPEG,
			frame.FormatNV12,
			frame.FormatNV21, // gives blue tinted image?
		}
	},
}

var depthConstraints = mediadevices.MediaStreamConstraints{
	Video: func(constraint *mediadevices.MediaTrackConstraints) {
		constraint.Width = prop.IntRanged{640, 4096, 1920}
		constraint.Height = prop.IntRanged{400, 2160, 1080}
		constraint.FrameRate = prop.FloatRanged{0, 200, 60}
		constraint.FrameFormat = prop.FrameFormatExact(frame.FormatZ16)
	},
}

type Aligner interface {
	Align(img *ImageWithDepth) (*ImageWithDepth, error)
}

type WebcamSource struct {
	reader media.VideoReadCloser
	depth  *webcamDepthSource
	align  Aligner
}

func maybeAddDepth(s *WebcamSource, debug, isIntel515 bool, attrs map[string]string) (*WebcamSource, error) {
	var err error
	s.depth, err = findWebcamDepth(debug)
	if isIntel515 {
		if err != nil {
			return nil, fmt.Errorf("found intel camera point no matching depth: %w", err)
		}
		s.align = &Intel515Align{}
	}

	return s, nil
}

func NewWebcamSource(attrs map[string]string) (gostream.ImageSource, error) {
	debug := attrs["debug"] == "true"
	// TODO(erh): this is gross, but not sure what right config options are yet
	isIntel515 := attrs["model"] == "intel515"

	constraints := cameraConstraints

	if isIntel515 {
		constraints = intelConstraints
	}

	path := attrs["path"]
	if path != "" {
		return tryWebcamOpen(path, debug, constraints)
	}

	pathPattern := attrs["path_pattern"]
	if pathPattern != "" {
		pattern, err := regexp.Compile(pathPattern)
		if err != nil {
			return nil, err
		}
		s, err := tryWebcamOpenPattern(pattern, debug, constraints)
		if err != nil {
			return nil, err
		}
		return maybeAddDepth(s, debug, isIntel515, attrs)
	}

	labels := media.QueryVideoDeviceLabels()
	for _, label := range labels {
		if debug {
			golog.Global.Debugf("%s", label)
		}

		s, err := tryWebcamOpen(label, debug, constraints)
		if err == nil {
			if debug {
				golog.Global.Debugf("\t USING")
			}

			return maybeAddDepth(s, debug, isIntel515, attrs)
		}
		if debug {
			golog.Global.Debugf("\t %s", err)
		}
	}

	return nil, fmt.Errorf("could find no webcams")
}

func (s *WebcamSource) Next(ctx context.Context) (image.Image, error) {
	img, err := media.VideoReadReleaser{s.reader}.Read()
	if err != nil {
		return nil, err
	}

	if s.depth == nil {
		return img, nil
	}

	dm, err := s.depth.Next()
	if err != nil {
		return nil, err
	}

	iwd, err := s.align.Align(&ImageWithDepth{ConvertImage(img), dm})
	if err != nil {
		return nil, err
	}

	return iwd, nil
}

func (s *WebcamSource) Close() error {
	var errs []error
	if s.depth != nil {
		if err := s.depth.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if err := s.reader.Close(); err != nil {
		errs = append(errs, err)
	}
	return multierr.Combine(errs...)
}

func tryWebcamOpen(path string, debug bool, constraints mediadevices.MediaStreamConstraints) (*WebcamSource, error) {
	reader, err := media.GetNamedVideoReader(filepath.Base(path), constraints)
	if err != nil {
		return nil, err
	}

	return &WebcamSource{reader, nil, nil}, nil
}

func tryWebcamOpenPattern(pattern *regexp.Regexp, debug bool, constraints mediadevices.MediaStreamConstraints) (*WebcamSource, error) {
	reader, err := media.GetPatternedVideoReader(pattern, constraints)
	if err != nil {
		return nil, err
	}

	return &WebcamSource{reader, nil, nil}, nil
}

type webcamDepthSource struct {
	reader media.VideoReadCloser
}

func findWebcamDepth(debug bool) (*webcamDepthSource, error) {
	reader, err := media.GetAnyVideoReader(depthConstraints)
	if debug {
		golog.Global.Debugf("depth cam: %v", reader)
	}
	if err != nil {
		return nil, fmt.Errorf("no depth camera found: %w", err)
	}

	return &webcamDepthSource{reader}, nil
}

func (w *webcamDepthSource) Next() (*DepthMap, error) {
	img, err := media.VideoReadReleaser{w.reader}.Read()
	if err != nil {
		return nil, err
	}

	return imageToDepthMap(img), nil
}

func (w *webcamDepthSource) Close() error {
	return w.reader.Close()
}

func imageToDepthMap(img image.Image) *DepthMap {
	bounds := img.Bounds()

	width, height := bounds.Dx(), bounds.Dy()

	// TODO: handle non realsense Z16 devices better
	// realsense seems to rotate
	dm := NewEmptyDepthMap(width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			z := color.Gray16Model.Convert(img.At(x, y)).(color.Gray16).Y
			dm.Set(x, y, int(z))
		}
	}

	return dm.Rotate90(false)
}
