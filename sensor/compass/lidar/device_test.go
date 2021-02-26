package lidar

import (
	"context"
	"errors"
	"fmt"
	"image"
	"math"
	"testing"

	"go.viam.com/robotcore/lidar"
	"go.viam.com/robotcore/pc"
	"go.viam.com/robotcore/sensor/compass"
	"go.viam.com/robotcore/testutils"
	"go.viam.com/robotcore/utils"

	"github.com/edaniels/test"
	"gonum.org/v1/gonum/mat"
)

func TestNew(t *testing.T) {
	// unknown type
	_, err := New(context.Background(), lidar.DeviceDescription{Type: "what"})
	test.That(t, err, test.ShouldNotBeNil)

	devType := lidar.DeviceType(utils.RandomAlphaString(5))
	var newFunc func(desc lidar.DeviceDescription) (lidar.Device, error)
	lidar.RegisterDeviceType(devType, lidar.DeviceTypeRegistration{
		New: func(ctx context.Context, desc lidar.DeviceDescription) (lidar.Device, error) {
			return newFunc(desc)
		},
	})

	desc := lidar.DeviceDescription{Type: devType, Path: "somewhere"}
	newErr := errors.New("woof")
	newFunc = func(innerDesc lidar.DeviceDescription) (lidar.Device, error) {
		test.That(t, innerDesc, test.ShouldResemble, desc)
		return nil, newErr
	}

	_, err = New(context.Background(), desc)
	test.That(t, err, test.ShouldEqual, newErr)

	injectDev := &injectDevice{}
	newFunc = func(innerDesc lidar.DeviceDescription) (lidar.Device, error) {
		return injectDev, nil
	}

	dev, err := New(context.Background(), desc)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, dev, test.ShouldNotBeNil)
}

func TestFrom(t *testing.T) {
	dev := &injectDevice{}
	compassDev := From(dev)
	var relDev *compass.RelativeDevice = nil
	test.That(t, compassDev, test.ShouldImplement, relDev)
}

func getInjected() (*Device, *injectDevice) {
	dev := &injectDevice{}
	return From(dev).(*Device), dev
}

func TestCompass(t *testing.T) {
	t.Run("StartCalibration", func(t *testing.T) {
		compassDev, _ := getInjected()
		test.That(t, compassDev.StartCalibration(context.Background()), test.ShouldBeNil)
	})

	t.Run("StopCalibration", func(t *testing.T) {
		compassDev, _ := getInjected()
		test.That(t, compassDev.StopCalibration(context.Background()), test.ShouldBeNil)
	})
}

func TestScanToVec2Matrix(t *testing.T) {
	t.Run("with no results should produce an empty matrix", func(t *testing.T) {
		compassDev, injectDev := getInjected()
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			return nil, nil
		}
		m, err := compassDev.scanToVec2Matrix()
		test.That(t, err, test.ShouldBeNil)
		test.That(t, (*mat.Dense)(m).IsEmpty(), test.ShouldBeTrue)
	})

	t.Run("should request a scan with more than 1 count and no filtering", func(t *testing.T) {
		compassDev, injectDev := getInjected()
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			test.That(t, options.Count, test.ShouldBeGreaterThan, 1)
			test.That(t, options.NoFilter, test.ShouldBeTrue)
			return nil, nil
		}
		m, err := compassDev.scanToVec2Matrix()
		test.That(t, err, test.ShouldBeNil)
		test.That(t, (*mat.Dense)(m).IsEmpty(), test.ShouldBeTrue)
	})

	t.Run("should error out if all of the scans fail", func(t *testing.T) {
		compassDev, injectDev := getInjected()
		count := 0
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			count++
			return nil, errors.New("oops")
		}
		_, err := compassDev.scanToVec2Matrix()
		test.That(t, err, test.ShouldBeError, "oops")
		test.That(t, count, test.ShouldBeGreaterThan, 1)
	})

	t.Run("should convert measurments into a matrix", func(t *testing.T) {
		compassDev, injectDev := getInjected()
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			return lidar.Measurements{
				lidar.NewMeasurement(1, 10),
				lidar.NewMeasurement(20, 2),
				lidar.NewMeasurement(30, 5),
			}, nil
		}
		m, err := compassDev.scanToVec2Matrix()
		test.That(t, err, test.ShouldBeNil)
		mD := (*mat.Dense)(m)
		test.That(t, mD.IsEmpty(), test.ShouldBeFalse)
		r, c := mD.Dims()
		test.That(t, r, test.ShouldEqual, 3)
		test.That(t, c, test.ShouldEqual, 3)
		test.That(t, mD.RawRowView(0), test.ShouldResemble, []float64{
			0.17452406437283513, 0.6840402866513374, 2.4999999999999996,
		}) // x
		test.That(t, mD.RawRowView(1), test.ShouldResemble, []float64{
			-9.998476951563912, -1.8793852415718169, -4.330127018922194,
		}) // y
		test.That(t, mD.RawRowView(2), test.ShouldResemble, []float64{1, 1, 1}) // fill
	})
}

