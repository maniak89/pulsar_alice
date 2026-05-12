package forward

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	ctx := context.Background()
	client := New(Config{
		Address:  "https://simfor.ru",
		Login:    "maniak89s@gmail.com",
		Password: "5vg5suqCAU2Absdf",
		Timout:   time.Second * 5,
	})

	session, err := client.StartSession(ctx)
	require.NoError(t, err)

	require.NoError(t, session.SendMetrics(ctx, 25.8, 13.9))
}
