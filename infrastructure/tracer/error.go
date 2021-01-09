package tracer

// EnableRequestInterceptionError describe error that occurs during tracing url.
type EnableRequestInterceptionError struct {
	err error
}

// Error function return error message.
func (e *EnableRequestInterceptionError) Error() string {
	return "`EnableRequestInterception` failed. " + e.err.Error()
}

// NewTabError describe error that occurs during tracing url.
type NewTabError struct {
	err error
}

// Error function return error message.
func (e *NewTabError) Error() string {
	return "`newTab` failed. " + e.err.Error()
}

// CloseTabError describe error that occurs during tracing url.
type CloseTabError struct {
	err error
}

// Error function return error message.
func (e *CloseTabError) Error() string {
	return "`CloseTab` failed. " + e.err.Error()
}

// NetworkEventsError describe error that occurs during tracing url.
type NetworkEventsError struct {
	err error
}

// Error function return error message.
func (e *NetworkEventsError) Error() string {
	return "`NetworkEvents failed. " + e.err.Error()
}

// ActiveTabError describe error that occurs during tracing url.
type ActiveTabError struct {
	err error
}

// Error function return error message.
func (e *ActiveTabError) Error() string {
	return "`ActivateTab` failed. " + e.err.Error()
}

// AllEventsError describe error that occurs during tracing url.
type AllEventsError struct {
	err error
}

// Error function return error message.
func (e *AllEventsError) Error() string {
	return "`AllEvents` failed. %s" + e.err.Error()
}

// NavigateError describe error that occurs during tracing url.
type NavigateError struct {
	err error
}

// Error function return error message.
func (e *NavigateError) Error() string {
	return "`Navigate` failed. " + e.err.Error()
}

// InvalidFrameIDError describe error that occurs during tracing url.
type InvalidFrameIDError struct {
}

// Error function return error message.
func (e *InvalidFrameIDError) Error() string {
	return "invalid mainframe id"
}

// RedirectParseError describe error that occurs during tracing url.
type RedirectParseError struct {
	err error
}

// Error function return error message.
func (e *RedirectParseError) Error() string {
	return "an error during parsing redirects. " + e.err.Error()
}

// ResponseParseError describe error that occurs during tracing url.
type ResponseParseError struct {
	err error
}

// Error function return error message.
func (e *ResponseParseError) Error() string {
	return "an error during parsing response. " + e.err.Error()
}

// NoResponseError describe error that occurs during tracing url.
type NoResponseError struct {
}

// Error function return error message.
func (e *NoResponseError) Error() string {
	return "no responses found for mainframe"
}

// RedirectResponseNotExistsInRawDataError describe error that occurs during tracing url.
type RedirectResponseNotExistsInRawDataError struct {
}

// Error function return error message.
func (e *RedirectResponseNotExistsInRawDataError) Error() string {
	return "invalid redirect. `redirectResponse` param not exists"
}

// RequestNotExistsInRawDataError describe error that occurs during tracing url.
type RequestNotExistsInRawDataError struct {
}

// Error function return error message.
func (e *RequestNotExistsInRawDataError) Error() string {
	return "invalid redirect. `request` param not exists"
}

// InvalidToURLDataError describe error that occurs during tracing url.
type InvalidToURLDataError struct {
}

// Error function return error message.
func (e *InvalidToURLDataError) Error() string {
	return "invalid redirect `To` url"
}

// InvalidFromURLDataError describe error that occurs during tracing url.
type InvalidFromURLDataError struct {
}

// Error function return error message.
func (e *InvalidFromURLDataError) Error() string {
	return "invalid redirect `From` url"
}

// URLParamNotExistsInRedirectDataError describe error that occurs during tracing url.
type URLParamNotExistsInRedirectDataError struct {
}

// Error function return error message.
func (e *URLParamNotExistsInRedirectDataError) Error() string {
	return "invalid redirect. `redirectResponse` param `url` not exists"
}

// HeaderParamNotExistsInRedirectResponseDataError describe error that occurs during tracing url.
type HeaderParamNotExistsInRedirectResponseDataError struct {
}

// Error function return error message.
func (e *HeaderParamNotExistsInRedirectResponseDataError) Error() string {
	return "invalid redirect. `redirectResponse` param `headers` not exists"
}

// InitiatorParamNotExistsInRedirectDataError describe error that occurs during tracing url.
type InitiatorParamNotExistsInRedirectDataError struct {
}

// Error function return error message.
func (e *InitiatorParamNotExistsInRedirectDataError) Error() string {
	return "invalid redirect. `initiator` param not exists"
}

// ResponseParamNotExistsInRedirectDataError describe error that occurs during tracing url.
type ResponseParamNotExistsInRedirectDataError struct {
}

// Error function return error message.
func (e *ResponseParamNotExistsInRedirectDataError) Error() string {
	return "invalid redirect. `response` param not exists"
}

// HeaderParamNotExistsInRedirectDataError describe error that occurs during tracing url.
type HeaderParamNotExistsInRedirectDataError struct {
}

// Error function return error message.
func (e *HeaderParamNotExistsInRedirectDataError) Error() string {
	return "invalid redirect. request param `headers` not exists"
}

// SetScreenSizeError describe error that occurs during tracing url.
type SetScreenSizeError struct {
	err error
}

// Error function return error message.
func (e *SetScreenSizeError) Error() string {
	return "set screen size error: " + e.err.Error()
}

// SetVisibilitySizeError describe error that occurs during tracing url.
type SetVisibilitySizeError struct {
	err error
}

// Error function return error message.
func (e *SetVisibilitySizeError) Error() string {
	return "set visibility size error: " + e.err.Error()
}

// CaptureScreenshotError describe error that occurs during tracing url.
type CaptureScreenshotError struct {
	err error
}

// Error function return error message.
func (e *CaptureScreenshotError) Error() string {
	return "cannot capture screenshot: " + e.err.Error()
}
