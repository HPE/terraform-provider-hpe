package paging_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

// recorder is a PageFunc that serves items from a fixed slice and records how
// it was called.
//
// Pages are fetched concurrently, so it is safe for concurrent use and records
// offsets as a set rather than a sequence — the order they arrive in is not
// something the caller controls or should depend on.
type recorder struct {
	items []int
	total int64

	calls    atomic.Int64
	inFlight atomic.Int64
	maxPara  atomic.Int64
}

func newRecorder(items []int, total int64) *recorder {
	return &recorder{items: items, total: total}
}

func (r *recorder) fn() paging.PageFunc[int] {
	return func(_ context.Context, offset, max int64) ([]int, int64, error) {
		n := r.inFlight.Add(1)
		defer r.inFlight.Add(-1)

		for {
			peak := r.maxPara.Load()
			if n <= peak || r.maxPara.CompareAndSwap(peak, n) {
				break
			}
		}

		r.calls.Add(1)

		if offset >= int64(len(r.items)) {
			return nil, r.total, nil
		}

		end := offset + max
		if end > int64(len(r.items)) {
			end = int64(len(r.items))
		}

		return r.items[offset:end], r.total, nil
	}
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}

	return out
}

// contiguous reports whether the collected values are 0..n-1 in order, which is
// what the fixture serves. It is the assertion that matters most for concurrent
// paging: pages assembled out of order, or an offset miscalculated, shows up
// here as a gap, a duplicate or a reordering.
func contiguous(got []int) bool {
	for i, v := range got {
		if v != i {
			return false
		}
	}

	return true
}

func TestCollect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		items    []int
		total    int64
		pageSize int64
		wantLen  int
	}{
		{name: "empty result", items: nil, total: 0, pageSize: 10, wantLen: 0},
		{name: "single partial page", items: seq(3), total: 3, pageSize: 10, wantLen: 3},
		{
			// The case that catches an off-by-one: the final page is full, so
			// the walk must stop on total rather than waiting for a short page
			// that never comes.
			name:  "exact multiple of page size",
			items: seq(20), total: 20, pageSize: 10, wantLen: 20,
		},
		{name: "several pages, ragged last", items: seq(25), total: 25, pageSize: 10, wantLen: 25},
		{
			// More pages than one wave holds, so the wave loop itself is
			// exercised rather than a single pass.
			name:  "spans several waves",
			items: seq(95), total: 95, pageSize: 10, wantLen: 95,
		},
		{
			// total under-reports. The walk stops when the server says it is
			// done, which is the documented contract.
			name: "total too low", items: seq(25), total: 5, pageSize: 10, wantLen: 10,
		},
		{
			// total over-reports, which without the short-page guard would ask
			// forever for records that do not exist.
			name: "total too high", items: seq(12), total: 9999, pageSize: 10, wantLen: 12,
		},
		{name: "total absent", items: seq(15), total: 0, pageSize: 10, wantLen: 15},
		{
			name:  "total absent, spans waves",
			items: seq(85), total: 0, pageSize: 10, wantLen: 85,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(tt.items, tt.total)

			got, err := paging.CollectWithPageSize(
				context.Background(), tt.pageSize, r.fn())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != tt.wantLen {
				t.Errorf("collected %d items, want %d", len(got), tt.wantLen)
			}

			if !contiguous(got) {
				t.Errorf("collected %v, which is not in the order served — "+
					"pages were assembled wrongly", got)
			}
		})
	}
}

// A wrong total must not cost a request per phantom page. Waves bound the waste
// to a constant, whatever the total claims.
func TestCollectBoundsWasteOnWildlyWrongTotal(t *testing.T) {
	t.Parallel()

	r := newRecorder(seq(12), 9999)

	got, err := paging.CollectWithPageSize(context.Background(), 10, r.fn())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 12 {
		t.Fatalf("collected %d items, want 12", len(got))
	}

	// One first page, then a wave that reaches the short page, then a wave that
	// confirms the end by coming back empty. The second wave is the price of
	// not assuming a short page is the last one — an endpoint returning fewer
	// records than asked for is a different thing from an endpoint running out,
	// and only an empty page tells them apart.
	//
	// What matters is that this is a constant. Believing the total would have
	// meant a thousand requests.
	if maxCalls := int64(1 + 2*paging.DefaultConcurrency); r.calls.Load() > maxCalls {
		t.Errorf("made %d requests for 12 records, want at most %d — "+
			"the reported total of 9999 was trusted", r.calls.Load(), maxCalls)
	}
}

// Pages within a wave must genuinely overlap, or the walk is only sequential
// code wearing a mutex.
//
// An in-memory fixture returns so fast that goroutines need not overlap by
// chance, so each call waits on a barrier until its siblings arrive. If the
// pages were fetched one at a time the barrier would never fill and the test
// would time out rather than quietly pass.
func TestCollectFetchesPagesConcurrently(t *testing.T) {
	t.Parallel()

	const pageSize = 10

	var (
		peak    atomic.Int64
		inFlite atomic.Int64
		arrived atomic.Int64
		barrier = make(chan struct{})
		once    sync.Once
	)

	all := seq(200)

	page := func(_ context.Context, offset, max int64) ([]int, int64, error) {
		n := inFlite.Add(1)
		defer inFlite.Add(-1)

		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}

		// The first page is fetched on its own, before any wave, so it must not
		// wait for siblings that do not exist.
		if offset > 0 {
			// Release everyone once a full wave has arrived. Counting arrivals
			// rather than concurrent occupancy is deliberate: goroutines can be
			// scheduled such that the in-flight count never reaches the wave
			// size even though every page was launched together.
			if arrived.Add(1) >= int64(paging.DefaultConcurrency) {
				once.Do(func() { close(barrier) })
			}

			select {
			case <-barrier:
			case <-time.After(2 * time.Second):
				// The wave never filled; fail on the assertion below rather
				// than hang for the whole test timeout.
			}
		}

		if offset >= int64(len(all)) {
			return nil, 200, nil
		}

		end := offset + max
		if end > int64(len(all)) {
			end = int64(len(all))
		}

		return all[offset:end], 200, nil
	}

	got, err := paging.CollectWithPageSize(context.Background(), pageSize, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 200 || !contiguous(got) {
		t.Errorf("collected %d items, want 200 in order", len(got))
	}

	if peak.Load() < 2 {
		t.Errorf("peak concurrency was %d, want at least 2 — "+
			"pages are not being fetched in parallel", peak.Load())
	}

	if peak.Load() > int64(paging.DefaultConcurrency) {
		t.Errorf("peak concurrency was %d, want at most %d — "+
			"the wave is not bounded", peak.Load(), paging.DefaultConcurrency)
	}
}

