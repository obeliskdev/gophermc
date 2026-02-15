package protocol

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Definition struct {
	ProtocolVersion int32
	PacketNames     map[State]map[Direction]map[int32]string
	PacketIDs       map[State]map[Direction]map[string]int32
}

var protocolRegistry = make(map[Version]*Definition)

func GetDefinition(v Version) *Definition {
	for current := v; current >= 0; current-- {
		if def, ok := protocolRegistry[current]; ok {
			if current != v {
				protocolRegistry[v] = def
			}
			return def
		}
	}
	return nil
}

func NewPacket(v Version, s State, d Direction, id int32) (Packet, error) {
	def := GetDefinition(v)
	if def == nil {
		return nil, fmt.Errorf("unknown protocol version %s", v)
	}

	name, ok := def.PacketNames[s][d][id]
	if !ok {
		return nil, ErrUnknownPacket
	}

	factory, ok := packetConstructors[name]
	if !ok {
		return nil, fmt.Errorf("no constructor for packet type '%s'", name)
	}
	return factory(), nil
}

func GetPacketID(v Version, s State, d Direction, p Packet) (int32, bool) {
	def := GetDefinition(v)
	if def == nil {
		return -1, false
	}

	packetType := reflect.TypeOf(p)
	name, ok := packetTypes[packetType]
	if !ok {
		return -1, false
	}

	id, ok := def.PacketIDs[s][d][name]
	return id, ok
}

func VersionFromString(s string) (Version, bool) {
	version, ok := stringToVersion[s]
	if ok {
		return version, true
	}

	major, minor, patch, hasPatch, ok := parseReleaseVersion(s)
	if !ok {
		return 0, false
	}

	bestPatch := -1
	var bestVersion Version
	for knownStr, knownVersion := range stringToVersion {
		kMajor, kMinor, kPatch, _, kOK := parseReleaseVersion(knownStr)
		if !kOK || kMajor != major || kMinor != minor {
			continue
		}

		if hasPatch && kPatch > patch {
			continue
		}
		if kPatch > bestPatch {
			bestPatch = kPatch
			bestVersion = knownVersion
		}
	}
	if bestPatch >= 0 {
		return bestVersion, true
	}

	return 0, false
}

func parseReleaseVersion(s string) (major int, minor int, patch int, hasPatch bool, ok bool) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, 0, 0, false, false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, false, false
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, false, false
	}

	if len(parts) == 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, 0, 0, false, false
		}
		return major, minor, patch, true, true
	}

	return major, minor, 0, false, true
}
