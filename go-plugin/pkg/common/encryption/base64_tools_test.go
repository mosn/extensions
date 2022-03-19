package encryption

import (
	"reflect"
	"testing"
)

func TestBase64Decoder(t *testing.T) {

	tests := []struct {
		msg  string
		want []int8
	}{
		{
			msg:  "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpWV11YGwcFCAEcAwQFBxkIARwDBRcIPQRiS0B8UFdTBjsOYVFBdEJLWHZSQFAIBQgDAAMFBAcLF2NXR3ZARV58UEZWCj8KZV1FYUpHcVdDXQ8AAwYHBgYJAA4cZlBCZEFCdlJAUAg9BGNXR2dMRWNRXFcNBQAGBgwHAQ8bZ1NDa0hBZ11YUwkyDWBWQGZTRnZeDHJ6fGMHCAACAwQFBgcJBwcFDQUKGGpURmBRRHhYBjsOYVFBCD0EY1dHd1pSUgZzBwMFCRllXUVxXFBQCD0EY1dHeUZRCQ8BAlBbWBhETV9LUkZRGFJAUldDQFxZWRZiR11xdntyQFJXQ0BcWVkCEdW5gtO2tt6Rtdubs9mLonQCAwcZ0Li32Y2D24ms0bGI1J+V0YmW3rCd14y/0ouY1JaC3IGTG920hdeJptOxvdScituJrN+XlNS6jdOalt2GgNeMv9KLmNScv9K9ptiEvdqcg9C9iNG2v9aQuNChiNSQrQgaZFJMfEFUCj8KGGpURg0+CRlkQUJ6VlVRCD0EcEJDfFBXUxcPOA8bcVlUTVxXXUALPA==",
			want: []int8{13, 13, 75, 89, 89, 22, 65, 93, 67, 65, 90, 91, 91, 11, 21, 9, 31, 2, 17, 20, 80, 88, 84, 87, 85, 91, 93, 83, 8, 20, 98, 108, 119, 31, 11, 22, 10, 8, 61, 4, 117, 93, 80, 65, 88, 83, 89, 76, 17, 74, 94, 88, 91, 69, 10, 26, 86, 87, 93, 88, 27, 7, 5, 8, 1, 28, 3, 4, 5, 7, 25, 8, 1, 28, 3, 5, 23, 8, 61, 4, 98, 75, 64, 124, 80, 87, 83, 6, 59, 14, 97, 81, 65, 116, 66, 75, 88, 118, 82, 64, 80, 8, 5, 8, 3, 0, 3, 5, 4, 7, 11, 23, 99, 87, 71, 118, 64, 69, 94, 124, 80, 70, 86, 10, 63, 10, 101, 93, 69, 97, 74, 71, 113, 87, 67, 93, 15, 0, 3, 6, 7, 6, 6, 9, 0, 14, 28, 102, 80, 66, 100, 65, 66, 118, 82, 64, 80, 8, 61, 4, 99, 87, 71, 103, 76, 69, 99, 81, 92, 87, 13, 5, 0, 6, 6, 12, 7, 1, 15, 27, 103, 83, 67, 107, 72, 65, 103, 93, 88, 83, 9, 50, 13, 96, 86, 64, 102, 83, 70, 118, 94, 12, 114, 122, 124, 99, 7, 8, 0, 2, 3, 4, 5, 6, 7, 9, 7, 7, 5, 13, 5, 10, 24, 106, 84, 70, 96, 81, 68, 120, 88, 6, 59, 14, 97, 81, 65, 8, 61, 4, 99, 87, 71, 119, 90, 82, 82, 6, 115, 7, 3, 5, 9, 25, 101, 93, 69, 113, 92, 80, 80, 8, 61, 4, 99, 87, 71, 121, 70, 81, 9, 15, 1, 2, 80, 91, 88, 24, 68, 77, 95, 75, 82, 70, 81, 24, 82, 64, 82, 87, 67, 64, 92, 89, 89, 22, 98, 71, 93, 113, 118, 123, 114, 64, 82, 87, 67, 64, 92, 89, 89, 2, 17, -43, -71, -126, -45, -74, -74, -34, -111, -75, -37, -101, -77, -39, -117, -94, 116, 2, 3, 7, 25, -48, -72, -73, -39, -115, -125, -37, -119, -84, -47, -79, -120, -44, -97, -107, -47, -119, -106, -34, -80, -99, -41, -116, -65, -46, -117, -104, -44, -106, -126, -36, -127, -109, 27, -35, -76, -123, -41, -119, -90, -45, -79, -67, -44, -100, -118, -37, -119, -84, -33, -105, -108, -44, -70, -115, -45, -102, -106, -35, -122, -128, -41, -116, -65, -46, -117, -104, -44, -100, -65, -46, -67, -90, -40, -124, -67, -38, -100, -125, -48, -67, -120, -47, -74, -65, -42, -112, -72, -48, -95, -120, -44, -112, -83, 8, 26, 100, 82, 76, 124, 65, 84, 10, 63, 10, 24, 106, 84, 70, 13, 62, 9, 25, 100, 65, 66, 122, 86, 85, 81, 8, 61, 4, 112, 66, 67, 124, 80, 87, 83, 23, 15, 56, 15, 27, 113, 89, 84, 77, 92, 87, 93, 64, 11, 60},
		},
	}

	for _, tt := range tests {
		got := Base64Decoder([]byte(tt.msg))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Base64Decoder() got = %v, want %v", got, tt.want)
		}
	}
}

