package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func extractFromUint32(source uint32, size int, shift int) uint32 {
	rigthShift := 4*byteSize - (shift + size)
	leftShift := rigthShift + shift
	return (source << rigthShift) >> leftShift
}

func extractFromByte(source byte, size int, shift int) byte {
	rigthShift := byteSize - (shift + size)
	leftShift := rigthShift + shift
	return (source << rigthShift) >> leftShift
}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		n := len(name)
		if n > maxNameLen {
			n = maxNameLen
		}
		for i := 0; i < n; i++ {
			person.nameBytes[i] = name[i]
		}
		person.personTypeAndNameLen |= byte(n) << nameLenShift
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		m := mana << manaShift
		person.manaAndHelthAndFlags[2] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[1] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[0] |= byte(m)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		m := health << healthShift
		person.manaAndHelthAndFlags[2] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[1] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[0] |= byte(m)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectAndStrength |= uint8(respect) << respectShift
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectAndStrength |= uint8(strength) << strengthShift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experienceAndLevel |= uint8(experience) << experienceShift
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experienceAndLevel |= uint8(level) << levelShift
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		m := 1 << hasHouseShift
		person.manaAndHelthAndFlags[2] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[1] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[0] |= byte(m)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		m := 1 << hasGunShift
		person.manaAndHelthAndFlags[2] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[1] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[0] |= byte(m)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		m := 1 << hasFamilyShift
		person.manaAndHelthAndFlags[2] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[1] |= byte(m)
		m = m >> byteSize
		person.manaAndHelthAndFlags[0] |= byte(m)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.personTypeAndNameLen |= uint8(personType) << personTypeShift
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	byteSize = 8

	personTypeSize  = 2
	personTypeShift = 6

	nameLenSize  = 6
	nameLenShift = 0

	respectSize  = 4
	respectShift = 4

	strengthSize  = 4
	strengthShift = 0

	experienceSize  = 4
	experienceShift = 4

	levelSize  = 4
	levelShift = 0

	manaSize  = 10
	manaShift = 13

	healthSize  = 10
	healthShift = 3

	hasHouseShift  = 2
	hasGunShift    = 1
	hasFamilyShift = 0
)

const maxNameLen = 42

type GamePerson struct {
	// раскладываем по половинке байта попарно
	respectAndStrength byte
	experienceAndLevel byte

	// пакуем mana (10 bits), health (10 bits), hasHouse + hasGun + hasFamily (3 bits total) в 3 байта (24 бита)
	manaAndHelthAndFlags [3]byte

	// пакуем тип персоны (2 бита) и длину имени (6 бит) в один байт
	personTypeAndNameLen byte

	// массив байт с максимальной емкостью, предусмотренной задачей (42 байта)
	nameBytes [maxNameLen]byte

	// так как выше битики ровно упакуются по 4 байта, то 4 байта выравнивания допустимо
	// в связи с этим можно тут хранить все в int32 числа без битовой магии
	x, y, z, gold int32

	// Всего получится как раз 64 байта
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}

	for _, opt := range options {
		opt(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	n := int(extractFromByte(p.personTypeAndNameLen, nameLenSize, nameLenShift))
	return unsafe.String(&p.nameBytes[0], n)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	m := uint32(0)
	m |= uint32(p.manaAndHelthAndFlags[0]) << (2 * byteSize)
	m |= uint32(p.manaAndHelthAndFlags[1]) << byteSize
	m |= uint32(p.manaAndHelthAndFlags[2])

	return int(extractFromUint32(m, manaSize, manaShift))
}

func (p *GamePerson) Health() int {
	m := uint32(0)
	m |= uint32(p.manaAndHelthAndFlags[0]) << (2 * byteSize)
	m |= uint32(p.manaAndHelthAndFlags[1]) << byteSize
	m |= uint32(p.manaAndHelthAndFlags[2])

	return int(extractFromUint32(m, healthSize, healthShift))
}

func (p *GamePerson) Respect() int {
	return int(extractFromByte(p.respectAndStrength, respectSize, respectShift))
}

func (p *GamePerson) Strength() int {
	return int(extractFromByte(p.respectAndStrength, strengthSize, strengthShift))
}

func (p *GamePerson) Experience() int {
	return int(extractFromByte(p.experienceAndLevel, experienceSize, experienceShift))
}

func (p *GamePerson) Level() int {
	return int(extractFromByte(p.experienceAndLevel, levelSize, levelShift))
}

func (p *GamePerson) HasHouse() bool {
	return extractFromByte(p.manaAndHelthAndFlags[2], 1, hasHouseShift) > 0
}

func (p *GamePerson) HasGun() bool {
	return extractFromByte(p.manaAndHelthAndFlags[2], 1, hasGunShift) > 0
}

func (p *GamePerson) HasFamilty() bool {
	return extractFromByte(p.manaAndHelthAndFlags[2], 1, hasFamilyShift) > 0
}

func (p *GamePerson) Type() int {
	return int(extractFromByte(p.personTypeAndNameLen, personTypeSize, personTypeShift))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
