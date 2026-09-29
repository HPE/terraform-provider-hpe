// Package paging walks an offset-paginated Morpheus list endpoint to
// completion.
//
// Most plural data sources in this provider issue a single request with a large
// Max and silently truncate anything beyond it. That is tolerable when the
// result is handed straight to the practitioner, but not when the provider
// filters the response itself: a client-side match applied to a truncated fetch
// quietly reports "no such item" for something that exists. Collecting every
// page first makes such filtering sound.
//
// # What is assumed of the endpoint
//
// Completeness rests on one assumption: that meta.total is never *lower* than
// the number of records the endpoint will actually serve. The walk ends when
// the offset reaches total, so an endpoint that under-reports would be
// truncated — reopening exactly the unsoundness this package exists to prevent.
//
// Over-reporting is safe, and so is omitting total altogether: both are caught
// by the empty-page guard, which ends the walk whatever total claims. Only
// under-reporting is dangerous, and no client-side check can detect it — an
// endpoint saying "5" and serving 25 is indistinguishable from one saying "5"
// and serving 5, until the records that were never asked for are missed.
//
// Morpheus reports a total consistent with what it pages, so this holds today.
// It is recorded because it is the package's one unverifiable premise.
package paging

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// DefaultPageSize is the number of records requested per page.
//
// Larger pages stop helping quickly. Walking 233 records on one appliance took
// 2888ms in pages of 25, 2054ms in pages of 50, and then 1873ms, 1875ms and
// 1908ms in pages of 100, 250 and 1000 — three requests and one request being
// indistinguishable. A second appliance holding 7040 records showed the same
// floor of roughly 8-10ms per record, reached by a page size of about 100.
//
// The cost is therefore per record on the server, not per request, and no page
// size avoids it. What does help is fetching pages at the same time: the same
// five pages took 2196ms in sequence and 703ms in parallel.
const DefaultPageSize int64 = 1000

// DefaultConcurrency is how many pages are fetched at once.
//
// Four is a deliberate compromise. The measurement above used five and gained
// threefold; more would gain a little more and cost the appliance a little
// more, and a data source is read on every refresh, plan and apply, so it is
// not the only thing asking.
//
// It also bounds how much work a wrong total can waste: pages are fetched a
// wave at a time, so an endpoint claiming far more records than it holds costs
// at most one wave of empty responses rather than a request per phantom page.
const DefaultConcurrency = 4

// PageFunc fetches one page beginning at offset and returns the items in that
// page together with the total number of records matching the query.
//
// total comes from the response envelope's meta.total. A total of 0, or one
// that disagrees with what arrives, is tolerated: collection then relies on
// reaching a short page instead.
type PageFunc[T any] func(ctx context.Context, offset, max int64) (items []T, total int64, err error)

// Collect walks page until every record has been read and returns them in the
// order the endpoint served them.
//
// Pages after the first are fetched a wave at a time, DefaultConcurrency at
// once. A wave stops the walk when any page in it comes back short, since a
// page smaller than requested is the end of the data.
//
// The walk is deliberately conservative about what it assumes. Offsets within a
// wave are computed from the page size, which supposes every earlier page was
// full; that supposition is checked rather than trusted, and a wave containing
// a short page anywhere but at its end is discarded and re-walked in sequence.
// The cost of being wrong is a wasted fetch, never a wrong answer.
//
// Note that offset paging is not atomic. A record created or deleted mid-walk
// can be seen twice or missed, and ordering by a non-unique key (Morpheus
// orders virtual images by name, which is not unique) allows rows to move
// between pages. Callers that need exactness should de-duplicate on a stable
// identifier.
func Collect[T any](ctx context.Context, page PageFunc[T]) ([]T, error) {
	return CollectWithPageSize(ctx, DefaultPageSize, page)
}

// CollectWithPageSize is Collect with an explicit page size, for callers with a
// reason to deviate from DefaultPageSize.
func CollectWithPageSize[T any](
	ctx context.Context,
	pageSize int64,
	page PageFunc[T],
) ([]T, error) {
	if pageSize <= 0 {
		return nil, fmt.Errorf("paging: page size must be positive, got %d", pageSize)
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("paging: %w", err)
	}

	first, total, err := page(ctx, 0, pageSize)
	if err != nil {
		return nil, err
	}

	result := first
	offset := int64(len(first))

	// An empty page is the end of the data, whatever total says.
	//
	// A short page is deliberately not treated as the end. An endpoint may
	// return fewer records than asked for — through a per-request cap or
	// permission filtering — without having run out, and stopping there would
	// silently drop everything it had not yet sent.
	if offset == 0 {
		return result, nil
	}

	// The server has already accounted for everything it claims to hold.
	if total > 0 && offset >= total {
		return result, nil
	}

	for {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("paging: %w", err)
		}

		items, done, err := fetchWave(ctx, offset, pageSize, page)
		if err != nil {
			return nil, err
		}

		result = append(result, items...)
		offset += int64(len(items))

		if done {
			return result, nil
		}

		if total > 0 && offset >= total {
			return result, nil
		}
	}
}

// fetchWave fetches up to DefaultConcurrency pages at once, beginning at
// offset, and returns them in offset order.
//
// done reports that the walk has reached the end, which is true only when a
// page comes back *empty*. A short page does not end the walk: an endpoint may
// return fewer records than asked for through a per-request cap or permission
// filtering without having run out, and stopping there would silently drop
// whatever it had not yet sent.
//
// A short page does end the *wave*, though. The pages after it were fetched at
// offsets computed on the assumption that it was full, so they may have skipped
// records and are discarded. The next wave resumes from what was actually read
// and settles the question by either returning more or coming back empty.
func fetchWave[T any](
	ctx context.Context,
	offset, pageSize int64,
	page PageFunc[T],
) ([]T, bool, error) {
	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		pages = make([][]T, DefaultConcurrency)
		errs  = make([]error, DefaultConcurrency)
	)

	for i := range DefaultConcurrency {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			items, _, err := page(wctx, offset+int64(i)*pageSize, pageSize)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs[i] = err
				cancel()

				return
			}

			pages[i] = items
		}(i)
	}

	wg.Wait()

	// Prefer a genuine failure over a cancellation. When one page fails the
	// rest of the wave is cancelled, so a lower-indexed sibling often carries
	// context.Canceled — reporting that instead of the API error it was caused
	// by would hide the reason the walk stopped.
	var cancelled error

	for _, err := range errs {
		if err == nil {
			continue
		}

		if errors.Is(err, context.Canceled) {
			if cancelled == nil {
				cancelled = err
			}

			continue
		}

		return nil, false, err
	}

	if cancelled != nil && ctx.Err() == nil {
		// Every failure was a cancellation this wave caused itself, with the
		// caller's context still live. Surface it rather than reporting success.
		return nil, false, cancelled
	}

	if err := ctx.Err(); err != nil {
		return nil, false, fmt.Errorf("paging: %w", err)
	}

	var out []T

	for _, items := range pages {
		out = append(out, items...)

		// An empty page is the end of the data, whatever total claims.
		if len(items) == 0 {
			return out, true, nil
		}

		// A short page is either the end or an endpoint returning fewer records
		// than asked for. The two are indistinguishable here, and the pages
		// after it in this wave were fetched at offsets that assumed it was
		// full, so they are discarded. The next wave resumes from what was
		// actually read, and settles the question by either returning more or
		// coming back empty.
		if int64(len(items)) < pageSize {
			return out, false, nil
		}
	}

	return out, false, nil
}
