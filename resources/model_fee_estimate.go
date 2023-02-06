/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type FeeEstimate struct {
	Key
	Attributes    FeeEstimateAttributes    `json:"attributes"`
	Relationships FeeEstimateRelationships `json:"relationships"`
}
type FeeEstimateResponse struct {
	Data     FeeEstimate `json:"data"`
	Included Included    `json:"included"`
}

type FeeEstimateListResponse struct {
	Data     []FeeEstimate `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustFeeEstimate - returns FeeEstimate from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustFeeEstimate(key Key) *FeeEstimate {
	var feeEstimate FeeEstimate
	if c.tryFindEntry(key, &feeEstimate) {
		return &feeEstimate
	}
	return nil
}
