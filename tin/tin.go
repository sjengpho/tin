package tin

// Comparable is the interface implemented by an object that can
// compare itself with another instance of the same type.
//
// Equal does a equality check.
type Comparable interface {
	Equal(t interface{}) bool
}
