package uulid

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUULIDIsEmpty(t *testing.T) {
	r := require.New(t)
	var id ULID
	r.True(id.IsEmpty())

	id = NewULID()
	r.False(id.IsEmpty())
}

func TestULID_Value(t *testing.T) {
	id := NewULID()
	fmt.Println(id.String())
	v, _ := id.Value()
	require.Equal(t, id.String(), v)
}

func TestUULID_ThreeNewIDsAreDifferent(t *testing.T) {
	r := require.New(t)

	goroutines := 100
	ids_per_goroutine := 1000

	ids := make(map[ULID]bool, ids_per_goroutine*goroutines)
	m := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			var oids []ULID
			for j := 0; j < ids_per_goroutine; j++ {
				oids = append(oids, NewULID())
			}

			m.Lock()
			for _, id := range oids {
				ids[id] = true
			}
			m.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	r.Equal(goroutines*ids_per_goroutine, len(ids))
}

func TestULID_ScanValue(t *testing.T) {
	r := require.New(t)

	expected := NewULID()
	v, err := expected.Value()
	r.NoError(err)

	var id ULID
	r.NoError(id.Scan(v))

	r.Equal(expected, id)
	r.Equal(expected.String(), id.String())

	r.NoError(id.Scan([]byte("015af13d-2271-fb69-2dcd-fb24a1fd7dcc")))
}

func BenchmarkNewULID(b *testing.B) {
	//num:=10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewULID()
	}
}