func TestHeading(t *testing.T) {
	t.Run("with no results should NaN", func(t *testing.T) {
		compassDev, injectDev := getInjected()
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			return nil, nil
		}
		h, err := compassDev.Heading(context.Background())
		test.That(t, err, test.ShouldBeNil)
		test.That(t, math.IsNaN(h), test.ShouldBeTrue)
	})

	t.Run("with some results should NaN without mark", func(t *testing.T) {
		compassDev, injectDev := getInjected()
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			return lidar.Measurements{
				lidar.NewMeasurement(1, 10),
				lidar.NewMeasurement(20, 2),
				lidar.NewMeasurement(30, 5),
			}, nil
		}
		h, err := compassDev.Heading(context.Background())
		test.That(t, err, test.ShouldBeNil)
		test.That(t, math.IsNaN(h), test.ShouldBeTrue)
	})

	t.Run("with mark", func(t *testing.T) {
		pointCloud, err := pc.NewPointCloudFromFile(testutils.ResolveFile("pc/data/test.las"))
		test.That(t, err, test.ShouldBeNil)

		mat2 := pointCloud.ToVec2Matrix()
		firstMs := lidar.MeasurementsFromVec2Matrix(mat2)
		compassDev, injectDev := getInjected()
		angularRes := .3375
		injectDev.AngularResolutionFunc = func(ctx context.Context) (float64, error) {
			return angularRes, nil
		}
		injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
			return firstMs, nil
		}
		test.That(t, compassDev.Mark(), test.ShouldBeNil)

		scannedM, err := compassDev.scanToVec2Matrix()
		test.That(t, err, test.ShouldBeNil)

		setup := func(t *testing.T) (*Device, *injectDevice) {
			t.Helper()
			_, injectDev := getInjected()
			injectDev.AngularResolutionFunc = func(ctx context.Context) (float64, error) {
				return angularRes, nil
			}
			injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
				return firstMs, nil
			}
			cloned := compassDev.clone()
			cloned.setDevice(injectDev)
			return cloned, injectDev
		}

		t.Run("heading should be 0", func(t *testing.T) {
			compassDev, _ := setup(t)
			heading, err := compassDev.Heading(context.Background())
			test.That(t, err, test.ShouldBeNil)
			test.That(t, heading, test.ShouldEqual, 0)
		})

		for i := 0; i < 360; i++ {
			iCopy := i
			t.Run(fmt.Sprintf("rotating %d heading should be %d", iCopy, iCopy), func(t *testing.T) {
				t.Parallel()
				compassDev, injectDev := setup(t)
				rot := scannedM.RotateMatrixAbout(0, 0, float64(iCopy))
				rotM := lidar.MeasurementsFromVec2Matrix(rot)

				injectDev.ScanFunc = func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
					return rotM, nil
				}

				heading, err := compassDev.Heading(context.Background())
				test.That(t, err, test.ShouldBeNil)
				test.That(t, heading, test.ShouldEqual, iCopy)
			})
		}
	})
}

type injectDevice struct {
	lidar.Device
	StartFunc             func(ctx context.Context) error
	StopFunc              func(ctx context.Context) error
	CloseFunc             func(ctx context.Context) error
	ScanFunc              func(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error)
	RangeFunc             func(ctx context.Context) (int, error)
	BoundsFunc            func(ctx context.Context) (image.Point, error)
	AngularResolutionFunc func(ctx context.Context) (float64, error)
}

func (ij *injectDevice) Start(ctx context.Context) error {
	if ij.StartFunc == nil {
		return ij.Device.Start(ctx)
	}
	return ij.StartFunc(ctx)
}

func (ij *injectDevice) Stop(ctx context.Context) error {
	if ij.StopFunc == nil {
		return ij.Device.Stop(ctx)
	}
	return ij.StopFunc(ctx)
}

func (ij *injectDevice) Close(ctx context.Context) error {
	if ij.CloseFunc == nil {
		return ij.Device.Close(ctx)
	}
	return ij.CloseFunc(ctx)
}

func (ij *injectDevice) Scan(ctx context.Context, options lidar.ScanOptions) (lidar.Measurements, error) {
	if ij.ScanFunc == nil {
		return ij.Device.Scan(ctx, options)
	}
	return ij.ScanFunc(ctx, options)
}

func (ij *injectDevice) Range(ctx context.Context) (int, error) {
	if ij.RangeFunc == nil {
		return ij.Device.Range(ctx)
	}
	return ij.RangeFunc(ctx)
}

func (ij *injectDevice) Bounds(ctx context.Context) (image.Point, error) {
	if ij.BoundsFunc == nil {
		return ij.Device.Bounds(ctx)
	}
	return ij.BoundsFunc(ctx)
}

func (ij *injectDevice) AngularResolution(ctx context.Context) (float64, error) {
	if ij.AngularResolutionFunc == nil {
		return ij.Device.AngularResolution(ctx)
	}
	return ij.AngularResolutionFunc(ctx)
}
