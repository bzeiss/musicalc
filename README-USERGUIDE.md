# User Guide

## Timecode

1. **Select frame rate** from the dropdown (23.976, 24, 25, 29.97, 30, etc.)
2. **Enter timecode values** in Hours:Minutes:Seconds:Frames format
   - Fields start empty - just click and type
   - **Right-justified entry**: Type digits and they fill from right to left (frames → seconds → minutes → hours)
   - **Dot shorthand**: Press `.` (period) to quickly insert `00` - great for fast entry (e.g., `1.` = 1 second)
   - **Backspace** removes the last entered digit
   - Empty fields are treated as 0 in calculations

3. **Perform calculations**:
   - **Add**: Combines Timecode 1 and Timecode 2, moves result to Timecode 1
   - **Subtract**: Subtracts Timecode 2 from Timecode 1, moves result to Timecode 1
   - **Reset**: Clears all timecode fields and history
   - **Clear History**: Clears only the history, keeps current timecode values

4. **Switch frame rates**: Change the FPS dropdown to see:
   - Same timecode notation (H:M:S:F stays constant)
   - Recalculated frame count for the new frame rate
   - History entry showing the conversion

5. **View history**: All calculations and conversions are logged with frame counts and FPS
   - Copy/paste history entries for documentation
   - History persists until cleared or reset

**Understanding the Display**:
- Format: `00:34:35:04 (62254 frames @ 30)`
- First part: Hours:Minutes:Seconds:Frames
- Parentheses: Total frame count at current FPS
- After @: Current frame rate

**Example Use Cases**:
- Add two video clips: `00:22:22:22 + 00:12:12:12 = 00:34:35:04 @ 30 fps`
- Convert timecode between formats: `00:34:35:04 (62254f) @30 → 00:34:35:04 (51879f) @25`
- Calculate exact frame counts for EDL or AAF workflows
- Verify timecode calculations for post-production deliverables

## Tempo to Delay

1. **Set your project tempo** in the BPM field at the top
2. **View delay times** for each note division in the table
3. **Use the values** to configure:
   - Delay effect times for rhythmic echoes
   - LFO rates for synchronized modulation
   - Gate/sequencer timing

**Example**: At 120 BPM, a 1/4 note = 500 ms, perfect for a quarter-note delay

## Note to Frequency

1. **Set reference frequency** (default: 440 Hz)
2. **Select reference note** if using non-standard tuning
3. **Choose middle C convention**: Switch between C3 and C4 naming conventions
4. **Browse the frequency table** to find:
   - Exact frequencies for synthesizer tuning
   - MIDI note numbers for programming
   - Pitch relationships between notes

**Example**: To tune a 808 kick to your track's key, find the root note frequency and adjust the kick's pitch to match

## Frequency to Note

1. **Enter frequency** in Hz (e.g., 440, 261.63, 1000)
2. **View the closest note** with pitch deviation in cents
3. **Choose notation system**:
   - **100 Cents**: Standard music notation (100 cents = 1 semitone)
   - **50 Cents**: Alternative notation used by some samplers (50 cents = 1 semitone)

4. **Adjust settings** (optional):
   - **Middle C**: Switch between C3 and C4 conventions
   - **Reference**: Set custom tuning reference (default A4 = 440 Hz)

5. **Use quick-select buttons** for common reference frequencies:
   - 220 Hz (A3), 440 Hz (A4), 880 Hz (A5), etc.
   - Instantly see the note and any pitch deviation

**Understanding the Output**:
- **100 Cents Notation**: Shows closest note and deviation (e.g., "F#3, -14 cents")
  - Positive cents = sharp (higher than note)
  - Negative cents = flat (lower than note)
- **50 Cents Notation**: Alternative display (e.g., "F3, +86 cents")
  - Same pitch, different reference point

**Example Use Cases**:
- Analyze a synth tone at 367 Hz → Result: F#3, -14 cents (or F3, +86 cents)
- Find what note a resonant frequency represents
- Tune acoustic instruments by analyzing recorded frequencies
- Identify pitch of sampled sounds for proper key mapping
- Match synth oscillators to specific frequencies

## Sample Length

1. **Choose your workflow**:
   - Enter **Tempo** to calculate sample length for a specific BPM
   - Enter **Length in samples** to find what tempo matches your loop
   - Enter **Length in ms** to calculate both tempo and sample count

2. **Set parameters**:
   - **Sample Rate**: Select your project's sample rate (typically 44100 or 48000)
   - **Beats**: Set loop length (4 = 4-beat loop, 8 = 8-beat loop, etc.)

3. **Read calculated values** instantly in all other fields

**Example Use Cases**:
- You have a 88200-sample loop at 44.1kHz → Calculator shows it's 2000ms and fits 120 BPM at 4 beats
- You want a 4-bar loop at 140 BPM → Calculator shows you need 151,543 samples at 44.1kHz
- You need to know what tempo your found sample is → Enter its sample count and get the BPM

## Multi-Mic Alignment Delay

This calculator helps you time-align multiple microphones (close mics) to a single **reference** microphone (usually the furthest mic, e.g. a room mic). It computes how much delay to add to each close mic so that all transients arrive at the DAW at the same time.

1. **Set sample rate**
   - Choose your DAW/project sample rate (e.g. 44100, 48000, 96000).

2. **Set temperature**
   - Enter the room temperature and select **C** or **F**.
   - This is used to estimate speed of sound.

3. **Set the reference distance**
   - Enter the distance to your furthest microphone (the one you align everything to).
   - Choose **m** or **ft**.

4. **Add target microphones**
   - Enter a mic name (pick from the dropdown or type a custom one).
   - Enter its distance.
   - Choose **m** or **ft**.
   - Press **+** to add it to the table.

5. **Read results**
   - Each row shows the mic distance, plus the required delay in:
     - **milliseconds (ms)**
     - **samples**

6. **Remove microphones**
   - Click the trash icon in the Actions column to remove a mic.

**Notes / edge cases**:
- If a mic is **further than the reference**, the delay is shown as **N/A** because you cannot add a positive delay to make a further mic arrive earlier.
- All values update automatically when you change temperature, units, sample rate, reference distance, or add/remove microphones.

## Tempo Change

1. **Set original tempo** of your audio material
2. **Choose your operation**:
   
   **Option A - Change Tempo Directly**:
   - Enter **New Tempo** to see pitch shift required
   - View **Time Stretch %** (e.g., 200% = double speed)
   - See **Tempo Delta** (percentage change)
   
   **Option B - Change Time Stretch**:
   - Enter **Time Stretching %** to see resulting tempo
   - View pitch changes in semitones and cents
   
   **Option C - Transpose**:
   - Enter **Transpose Semis** and **Cents** for desired pitch shift
   - View resulting tempo and time stretch amount

3. **Use the Swap button** to quickly reverse tempo/new tempo values
4. **Use Reset** to return to default values (Tempo: 140, New Tempo: 22)

**Understanding the Output**:
- **Transpose Semis/Cents**: Standard pitch notation (100 cents = 1 semitone)
- **50 Cents Notation**: Alternative pitch display (50 cents = 1 semitone) used by some samplers
- **Tempo Delta**: Percentage change from original tempo (+/- %)

**Example Use Cases**:
- Speed up a 140 BPM loop to 170 BPM → See it requires +3.45 semitones pitch shift
- You need to pitch a sample up 7 semitones → See it will play at 150% speed
- Match a sample's tempo to your project without changing pitch → Use time-stretching calculations
