package forcematrix_test

import (
	"context"
	"testing"

	"go.viam.com/test"
	"go.viam.com/utils"

	"go.viam.com/rdk/component/arm"
	"go.viam.com/rdk/component/forcematrix"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/testutils/inject"
)

const (
	testForceMatrixName    = "forcematrix1"
	testForceMatrixName2   = "forcematrix2"
	fakeForceMatrixName    = "forcematrix3"
	missingForceMatrixName = "forcematrix4"
)

func setupInjectRobot() *inject.Robot {
	forcematrix1 := &mock{Name: testForceMatrixName}
	r := &inject.Robot{}
	r.ResourceByNameFunc = func(name resource.Name) (interface{}, bool) {
		switch name {
		case forcematrix.Named(testForceMatrixName):
			return forcematrix1, true
		case forcematrix.Named(fakeForceMatrixName):
			return "not a forcematrix", false
		default:
			return nil, false
		}
	}
	r.ResourceNamesFunc = func() []resource.Name {
		return []resource.Name{forcematrix.Named(testForceMatrixName), arm.Named("arm1")}
	}
	return r
}

func TestFromRobot(t *testing.T) {
	r := setupInjectRobot()

	s, ok := forcematrix.FromRobot(r, testForceMatrixName)
	test.That(t, s, test.ShouldNotBeNil)
	test.That(t, ok, test.ShouldBeTrue)

	result, err := s.DetectSlip(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, result, test.ShouldResemble, slip)

	s, ok = forcematrix.FromRobot(r, fakeForceMatrixName)
	test.That(t, s, test.ShouldBeNil)
	test.That(t, ok, test.ShouldBeFalse)

	s, ok = forcematrix.FromRobot(r, missingForceMatrixName)
	test.That(t, s, test.ShouldBeNil)
	test.That(t, ok, test.ShouldBeFalse)
}

func TestNamesFromRobot(t *testing.T) {
	r := setupInjectRobot()

	names := forcematrix.NamesFromRobot(r)
	test.That(t, names, test.ShouldResemble, []string{testForceMatrixName})
}

func TestForceMatrixName(t *testing.T) {
	for _, tc := range []struct {
		TestName string
		Name     string
		Expected resource.Name
	}{
		{
			"missing name",
			"",
			resource.Name{
				UUID: "08174524-a3f0-585d-a7da-9763d9534dd1",
				Subtype: resource.Subtype{
					Type:            resource.Type{Namespace: resource.ResourceNamespaceRDK, ResourceType: resource.ResourceTypeComponent},
					ResourceSubtype: forcematrix.SubtypeName,
				},
				Name: "",
			},
		},
		{
			"all fields included",
			testForceMatrixName,
			resource.Name{
				UUID: "a5f3c7aa-4267-5856-81ae-565e7ad44916",
				Subtype: resource.Subtype{
					Type:            resource.Type{Namespace: resource.ResourceNamespaceRDK, ResourceType: resource.ResourceTypeComponent},
					ResourceSubtype: forcematrix.SubtypeName,
				},
				Name: testForceMatrixName,
			},
		},
	} {
		t.Run(tc.TestName, func(t *testing.T) {
			observed := forcematrix.Named(tc.Name)
			test.That(t, observed, test.ShouldResemble, tc.Expected)
		})
	}
}

func TestWrapWithReconfigurable(t *testing.T) {
	var actualForceMatrix1 forcematrix.ForceMatrix = &mock{Name: testForceMatrixName}
	reconfForceMatrix1, err := forcematrix.WrapWithReconfigurable(actualForceMatrix1)
	test.That(t, err, test.ShouldBeNil)

	_, err = forcematrix.WrapWithReconfigurable(nil)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, "expected resource")

	reconfForceMatrix2, err := forcematrix.WrapWithReconfigurable(reconfForceMatrix1)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, reconfForceMatrix2, test.ShouldEqual, reconfForceMatrix1)
}

