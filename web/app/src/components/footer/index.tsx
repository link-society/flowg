import { useFeatureFlags } from '@/lib/hooks/featureflags'

export const Footer = () => {
  const { demoMode } = useFeatureFlags()

  return (
    <footer className="flex flex-row p-3 bg-gray-300">
      {demoMode && (
        <div className="font-semibold">
          <p>Demo Mode Enabled</p>
        </div>
      )}
      <div className="ml-auto font-semibold">
        <p>{import.meta.env.FLOWG_VERSION}</p>
      </div>
    </footer>
  )
}
