/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type RelayTask struct {
	Key
	Attributes    map[string]interface{} `json:"attributes"`
	Relationships RelayTaskRelationships `json:"relationships"`
}
type RelayTaskResponse struct {
	Data     RelayTask `json:"data"`
	Included Included  `json:"included"`
}

type RelayTaskListResponse struct {
	Data     []RelayTask `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustRelayTask - returns RelayTask from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRelayTask(key Key) *RelayTask {
	var relayTask RelayTask
	if c.tryFindEntry(key, &relayTask) {
		return &relayTask
	}
	return nil
}