func TestReconfigurableForceMatrix(t *testing.T) {
	actualForceMatrix1 := &mock{Name: testForceMatrixName}
	reconfForceMatrix1, err := forcematrix.WrapWithReconfigurable(actualForceMatrix1)
	test.That(t, err, test.ShouldBeNil)

	actualForceMatrix2 := &mock{Name: testForceMatrixName2}
	reconfForceMatrix2, err := forcematrix.WrapWithReconfigurable(actualForceMatrix2)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, actualForceMatrix1.reconfCount, test.ShouldEqual, 0)

	err = reconfForceMatrix1.Reconfigure(context.Background(), reconfForceMatrix2)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, reconfForceMatrix1, test.ShouldResemble, reconfForceMatrix2)
	test.That(t, actualForceMatrix1.reconfCount, test.ShouldEqual, 1)

	test.That(t, actualForceMatrix1.slipCount, test.ShouldEqual, 0)
	test.That(t, actualForceMatrix2.slipCount, test.ShouldEqual, 0)
	result, err := reconfForceMatrix1.(forcematrix.ForceMatrix).DetectSlip(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, result, test.ShouldEqual, slip)
	test.That(t, actualForceMatrix1.slipCount, test.ShouldEqual, 0)
	test.That(t, actualForceMatrix2.slipCount, test.ShouldEqual, 1)

	err = reconfForceMatrix1.Reconfigure(context.Background(), nil)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, "expected *forcematrix.reconfigurableForceMatrix")
}

func TestReadMatrix(t *testing.T) {
	actualForceMatrix1 := &mock{Name: testForceMatrixName}
	reconfForceMatrix1, _ := forcematrix.WrapWithReconfigurable(actualForceMatrix1)

	test.That(t, actualForceMatrix1.matrixCount, test.ShouldEqual, 0)
	matrix1, err := reconfForceMatrix1.(forcematrix.ForceMatrix).ReadMatrix(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, matrix1, test.ShouldResemble, matrix)
	test.That(t, actualForceMatrix1.matrixCount, test.ShouldEqual, 1)
}

func TestDetectSlip(t *testing.T) {
	actualForceMatrix1 := &mock{Name: testForceMatrixName}
	reconfForceMatrix1, _ := forcematrix.WrapWithReconfigurable(actualForceMatrix1)

	test.That(t, actualForceMatrix1.slipCount, test.ShouldEqual, 0)
	slip1, err := reconfForceMatrix1.(forcematrix.ForceMatrix).DetectSlip(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, slip1, test.ShouldEqual, slip)
	test.That(t, actualForceMatrix1.slipCount, test.ShouldEqual, 1)
}

func TestGetReadings(t *testing.T) {
	actualForceMatrix1 := &mock{Name: testForceMatrixName}
	reconfForceMatrix1, _ := forcematrix.WrapWithReconfigurable(actualForceMatrix1)

	test.That(t, actualForceMatrix1.readingsCount, test.ShouldEqual, 0)
	result, err := reconfForceMatrix1.(forcematrix.ForceMatrix).GetReadings(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, result, test.ShouldResemble, []interface{}{matrix})
	test.That(t, actualForceMatrix1.readingsCount, test.ShouldEqual, 1)
}

func TestClose(t *testing.T) {
	actualForceMatrix1 := &mock{Name: testForceMatrixName}
	reconfForceMatrix1, _ := forcematrix.WrapWithReconfigurable(actualForceMatrix1)

	test.That(t, actualForceMatrix1.reconfCount, test.ShouldEqual, 0)
	test.That(t, utils.TryClose(context.Background(), reconfForceMatrix1), test.ShouldBeNil)
	test.That(t, actualForceMatrix1.reconfCount, test.ShouldEqual, 1)
}

var (
	matrix = [][]int{{2, 1}}
	slip   = false
)

type mock struct {
	forcematrix.ForceMatrix
	Name          string
	matrixCount   int
	slipCount     int
	readingsCount int
	reconfCount   int
}

// ReadMatrix returns the set value.
func (m *mock) ReadMatrix(ctx context.Context) ([][]int, error) {
	m.matrixCount++
	return matrix, nil
}

// DetectSlip returns the set value.
func (m *mock) DetectSlip(ctx context.Context) (bool, error) {
	m.slipCount++
	return slip, nil
}

func (m *mock) GetReadings(ctx context.Context) ([]interface{}, error) {
	m.readingsCount++
	return []interface{}{matrix}, nil
}

func (m *mock) Close() { m.reconfCount++ }
