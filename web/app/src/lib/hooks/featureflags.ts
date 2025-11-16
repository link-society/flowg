type FeatureFlags = {
  demoMode: boolean
}

export const useFeatureFlags = (): FeatureFlags => {
  const demoModeMeta = document
    .querySelector('meta[name="X-FlowG-DemoMode"]')
    ?.getAttribute('content')

  return {
    demoMode: demoModeMeta === 'true',
  }
}
