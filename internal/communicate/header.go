package communicate

import "net/http"

func makeGetEndpointRequestHeader(signature string) http.Header {
	header := make(http.Header)
	header.Set("Accept-Language", "zh-Hans")
	header.Set("X-ClientVersion", clientVersion)
	header.Set("X-UserId", userId)
	header.Set("X-HomeGeographicRegion", homeGeographicRegion)
	header.Set("X-ClientTraceId", clientTraceId)
	header.Set("X-MT-Signature", signature)
	header.Set("User-Agent", userAgent)
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Content-Length", "0")
	header.Set("Accept-Encoding", "gzip")
	return header
}
