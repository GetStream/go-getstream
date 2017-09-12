package getstream

// ScopeAction defines the Actions allowed by a scope token
type ScopeAction uint32

const (
	// ScopeActionRead : GET, OPTIONS, HEAD
	ScopeActionRead ScopeAction = 1
	// ScopeActionWrite : POST, PUT, PATCH
	ScopeActionWrite ScopeAction = 2
	// ScopeActionDelete : DELETE
	ScopeActionDelete ScopeAction = 4
	// ScopeActionAll : The JWT has permission to all HTTP verbs
	ScopeActionAll ScopeAction = 8
)

// Value returns a string representation
func (a ScopeAction) Value() string {
	switch a {
	case ScopeActionRead:
		return "read"
	case ScopeActionWrite:
		return "write"
	case ScopeActionDelete:
		return "delete"
	case ScopeActionAll:
		return "*"
	default:
		return ""
	}
}

// ScopeContext defines the resources accessible by a scope token
type ScopeContext uint32

const (
	ScopeContextNoContext ScopeContext = 0
	// ScopeContextActivities :  Activities Endpoint
	ScopeContextActivities ScopeContext = 1
	// ScopeContextFeed : Feed Endpoint
	ScopeContextFeed ScopeContext = 2
	// ScopeContextFollower : Following + Followers Endpoint
	ScopeContextFollower ScopeContext = 4
	// ScopeContextFeedTargets : UpdateFeedToTargets
	ScopeContextFeedTargets ScopeContext = 8
	// ScopeContextAll : Allow access to any resource
	ScopeContextAll ScopeContext = 16
)

// Value returns a string representation
func (a ScopeContext) Value() string {
	switch a {
	case ScopeContextActivities:
		return "activities"
	case ScopeContextFeed:
		return "feed"
	case ScopeContextFollower:
		return "follower"
	case ScopeContextFeedTargets:
		return "feed_targets"
	case ScopeContextAll:
		return "*"
	default:
		return ""
	}
}
