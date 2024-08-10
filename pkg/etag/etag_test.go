package etag

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ParseTag(t *testing.T) {
	t.Parallel()

	t.Run("Valid Raw ETag", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		timestamp := time.Now()
		raw := fmt.Sprintf(
			`"%s%s%s"`,
			id,
			eTagSeparator,
			timestamp.Format(time.RFC3339Nano),
		)

		got, err := Parse(raw)
		assert.NoError(t, err)
		assert.Equal(t, id, got.id)

		assert.Truef(
			t,
			timestamp.Equal(got.updatedAt),
			"expected equal timestamps\n\twant:\t%s\n\tgot:\t%s\n",
			timestamp.Format(time.RFC3339Nano),
			got.updatedAt.Format(time.RFC3339Nano),
		)

		t.Run("errors", func(t *testing.T) {
			t.Parallel()

			testCases := []struct {
				name string
				raw  string
			}{
				{
					name: "not a double-quoted string",
					raw: fmt.Sprintf(
						"%s%s%s",
						uuid.New(),
						eTagSeparator,
						time.Now().Format(time.RFC3339Nano),
					),
				},
				{
					name: "has < 2 components",
					raw:  fmt.Sprintf(`"%s"`, uuid.New()),
				},
				{
					name: "has > 2 components",
					raw: fmt.Sprintf(
						`"%s%s%s%s"`,
						uuid.New(),
						eTagSeparator,
						time.Now().Format(time.RFC3339Nano),
						eTagSeparator,
					),
				},
				{
					name: "has invalid UUID",
					raw: fmt.Sprintf(
						`"%s%s%s"`,
						"not a UUID",
						eTagSeparator,
						time.Now().Format(time.RFC3339Nano),
					),
				},
				{
					name: "has invalid timestamp",
					raw: fmt.Sprintf(
						`"%s%s%s"`,
						uuid.New(),
						eTagSeparator,
						"not a timestamp",
					),
				},
			}

			for _, tc := range testCases {
				tc := tc

				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					tag, err := Parse(tc.raw)

					var parseErr *ParseETagError
					assert.ErrorAs(t, err, &parseErr)
					assert.Empty(t, tag)
				})
			}
		})
	})
}

func Test_ETag_String(t *testing.T) {
	t.Parallel()

	eTag := ETag{
		id:        uuid.New(),
		updatedAt: time.Now(),
	}
	expected := fmt.Sprintf(
		`"%s%s%s"`,
		eTag.id,
		eTagSeparator,
		eTag.updatedAt.Format(time.RFC3339Nano),
	)

	got := eTag.String()
	assert.Equal(t, expected, got)
}

func TestNew(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	timestamp := time.Now()

	etag := New(id, timestamp)

	// test ID
	assert.Equal(t, id, etag.ID())

	// test the updated at time is recent
	assert.WithinDuration(t, time.Now(), etag.UpdatedAt(), 2*time.Second)
	// Test that the string representation is correct
	expectedFormat := fmt.Sprintf(
		`"%s%s%s"`,
		etag.ID(),
		eTagSeparator,
		etag.UpdatedAt().Format(time.RFC3339Nano),
	)
	assert.Equal(t, expectedFormat, etag.String())
}

func TestRandom(t *testing.T) {
	t.Parallel()

	etag := Random()

	// Test that ID is not nil
	assert.NotEqual(t, uuid.Nil, etag.ID())

	// Test that UpdatedAt time is within the last year
	now := time.Now()
	lastYear := now.AddDate(-1, 0, 0)
	assert.True(t, etag.UpdatedAt().After(lastYear))
	assert.True(t, etag.UpdatedAt().Before(now) || etag.UpdatedAt().Equal(now))

	// Test that the string representation is correct
	expectedFormat := fmt.Sprintf(
		`"%s%s%s"`,
		etag.ID(),
		eTagSeparator,
		etag.UpdatedAt().Format(time.RFC3339Nano),
	)
	assert.Equal(t, expectedFormat, etag.String())

	// Test multiple random ETags to ensure they're different
	anotherETag := Random()
	assert.NotEqual(t, etag.ID(), anotherETag.ID())
	assert.NotEqual(t, etag.UpdatedAt(), anotherETag.UpdatedAt())
}
