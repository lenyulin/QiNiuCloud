package textshrink

import "errors"

var (
	ErrShrinkerNotAvailable            = errors.New("Shrinker Creat Http Request Failed")
	ErrShrinkerMarshalData             = errors.New("Shrinker Marshal Data Failed")
	ErrShrinkerSendHttpRequestFailed   = errors.New("Shrinker Send Http Request Failed")
	ErrShrinkerNotGetDataFromRemoteAPI = errors.New("Shrinker Not Get Data From Remote API")
)
