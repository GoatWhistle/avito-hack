import { useGetMyRaccoon } from '#/features/tamagotchi/hooks/useGetMyRaccoon'
import { DotLottieReact } from '@lottiefiles/dotlottie-react'

export function Raccoon() {
  const { data, isLoading, error } = useGetMyRaccoon()

  if (!data || isLoading || error) return null

  return (
    <DotLottieReact
      renderConfig={{ autoResize: true }}
      src={`/data-${data.stage}.json`}
      loop
      autoplay
      height={150}
    />
  )
}
