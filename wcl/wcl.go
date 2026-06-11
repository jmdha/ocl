package wcl

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Kind int

const (
	KindVersion Kind = iota
	KindChallengeModeStart
	KindChallengeModeEnd
	KindEncounterStart
	KindEncounterEnd
	KindCombatantInfo
	KindMapChange
	KindZoneChange
	KindUnitDied
	KindPartyKill
	KindSpellCastSuccess
	KindSpellSummon
	KindSpellInstakill
	KindSpellEnergize
	KindSpellHeal
	KindSpellDamage
	KindSpellPeriodicDamage
	KindSpellPeriodicHeal
	KindSpellAuraApplied
	KindSpellAuraRefresh
	KindSpellAuraRemoved
	KindUnknown
)

type Event struct {
	Kind             Kind
	Version          EventVersion
	EncounterStart   EventEncounterStart
	EncounterEnd     EventEncounterEnd
	MapChange        EventMapChange
	ZoneChange       EventZoneChange
	UnitDied         EventUnitDied
	SpellCastSuccess EventSpellCastSuccess
	SpellDamage      EventSpellDamage
	SpellHeal        EventSpellHeal
	SpellAura        EventSpellAura
	Unknown          EventUnknown
}

type EventVersion struct {
	Timestamp time.Time
	Log       int64
	Advanced  int64
	Major     int64
	Minor     int64
	Patch     int64
	Project   int64
}

type EventEncounterStart struct {
	Timestamp     time.Time
	EncounterID   int64
	EncounterName string
	DifficultyID  int64
	GroupSize     int64
	InstanceID    int64
}

type EventEncounterEnd struct {
	Timestamp     time.Time
	EncounterID   int64
	EncounterName string
	DifficultyID  int64
	GroupSize     int64
	Success       bool
	Duration      int64
}

type EventMapChange struct {
	Timestamp  time.Time
	InstanceID int64
	Name       string
}

type EventZoneChange struct {
	Timestamp    time.Time
	InstanceID   int64
	Name         string
	DifficultyID int64
}

type EventUnitDied struct {
	Timestamp  time.Time
	TargetGUID string
	TargetName string
}

type EventSpellCastSuccess struct {
	Timestamp  time.Time
	SourceGUID string
	SourceName string
	TargetGUID string
	TargetName string
	SpellID    int64
	SpellName  string
}

type EventSpellDamage struct {
	Timestamp  time.Time
	SourceGUID string
	SourceName string
	TargetGUID string
	TargetName string
	SpellID    int64
	SpellName  string
	Damage     int64
	DamageRaw  int64
}

type EventSpellHeal struct {
	Timestamp  time.Time
	SourceGUID string
	SourceName string
	TargetGUID string
	TargetName string
	SpellID    int64
	SpellName  string
	Amount     int64
}

type EventSpellAura struct {
	Timestamp  time.Time
	SourceGUID string
	SourceName string
	TargetGUID string
	TargetName string
	SpellID    int64
	SpellName  string
}

type EventUnknown struct {
	Timestamp time.Time
	Kind      string
}

func Parse(event *Event, line string) error {
	sTimestamp, sFields, ok := strings.Cut(line, "  ")
	if !ok {
		return fmt.Errorf("missing event seperator \"  \"")
	}

	timestamp, err := ParseTimestamp(sTimestamp)
	if err != nil {
		return err
	}

	r := csv.NewReader(strings.NewReader(sFields))
	fields, err := r.Read()
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return fmt.Errorf("invalid event: no fields")
	}

	event.Kind = MatchKind(fields[0])
	switch event.Kind {
	case KindVersion:
		return parseVersion(event, timestamp, fields)
	case KindEncounterStart:
		return parseEncounterStart(event, timestamp, fields)
	case KindEncounterEnd:
		return parseEncounterEnd(event, timestamp, fields)
	case KindMapChange:
		return parseMapChange(event, timestamp, fields)
	case KindZoneChange:
		return parseZoneChange(event, timestamp, fields)
	case KindUnitDied:
		return parseUnitDied(event, timestamp, fields)
	case KindSpellCastSuccess:
		return parseSpellCastSuccess(event, timestamp, fields)
	case KindSpellDamage, KindSpellPeriodicDamage:
		return parseSpellDamage(event, timestamp, fields)
	case KindSpellHeal, KindSpellPeriodicHeal:
		return parseSpellHeal(event, timestamp, fields)
	case KindSpellAuraApplied, KindSpellAuraRefresh, KindSpellAuraRemoved:
		return parseSpellAura(event, timestamp, fields)
	default:
		event.Unknown.Timestamp = timestamp
		event.Unknown.Kind = fields[0]
		return nil
	}
}

