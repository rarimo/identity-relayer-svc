/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type FeeEstimateRequest struct {
	Key
	Attributes    map[string]interface{}          `json:"attributes"`
	Relationships FeeEstimateRequestRelationships `json:"relationships"`
}
type FeeEstimateRequestResponse struct {
	Data     FeeEstimateRequest `json:"data"`
	Included Included           `json:"included"`
}

type FeeEstimateRequestListResponse struct {
	Data     []FeeEstimateRequest `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustFeeEstimateRequest - returns FeeEstimateRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustFeeEstimateRequest(key Key) *FeeEstimateRequest {
	var feeEstimateRequest FeeEstimateRequest
	if c.tryFindEntry(key, &feeEstimateRequest) {
		return &feeEstimateRequest
	}
	return nil
}
