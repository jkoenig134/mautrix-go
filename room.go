package mautrix

import (
	"github.com/jkoenig134/mautrix-go/event"
	"github.com/jkoenig134/mautrix-go/id"
)

type RoomStateMap = map[event.Type]map[string]*event.Event

// Room represents a single Matrix room.
type Room struct {
	ID    id.RoomID
	State RoomStateMap
}

// PublicRoom represents the information about a public room obtainable from the room directory
type PublicRoom struct {
	CanonicalAlias   string   `json:"canonical_alias"`
	Name             string   `json:"name"`
	WorldReadable    bool     `json:"world_readable"`
	Topic            string   `json:"topic"`
	NumJoinedMembers int      `json:"num_joined_members"`
	AvatarURL        string   `json:"avatar_url"`
	RoomID           string   `json:"room_id"`
	GuestCanJoin     bool     `json:"guest_can_join"`
	Aliases          []string `json:"aliases"`
}

// UpdateState updates the room's current state with the given Event. This will clobber events based
// on the type/state_key combination.
func (room Room) UpdateState(evt *event.Event) {
	_, exists := room.State[evt.Type]
	if !exists {
		room.State[evt.Type] = make(map[string]*event.Event)
	}
	room.State[evt.Type][*evt.StateKey] = evt
}

// GetStateEvent returns the state event for the given type/state_key combo, or nil.
func (room Room) GetStateEvent(eventType event.Type, stateKey string) *event.Event {
	stateEventMap, _ := room.State[eventType]
	evt, _ := stateEventMap[stateKey]
	return evt
}

// GetMembershipState returns the membership state of the given user ID in this room. If there is
// no entry for this member, 'leave' is returned for consistency with left users.
func (room Room) GetMembershipState(userID id.UserID) event.Membership {
	state := event.MembershipLeave
	evt := room.GetStateEvent(event.StateMember, string(userID))
	if evt != nil {
		membership, ok := evt.Content.Raw["membership"].(string)
		if ok {
			state = event.Membership(membership)
		}
	}
	return state
}

// NewRoom creates a new Room with the given ID
func NewRoom(roomID id.RoomID) *Room {
	// Init the State map and return a pointer to the Room
	return &Room{
		ID:    roomID,
		State: make(RoomStateMap),
	}
}