func ParseTimestamp(text string) (time.Time, error) {
	const layout = "1/2/2006 15:04:05.0000"
	t, err := time.Parse(layout, text)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func MatchKind(text string) Kind {
	switch text {
	case "COMBAT_LOG_VERSION":
		return KindVersion
	case "CHALLENGE_MODE_START":
		return KindChallengeModeStart
	case "CHALLENGE_MODE_END":
		return KindChallengeModeEnd
	case "ENCOUNTER_START":
		return KindEncounterStart
	case "ENCOUNTER_END":
		return KindEncounterEnd
	case "MAP_CHANGE":
		return KindMapChange
	case "ZONE_CHANGE":
		return KindZoneChange
	case "UNIT_DIED":
		return KindUnitDied
	case "PARTY_KILL":
		return KindPartyKill
	case "SPELL_CAST_SUCCESS":
		return KindSpellCastSuccess
	case "SPELL_SUMMON":
		return KindSpellSummon
	case "SPELL_INSTAKILL":
		return KindSpellInstakill
	case "SPELL_ENERGIZE":
		return KindSpellEnergize
	case "SPELL_HEAL":
		return KindSpellHeal
	case "SPELL_DAMAGE":
		return KindSpellDamage
	case "SPELL_PERIODIC_DAMAGE":
		return KindSpellPeriodicDamage
	case "SPELL_PERIODIC_HEAL":
		return KindSpellPeriodicHeal
	case "SPELL_AURA_APPLIED":
		return KindSpellAuraApplied
	case "SPELL_AURA_REFRESH":
		return KindSpellAuraRefresh
	case "SPELL_AURA_REMOVED":
		return KindSpellAuraRemoved
	default:
		return KindUnknown
	}
}

// need returns an error if fields doesn't have enough elements.
func need(fields []string, n int, event string) error {
	if len(fields) < n {
		return fmt.Errorf("%s: need %d fields, got %d", event, n, len(fields))
	}
	return nil
}

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseBool(s string) (bool, error) {
	n, err := parseInt(s)
	return n != 0, err
}

// COMBAT_LOG_VERSION layout:
// [0] event, [1] log, [2] advanced_key, [3] advanced, [4] build_key, [5] major.minor.patch, [6] project_key, [7] project
func parseVersion(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 8, "COMBAT_LOG_VERSION"); err != nil {
		return err
	}
	log, err := parseInt(fields[1])
	if err != nil {
		return err
	}
	advanced, err := parseInt(fields[3])
	if err != nil {
		return err
	}
	parts := strings.SplitN(fields[5], ".", 3)
	if len(parts) != 3 {
		return fmt.Errorf("COMBAT_LOG_VERSION: invalid build version %q", fields[5])
	}
	major, err := parseInt(parts[0])
	if err != nil {
		return err
	}
	minor, err := parseInt(parts[1])
	if err != nil {
		return err
	}
	patch, err := parseInt(parts[2])
	if err != nil {
		return err
	}
	project, err := parseInt(fields[7])
	if err != nil {
		return err
	}
	event.Version = EventVersion{
		Timestamp: ts,
		Log:       log,
		Advanced:  advanced,
		Major:     major,
		Minor:     minor,
		Patch:     patch,
		Project:   project,
	}
	return nil
}

// Unit event prefix (shared by all unit/spell events):
// [1] sourceGUID, [2] sourceName, [3] sourceFlags, [4] sourceRaidFlags,
// [5] targetGUID, [6] targetName, [7] targetFlags, [8] targetRaidFlags

// ENCOUNTER_START: [0] event, [1] encounterID, [2] name, [3] difficultyID, [4] groupSize, [5] instanceID
func parseEncounterStart(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 6, "ENCOUNTER_START"); err != nil {
		return err
	}
	encounterID, err := parseInt(fields[1])
	if err != nil {
		return err
	}
	difficultyID, err := parseInt(fields[3])
	if err != nil {
		return err
	}
	groupSize, err := parseInt(fields[4])
	if err != nil {
		return err
	}
	instanceID, err := parseInt(fields[5])
	if err != nil {
		return err
	}
	event.EncounterStart = EventEncounterStart{
		Timestamp:     ts,
		EncounterID:   encounterID,
		EncounterName: fields[2],
		DifficultyID:  difficultyID,
		GroupSize:     groupSize,
		InstanceID:    instanceID,
	}
	return nil
}

