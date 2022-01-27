package gantry_test

import (
	"context"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/component/gantry"
	pb "go.viam.com/rdk/proto/api/component/v1"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/subtype"
	"go.viam.com/rdk/testutils/inject"
)

func newServer() (pb.GantryServiceServer, *inject.Gantry, *inject.Gantry, error) {
	injectGantry := &inject.Gantry{}
	injectGantry2 := &inject.Gantry{}
	gantries := map[resource.Name]interface{}{
		gantry.Named("gantry1"): injectGantry,
		gantry.Named("gantry2"): injectGantry2,
		gantry.Named("gantry3"): "notGantry",
	}
	gantrySvc, err := subtype.New(gantries)
	if err != nil {
		return nil, nil, nil, err
	}
	return gantry.NewServer(gantrySvc), injectGantry, injectGantry2, nil
}

func TestServer(t *testing.T) {
	gantryServer, injectGantry, injectGantry2, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	var gantryPos []float64

	gantry1 := "gantry1"
	pos1 := []float64{1.0, 2.0, 3.0}
	len1 := []float64{2.0, 3.0, 4.0}
	injectGantry.GetPositionFunc = func(ctx context.Context) ([]float64, error) {
		return pos1, nil
	}
	injectGantry.MoveToPositionFunc = func(ctx context.Context, pos []float64) error {
		gantryPos = pos
		return nil
	}
	injectGantry.GetLengthsFunc = func(ctx context.Context) ([]float64, error) {
		return len1, nil
	}

	gantry2 := "gantry2"
	pos2 := []float64{4.0, 5.0, 6.0}
	len2 := []float64{5.0, 6.0, 7.0}
	injectGantry2.GetPositionFunc = func(ctx context.Context) ([]float64, error) {
		return pos2, nil
	}
	injectGantry2.MoveToPositionFunc = func(ctx context.Context, pos []float64) error {
		gantryPos = pos
		return nil
	}
	injectGantry2.GetLengthsFunc = func(ctx context.Context) ([]float64, error) {
		return len2, nil
	}

	t.Run("gantry position", func(t *testing.T) {
		_, err := gantryServer.GetPosition(context.Background(), &pb.GantryServiceGetPositionRequest{Name: "g4"})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, "no gantry")

		_, err = gantryServer.GetPosition(context.Background(), &pb.GantryServiceGetPositionRequest{Name: "gantry3"})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, "not a gantry")

		resp, err := gantryServer.GetPosition(context.Background(), &pb.GantryServiceGetPositionRequest{Name: gantry1})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, resp.PositionsMm, test.ShouldResemble, pos1)

		resp, err = gantryServer.GetPosition(context.Background(), &pb.GantryServiceGetPositionRequest{Name: gantry2})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, resp.PositionsMm, test.ShouldResemble, pos2)
	})

	t.Run("move to position", func(t *testing.T) {
		_, err := gantryServer.MoveToPosition(context.Background(), &pb.GantryServiceMoveToPositionRequest{Name: "g4", PositionsMm: pos2})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, "no gantry")

		_, err = gantryServer.MoveToPosition(context.Background(), &pb.GantryServiceMoveToPositionRequest{Name: gantry1, PositionsMm: pos2})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, gantryPos, test.ShouldResemble, pos2)

		_, err = gantryServer.MoveToPosition(context.Background(), &pb.GantryServiceMoveToPositionRequest{Name: gantry2, PositionsMm: pos1})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, gantryPos, test.ShouldResemble, pos1)
	})

	t.Run("lengths", func(t *testing.T) {
		_, err := gantryServer.GetLengths(context.Background(), &pb.GantryServiceGetLengthsRequest{Name: "g4"})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, "no gantry")

		resp, err := gantryServer.GetLengths(context.Background(), &pb.GantryServiceGetLengthsRequest{Name: gantry1})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, resp.LengthsMm, test.ShouldResemble, len1)

		resp, err = gantryServer.GetLengths(context.Background(), &pb.GantryServiceGetLengthsRequest{Name: gantry2})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, resp.LengthsMm, test.ShouldResemble, len2)
	})
}
