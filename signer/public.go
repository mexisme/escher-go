package signer

import (
	escher "github.com/adamluzsi/escher-go"
)

func (s *signer) SignRequest(request escher.Request, headersToSign []string) escher.Request {
	var authHeader = s.GenerateHeader(request, headersToSign)
	for _, header := range s.getDefaultHeaders(request.Headers) {
		request.Headers = append(request.Headers, header)
	}
	request.Headers = append(request.Headers, [2]string{s.config.AuthHeaderName, authHeader})
	return request
}

func (s *signer) CanonicalizeRequest(request escher.Request, headersToSign []string) string {
	var url = parsePathQuery(request.Url)
	var canonicalizedRequest = request.Method + "\n" +
		canonicalizePath(url.Path) + "\n" +
		canonicalizeQuery(url.Query) + "\n" +
		s.canonicalizeHeaders(request.Headers, headersToSign) + "\n" +
		s.canonicalizeHeadersToSign(headersToSign) + "\n" +
		s.computeDigest(request.Body)
	return canonicalizedRequest
}

func (s *signer) GenerateHeader(request escher.Request, headersToSign []string) string {
	return s.config.AlgoPrefix + "-HMAC-" + s.config.HashAlgo + " " +
		"Credential=" + s.generateCredentials() + ", " +
		"SignedHeaders=" + s.canonicalizeHeadersToSign(headersToSign) + ", " +
		"Signature=" + s.GenerateSignature(request, headersToSign)
}

func (s *signer) GenerateSignature(request escher.Request, headersToSign []string) string {
	var stringToSign = s.GetStringToSign(request, headersToSign)
	var signingKey = s.calculateSigningKey()
	return s.calculateSignature(stringToSign, signingKey)
}