func TestBase64Encoder(t *testing.T) {
	tests := []struct {
		want string
		msg  []int8
	}{
		{
			want: "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpWV11YGwcFCAEcAwQFBxkIARwDBRcIPQRiS0B8UFdTBjsOYVFBdEJLWHZSQFAIBQgDAAMFBAcLF2NXR3ZARV58UEZWCj8KZV1FYUpHcVdDXQ8AAwYHBgYJAA4cZlBCZEFCdlJAUAg9BGNXR2dMRWNRXFcNBQAGBgwHAQ8bZ1NDa0hBZ11YUwkyDWBWQGZTRnZeDHJ6fGMHCAACAwQFBgcJBwcFDQUKGGpURmBRRHhYBjsOYVFBCD0EY1dHd1pSUgZzBwMFCRllXUVxXFBQCD0EY1dHeUZRCQ8BAlBbWBhETV9LUkZRGFJAUldDQFxZWRZiR11xdntyQFJXQ0BcWVkCEdW5gtO2tt6Rtdubs9mLonQCAwcZ0Li32Y2D24ms0bGI1J+V0YmW3rCd14y/0ouY1JaC3IGTG920hdeJptOxvdScituJrN+XlNS6jdOalt2GgNeMv9KLmNScv9K9ptiEvdqcg9C9iNG2v9aQuNChiNSQrQgaZFJMfEFUCj8KGGpURg0+CRlkQUJ6VlVRCD0EcEJDfFBXUxcPOA8bcVlUTVxXXUALPA==",
			msg:  []int8{13, 13, 75, 89, 89, 22, 65, 93, 67, 65, 90, 91, 91, 11, 21, 9, 31, 2, 17, 20, 80, 88, 84, 87, 85, 91, 93, 83, 8, 20, 98, 108, 119, 31, 11, 22, 10, 8, 61, 4, 117, 93, 80, 65, 88, 83, 89, 76, 17, 74, 94, 88, 91, 69, 10, 26, 86, 87, 93, 88, 27, 7, 5, 8, 1, 28, 3, 4, 5, 7, 25, 8, 1, 28, 3, 5, 23, 8, 61, 4, 98, 75, 64, 124, 80, 87, 83, 6, 59, 14, 97, 81, 65, 116, 66, 75, 88, 118, 82, 64, 80, 8, 5, 8, 3, 0, 3, 5, 4, 7, 11, 23, 99, 87, 71, 118, 64, 69, 94, 124, 80, 70, 86, 10, 63, 10, 101, 93, 69, 97, 74, 71, 113, 87, 67, 93, 15, 0, 3, 6, 7, 6, 6, 9, 0, 14, 28, 102, 80, 66, 100, 65, 66, 118, 82, 64, 80, 8, 61, 4, 99, 87, 71, 103, 76, 69, 99, 81, 92, 87, 13, 5, 0, 6, 6, 12, 7, 1, 15, 27, 103, 83, 67, 107, 72, 65, 103, 93, 88, 83, 9, 50, 13, 96, 86, 64, 102, 83, 70, 118, 94, 12, 114, 122, 124, 99, 7, 8, 0, 2, 3, 4, 5, 6, 7, 9, 7, 7, 5, 13, 5, 10, 24, 106, 84, 70, 96, 81, 68, 120, 88, 6, 59, 14, 97, 81, 65, 8, 61, 4, 99, 87, 71, 119, 90, 82, 82, 6, 115, 7, 3, 5, 9, 25, 101, 93, 69, 113, 92, 80, 80, 8, 61, 4, 99, 87, 71, 121, 70, 81, 9, 15, 1, 2, 80, 91, 88, 24, 68, 77, 95, 75, 82, 70, 81, 24, 82, 64, 82, 87, 67, 64, 92, 89, 89, 22, 98, 71, 93, 113, 118, 123, 114, 64, 82, 87, 67, 64, 92, 89, 89, 2, 17, -43, -71, -126, -45, -74, -74, -34, -111, -75, -37, -101, -77, -39, -117, -94, 116, 2, 3, 7, 25, -48, -72, -73, -39, -115, -125, -37, -119, -84, -47, -79, -120, -44, -97, -107, -47, -119, -106, -34, -80, -99, -41, -116, -65, -46, -117, -104, -44, -106, -126, -36, -127, -109, 27, -35, -76, -123, -41, -119, -90, -45, -79, -67, -44, -100, -118, -37, -119, -84, -33, -105, -108, -44, -70, -115, -45, -102, -106, -35, -122, -128, -41, -116, -65, -46, -117, -104, -44, -100, -65, -46, -67, -90, -40, -124, -67, -38, -100, -125, -48, -67, -120, -47, -74, -65, -42, -112, -72, -48, -95, -120, -44, -112, -83, 8, 26, 100, 82, 76, 124, 65, 84, 10, 63, 10, 24, 106, 84, 70, 13, 62, 9, 25, 100, 65, 66, 122, 86, 85, 81, 8, 61, 4, 112, 66, 67, 124, 80, 87, 83, 23, 15, 56, 15, 27, 113, 89, 84, 77, 92, 87, 93, 64, 11, 60},
		},
	}

	for _, tt := range tests {
		got := string(Base64Encoder(tt.msg))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Base64Encoder() got = %s, want %v", got, tt.want)
		}
	}
}
