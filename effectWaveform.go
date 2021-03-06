package main

import "math"

/*

4XY

VIBRATO means to "oscillate the sample pitch using a  particular waveform
with amplitude yyyy notes, such that (xxxx * speed)/64  full oscillations
occur in the line".  The waveform to use in vibrating is set using effect
E4 (see  below). By  placing vibrato  effects on  consecutive lines,  the
vibrato effect can be sustained for  any length of time.  If  either xxxx
or yyyy are 0,  then values from  the most recent  prior vibrato will  be
used.

An example is: Note C-3, with xxxx=12 and yyyy=1 when speed=8.  This will
play tones  around  C-3,  vibrating through  D-3  and  B-2 to  C-3  again
(amplitude  yyyy  is 1), with (12*8)/64 = 1.5 full oscillations per line.

8 /

(x*sp)/64 = 2𝛑

Please see effect E4 for the waveform to use for vibrating.

FIXME: notes or HALF-NOTES (SEMITONES)?


7XY

TREMOLO means to "oscillate the sample volume using a particular waveform
with   amplitude   yyyy*(speed-1),   such   that   (xxxx*speed)/64   full
oscillations occur in the line".  The waveform to use to oscillate is set
using the  effect  E7  (see  below).    By  placing  tremolo  effects  on
consecutive lines, the tremolo effect can be sustained for any  length of
time.  If either  xxxx or yyyy are  0, then values  from the most  recent
prior tremolo will be used.

The usage of this effect is similar to that of effect 4:Vibrato.

*/

// WaveformType is the type of the envelope waveform used for vibrato/tremolo
type WaveformType int

const (
	// Sine - sine waveform
	Sine WaveformType = iota
	// RampDown - "sawtooth" waveform
	RampDown
	// Square - square wave
	Square
	// Random - select one of Sine/RampDown/Square randomly
	Random
)

// EffectWaveform contains the parameters for a waveform assigned to an effect
type EffectWaveform struct {
	SamplesPerTick int

	Active bool

	Type   WaveformType
	Retrig bool

	CurType   WaveformType
	Pos       float64
	Step      float64
	Amplitude float64
}

// DoStep gets the next value for our waveform
func (ew *EffectWaveform) DoStep() int {
	if !ew.Active {
		return 0
	}
	ew.Pos += ew.Step
	switch ew.CurType {
	case Sine:
		return int(math.Round(ew.Amplitude * math.Sin(ew.Pos)))
	case Square, RampDown: // FIXME implement RampDown!
		if math.Sin(ew.Pos) > 0 {
			return int(ew.Amplitude)
		}
		return int(-ew.Amplitude)
	}
	return 0
}

func (ew *EffectWaveform) initWaveform(X, amplitude int) {
	ew.Active = true
	if ew.Retrig {
		ew.Pos = 0
	}
	if X > 0 && amplitude > 0 {
		ew.CurType = ew.Type
		if ew.Type == Random {
			ew.CurType = Sine // TODO: really set type randomly!
		}
		ew.Step = (math.Pi * float64(X)) / (32.0 * float64(ew.SamplesPerTick))
		ew.Amplitude = float64(amplitude)
	}
}

// InitTremoloWaveform (re)initializes a waveform for a tremolo (volume) effect at the beginning of a new note
func (ew *EffectWaveform) InitTremoloWaveform(X, Y int) {
	ew.initWaveform(X, Y)
}

// InitVibratoWaveform initializes a waveform for a vibrato (pitch) effect
func (ew *EffectWaveform) InitVibratoWaveform(X, Y, period int, ins Instrument) {
	ew.initWaveform(X, ins.GetPeriodDelta(period, Y))
}

// DecodeWaveformType sets the type of an EffectWaveform from a "set waveform" command parameter par
func (ew *EffectWaveform) DecodeWaveformType(par int) {
	ew.Retrig = par < 4
	switch par {
	case 0, 4:
		ew.Type = Sine
	case 1, 5:
		ew.Type = RampDown
	case 2, 6:
		ew.Type = Square
	case 3, 7:
		ew.Type = Random
	}
}

// NewEffectWaveform creates a new EffectWaveform
func NewEffectWaveform(SamplesPerTick int) (ew EffectWaveform) {
	ew.SamplesPerTick = SamplesPerTick
	ew.DecodeWaveformType(0)
	return
}
