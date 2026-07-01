// Package kobo provides a zero-dependency Go client for the Kobo API v1.
// This file provides helpers for accessing the SDK.
package kobo

// ptr returns a pointer to the provided value. Useful for optional fields.
//
//	req := kobo.UpdateIdentityRequest{
//	    DisplayName: kobo.Ptr("Jane Doe"),
//	}
func Ptr[T any](v T) *T { return &v }
