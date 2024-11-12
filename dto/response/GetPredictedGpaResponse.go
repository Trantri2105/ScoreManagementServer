package response

type GetPredictedGpaResponse struct {
	PredictedFutureGpa float64 `json:"predicted_future_gpa"`
	Error              string  `json:"error,omitempty"`
}
