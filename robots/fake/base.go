package fake

import "context"

// tracks in CM
type Base struct {
}

func (b *Base) MoveStraight(ctx context.Context, distanceMillis int, millisPerSec float64, block bool) error {
	return nil
}

func (b *Base) Spin(ctx context.Context, angleDeg float64, speed int, block bool) error {
	return nil
}

func (b *Base) WidthMillis(ctx context.Context) (int, error) {
	return 600, nil
}

func (b *Base) Stop(ctx context.Context) error {
	return nil
}

func (b *Base) Close(ctx context.Context) error {
	return nil
}