// ENCOUNTER_END: [0] event, [1] encounterID, [2] name, [3] difficultyID, [4] groupSize, [5] success, [6] duration
func parseEncounterEnd(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 7, "ENCOUNTER_END"); err != nil {
		return err
	}
	encounterID, err := parseInt(fields[1])
	if err != nil {
		return err
	}
	difficultyID, err := parseInt(fields[3])
	if err != nil {
		return err
	}
	groupSize, err := parseInt(fields[4])
	if err != nil {
		return err
	}
	success, err := parseBool(fields[5])
	if err != nil {
		return err
	}
	duration, err := parseInt(fields[6])
	if err != nil {
		return err
	}
	event.EncounterEnd = EventEncounterEnd{
		Timestamp:     ts,
		EncounterID:   encounterID,
		EncounterName: fields[2],
		DifficultyID:  difficultyID,
		GroupSize:     groupSize,
		Success:       success,
		Duration:      duration,
	}
	return nil
}

// MAP_CHANGE: [0] event, [1] instanceID, [2] name, [3..5] ignored
func parseMapChange(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 3, "MAP_CHANGE"); err != nil {
		return err
	}
	instanceID, err := parseInt(fields[1])
	if err != nil {
		return err
	}
	event.MapChange = EventMapChange{
		Timestamp:  ts,
		InstanceID: instanceID,
		Name:       fields[2],
	}
	return nil
}

// ZONE_CHANGE: [0] event, [1] instanceID, [2] name, [3] difficultyID
func parseZoneChange(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 4, "ZONE_CHANGE"); err != nil {
		return err
	}
	instanceID, err := parseInt(fields[1])
	if err != nil {
		return err
	}
	difficultyID, err := parseInt(fields[3])
	if err != nil {
		return err
	}
	event.ZoneChange = EventZoneChange{
		Timestamp:    ts,
		InstanceID:   instanceID,
		Name:         fields[2],
		DifficultyID: difficultyID,
	}
	return nil
}

// UNIT_DIED: uses the unit prefix; target is at [5],[6]
func parseUnitDied(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 9, "UNIT_DIED"); err != nil {
		return err
	}
	event.UnitDied = EventUnitDied{
		Timestamp:  ts,
		TargetGUID: fields[5],
		TargetName: fields[6],
	}
	return nil
}

// Spell prefix appended after unit prefix: [9] spellID, [10] spellName, [11] spellSchool
func parseSpellCastSuccess(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 12, "SPELL_CAST_SUCCESS"); err != nil {
		return err
	}
	spellID, err := parseInt(fields[9])
	if err != nil {
		return err
	}
	event.SpellCastSuccess = EventSpellCastSuccess{
		Timestamp:  ts,
		SourceGUID: fields[1],
		SourceName: fields[2],
		TargetGUID: fields[5],
		TargetName: fields[6],
		SpellID:    spellID,
		SpellName:  fields[10],
	}
	return nil
}

// SPELL_DAMAGE / SPELL_PERIODIC_DAMAGE:
// After the spell prefix (fields[9..11]), advanced fields follow (fields[12..29]).
// Damage is at [30], raw damage at [31].
func parseSpellDamage(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 32, "SPELL_DAMAGE"); err != nil {
		return err
	}
	spellID, err := parseInt(fields[9])
	if err != nil {
		return err
	}
	damage, err := parseInt(fields[30])
	if err != nil {
		return err
	}
	damageRaw, err := parseInt(fields[31])
	if err != nil {
		return err
	}
	event.SpellDamage = EventSpellDamage{
		Timestamp:  ts,
		SourceGUID: fields[1],
		SourceName: fields[2],
		TargetGUID: fields[5],
		TargetName: fields[6],
		SpellID:    spellID,
		SpellName:  fields[10],
		Damage:     damage,
		DamageRaw:  damageRaw,
	}
	return nil
}

// SPELL_HEAL / SPELL_PERIODIC_HEAL: amount at [30]
func parseSpellHeal(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 31, "SPELL_HEAL"); err != nil {
		return err
	}
	spellID, err := parseInt(fields[9])
	if err != nil {
		return err
	}
	amount, err := parseInt(fields[30])
	if err != nil {
		return err
	}
	event.SpellHeal = EventSpellHeal{
		Timestamp:  ts,
		SourceGUID: fields[1],
		SourceName: fields[2],
		TargetGUID: fields[5],
		TargetName: fields[6],
		SpellID:    spellID,
		SpellName:  fields[10],
		Amount:     amount,
	}
	return nil
}

// SPELL_AURA_APPLIED / _REFRESH / _REMOVED
func parseSpellAura(event *Event, ts time.Time, fields []string) error {
	if err := need(fields, 12, "SPELL_AURA"); err != nil {
		return err
	}
	spellID, err := parseInt(fields[9])
	if err != nil {
		return err
	}
	event.SpellAura = EventSpellAura{
		Timestamp:  ts,
		SourceGUID: fields[1],
		SourceName: fields[2],
		TargetGUID: fields[5],
		TargetName: fields[6],
		SpellID:    spellID,
		SpellName:  fields[10],
	}
	return nil
}
