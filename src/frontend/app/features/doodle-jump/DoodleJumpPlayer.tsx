import { useState, type RefObject } from 'react'
import { DotLottieReact } from '@lottiefiles/dotlottie-react'
import { type GameStatus, PLAYER_H, PLAYER_W } from './game'

interface DoodleJumpPlayerProps {
  containerRef: RefObject<HTMLDivElement | null>
  status: GameStatus
}

export function DoodleJumpPlayer({
  containerRef,
  status,
}: DoodleJumpPlayerProps) {
  const [failed, setFailed] = useState(false)
  const playing = status === 'playing'

  return (
    <div
      ref={containerRef}
      style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: PLAYER_W,
        height: PLAYER_H,
        pointerEvents: 'none',
        zIndex: 5,
      }}
    >
      {failed ? (
        <span
          data-testid="raccoon-fallback"
          aria-hidden="true"
          style={{
            display: 'flex',
            width: '100%',
            height: '100%',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: PLAYER_H * 0.8,
            lineHeight: 1,
          }}
        >
          🦝
        </span>
      ) : (
        <DotLottieReact
          src="/lottie/raccoon-adult.json"
          loop
          autoplay={playing}
          speed={playing ? 1 : 0}
          style={{ width: '100%', height: '100%', display: 'block' }}
          onError={() => setFailed(true)}
        />
      )}
    </div>
  )
}
