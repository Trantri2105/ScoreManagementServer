package request

type GetPredictedGpaRequest struct {
	GpaT1 float64 `json:"gpa_t1"`
	GpaT2 float64 `json:"gpa_t2"`
	GpaT3 float64 `json:"gpa_t3"`
	GpaT4 float64 `json:"gpa_t4"`
}
