import { useState } from 'react'
import { PetAvatar } from './PetAvatar'
import type { PetEmotion, PetStage } from './types'

const STAGES: PetStage[] = ['egg', 'baby', 'teen', 'adult', 'legend']

const MOOD_PRESETS = [
  { key: 'happy', title: 'happy', satiety: 90, happiness: 90, energy: 90 },
  { key: 'idle', title: 'idle', satiety: 60, happiness: 50, energy: 60 },
  { key: 'hungry', title: 'hungry', satiety: 15, happiness: 70, energy: 70 },
  { key: 'sad', title: 'sad', satiety: 70, happiness: 12, energy: 70 },
  { key: 'sleeping', title: 'sleeping', satiety: 70, happiness: 70, energy: 8 },
] as const

const EMOTIONS: PetEmotion[] = ['eating', 'celebrate', 'levelup', 'hatching']

const cellStyle = {
  display: 'flex',
  flexDirection: 'column' as const,
  alignItems: 'center',
  gap: '0.5rem',
  padding: '0.75rem',
  borderRadius: '0.75rem',
  background: 'var(--card, #fff)',
  border: '1px solid var(--border, #e5e5e5)',
}

const captionStyle = {
  fontSize: '0.75rem',
  color: 'var(--muted-foreground, #666)',
  fontVariantNumeric: 'tabular-nums' as const,
}

export function PetAvatarDemo() {
  const [emotion, setEmotion] = useState<PetEmotion | null>(null)
  const [strokes, setStrokes] = useState(0)

  return (
    <div
      style={{
        padding: '2rem',
        display: 'flex',
        flexDirection: 'column',
        gap: '2rem',
        background: 'var(--background, #fafafa)',
        color: 'var(--foreground, #111)',
        minHeight: '100vh',
      }}
    >
      <header
        style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}
      >
        <h1 style={{ fontSize: '1.5rem', fontWeight: 600 }}>
          Енот Ноти — витрина аватара
        </h1>
        <p style={captionStyle}>
          Поглаживаний: {strokes}. Наведите курсор — зрачки следят. Не трогайте
          мышь 45 секунд — питомец заснёт.
        </p>
      </header>

      <section style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
        {EMOTIONS.map((value) => (
          <button
            key={value}
            type="button"
            onClick={() => setEmotion(value)}
            style={{
              padding: '0.5rem 0.9rem',
              borderRadius: '0.5rem',
              border: '1px solid var(--border, #e5e5e5)',
              background:
                emotion === value ? 'var(--primary, #0af)' : 'transparent',
              color: emotion === value ? '#fff' : 'inherit',
              cursor: 'pointer',
            }}
          >
            {value}
          </button>
        ))}
        <button
          type="button"
          onClick={() => setEmotion(null)}
          style={{
            padding: '0.5rem 0.9rem',
            borderRadius: '0.5rem',
            border: '1px solid var(--border, #e5e5e5)',
            background: 'transparent',
            color: 'inherit',
            cursor: 'pointer',
          }}
        >
          сбросить
        </button>
      </section>

      {MOOD_PRESETS.map((preset) => (
        <section
          key={preset.key}
          style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}
        >
          <h2 style={{ fontSize: '1rem', fontWeight: 600 }}>{preset.title}</h2>
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(160px, 1fr))',
              gap: '1rem',
            }}
          >
            {STAGES.map((stage) => (
              <div key={stage} style={cellStyle}>
                <PetAvatar
                  stage={stage}
                  emotion={emotion}
                  satiety={preset.satiety}
                  happiness={preset.happiness}
                  energy={preset.energy}
                  size="md"
                  idleSleepEnabled={false}
                  onStroke={() => setStrokes((value) => value + 1)}
                />
                <span style={captionStyle}>{stage}</span>
              </div>
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}
