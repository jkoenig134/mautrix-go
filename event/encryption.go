// Copyright (c) 2020 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package event

import (
	"encoding/json"

	"maunium.net/go/mautrix/id"
)

// Algorithm is a Matrix message encryption algorithm.
// https://matrix.org/docs/spec/client_server/r0.6.0#messaging-algorithm-names
type Algorithm string

const (
	AlgorithmOlmV1    Algorithm = "m.olm.v1.curve25519-aes-sha2"
	AlgorithmMegolmV1 Algorithm = "m.megolm.v1.aes-sha2"
)

// EncryptionEventContent represents the content of a m.room.encryption state event.
// https://matrix.org/docs/spec/client_server/r0.6.0#m-room-encryption
type EncryptionEventContent struct {
	// The encryption algorithm to be used to encrypt messages sent in this room. Must be 'm.megolm.v1.aes-sha2'.
	Algorithm Algorithm `json:"algorithm"`
	// How long the session should be used before changing it. 604800000 (a week) is the recommended default.
	RotationPeriodMillis int64 `json:"rotation_period_ms,omitempty"`
	// How many messages should be sent before changing the session. 100 is the recommended default.
	RotationPeriodMessages int `json:"rotation_period_messages,omitempty"`
}

// EncryptedEventContent represents the content of a m.room.encrypted message event.
// https://matrix.org/docs/spec/client_server/r0.6.0#m-room-encrypted
type EncryptedEventContent struct {
	Algorithm  Algorithm       `json:"algorithm"`
	SenderKey  string          `json:"sender_key"`
	DeviceID   id.DeviceID     `json:"device_id"`
	SessionID  string          `json:"session_id"`
	Ciphertext json.RawMessage `json:"ciphertext"`

	MegolmCiphertext string         `json:"-"`
	OlmCiphertext    OlmCiphertexts `json:"-"`
}

type OlmMessageType int

const (
	OlmPreKeyMessage OlmMessageType = 0
	OlmNormalMessage OlmMessageType = 1
)

type OlmCiphertexts map[string]struct {
	Body string         `json:"body"`
	Type OlmMessageType `json:"type"`
}

type serializableEncryptedEventContent EncryptedEventContent

func (content *EncryptedEventContent) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*serializableEncryptedEventContent)(content))
	if err != nil {
		return err
	}
	switch content.Algorithm {
	case AlgorithmOlmV1:
		content.OlmCiphertext = make(OlmCiphertexts)
		return json.Unmarshal(content.Ciphertext, &content.OlmCiphertext)
	case AlgorithmMegolmV1:
		return json.Unmarshal(content.Ciphertext, &content.MegolmCiphertext)
	default:
		return nil
	}
}

func (content *EncryptedEventContent) MarshalJSON() ([]byte, error) {
	var err error
	switch content.Algorithm {
	case AlgorithmOlmV1:
		content.Ciphertext, err = json.Marshal(content.OlmCiphertext)
	case AlgorithmMegolmV1:
		content.Ciphertext, err = json.Marshal(content.MegolmCiphertext)
	}
	if err != nil {
		return nil, err
	}
	return json.Marshal((*serializableEncryptedEventContent)(content))
}

// RoomKeyEventContent represents the content of a m.room_key to_device event.
// https://matrix.org/docs/spec/client_server/r0.6.0#m-room-key
type RoomKeyEventContent struct {
	Algorithm  Algorithm `json:"algorithm"`
	RoomID     id.RoomID `json:"room_id"`
	SessionID  string    `json:"session_id"`
	SessionKey string    `json:"session_key"`
}

// ForwardedRoomKeyEventContent represents the content of a m.forwarded_room_key to_device event.
// https://matrix.org/docs/spec/client_server/r0.6.0#m-forwarded-room-key
type ForwardedRoomKeyEventContent struct {
	RoomKeyEventContent
	SenderClaimedKey   string   `json:"sender_claimed_ed25519_key"`
	ForwardingKeyChain []string `json:"forwarding_curve25519_key_chain"`
}

type KeyRequestAction string

const (
	KeyRequestActionRequest = "request"
	KeyRequestActionCancel  = "request_cancellation"
)

// RoomKeyRequestEventContent represents the content of a m.room_key_request to_device event.
// https://matrix.org/docs/spec/client_server/r0.6.0#m-room-key-request
type RoomKeyRequestEventContent struct {
	Body               RequestedKeyInfo `json:"body"`
	Action             KeyRequestAction `json:"action"`
	RequestingDeviceID id.DeviceID      `json:"requesting_device_id"`
	RequestID          string           `json:"request_id"`
}

type RequestedKeyInfo struct {
	Algorithm Algorithm `json:"algorithm"`
	RoomID    id.RoomID `json:"room_id"`
	SenderKey string    `json:"sender_key"`
	SessionID string    `json:"session_id"`
}