// An endpoint that caps a response returns fewer records than asked for without
// being at the end. The offsets after such a page were computed assuming it was
// full, so they must be discarded and re-read rather than skipped over.
func TestCollectRecoversFromShortPageMidWalk(t *testing.T) {
	t.Parallel()

	const (
		pageSize = 10
		capped   = 3
		want     = 40
	)

	var calls atomic.Int64

	all := seq(want)

	// Serves at most `capped` records per call, however many are asked for.
	page := func(_ context.Context, offset, _ int64) ([]int, int64, error) {
		calls.Add(1)

		if offset >= int64(len(all)) {
			return nil, want, nil
		}

		end := offset + capped
		if end > int64(len(all)) {
			end = int64(len(all))
		}

		return all[offset:end], want, nil
	}

	got, err := paging.CollectWithPageSize(context.Background(), pageSize, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != want {
		t.Fatalf("collected %d items, want %d — short pages were skipped over",
			len(got), want)
	}

	if !contiguous(got) {
		t.Errorf("collected %v, which has gaps or duplicates", got)
	}
}

func TestCollectPropagatesError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")

	var calls atomic.Int64

	page := func(_ context.Context, offset, max int64) ([]int, int64, error) {
		if calls.Add(1) > 1 {
			return nil, 0, sentinel
		}

		return seq(int(max)), 1000, nil
	}

	got, err := paging.CollectWithPageSize(context.Background(), 10, page)
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want %v", err, sentinel)
	}

	if got != nil {
		t.Errorf("collected %v on error, want nil — partial results must not "+
			"be returned", got)
	}
}

// One page failing must not leave its siblings running.
func TestCollectCancelsSiblingsOnError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")

	var (
		calls    atomic.Int64
		observed atomic.Int64
	)

	page := func(ctx context.Context, offset, max int64) ([]int, int64, error) {
		// First call establishes a full page so a wave is launched.
		if calls.Add(1) == 1 {
			return seq(int(max)), 1_000_000, nil
		}

		// One page fails; the rest report whether they saw the cancellation.
		if offset == max {
			return nil, 0, sentinel
		}

		<-ctx.Done()
		observed.Add(1)

		return nil, 0, ctx.Err()
	}

	if _, err := paging.CollectWithPageSize(
		context.Background(), 10, page,
	); !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want the sentinel", err)
	}

	if observed.Load() == 0 {
		t.Error("no sibling page observed the cancellation")
	}
}

// When one page fails, the rest of the wave is cancelled — so a sibling often
// carries context.Canceled. The real failure is the one worth reporting;
// surfacing a cancellation it caused would hide why the walk stopped.
func TestCollectPrefersTheRealErrorOverACancellation(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("the actual API failure")

	var calls atomic.Int64

	page := func(ctx context.Context, offset, max int64) ([]int, int64, error) {
		// First call establishes a full page so a wave is launched.
		if calls.Add(1) == 1 {
			return seq(int(max)), 1_000_000, nil
		}

		// A later page is the genuine failure. Earlier-indexed siblings block
		// until the wave is cancelled, so without a preference the caller
		// would be told about the cancellation instead.
		if offset >= max*3 {
			return nil, 0, sentinel
		}

		<-ctx.Done()

		return nil, 0, ctx.Err()
	}

	_, err := paging.CollectWithPageSize(context.Background(), 10, page)

	if errors.Is(err, context.Canceled) {
		t.Fatal("reported a cancellation, which hides the failure that caused it")
	}

	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want the underlying failure", err)
	}
}

func TestCollectHonoursContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	var calls atomic.Int64

	page := func(_ context.Context, offset, max int64) ([]int, int64, error) {
		if calls.Add(1) >= 2 {
			cancel()
		}

		return seq(int(max)), 1_000_000, nil
	}

	if _, err := paging.CollectWithPageSize(
		ctx, 10, page,
	); !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestCollectRejectsNonPositivePageSize(t *testing.T) {
	t.Parallel()

	for _, size := range []int64{0, -1} {
		r := newRecorder(seq(5), 5)

		if _, err := paging.CollectWithPageSize(
			context.Background(), size, r.fn(),
		); err == nil {
			t.Errorf("page size %d: expected an error, got nil", size)
		}
	}
}

func TestCollectUsesDefaultPageSize(t *testing.T) {
	t.Parallel()

	var gotMax atomic.Int64

	page := func(_ context.Context, _, max int64) ([]int, int64, error) {
		gotMax.Store(max)

		return nil, 0, nil
	}

	if _, err := paging.Collect(context.Background(), page); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMax.Load() != paging.DefaultPageSize {
		t.Errorf("requested max %d, want DefaultPageSize %d",
			gotMax.Load(), paging.DefaultPageSize)
	}
}
